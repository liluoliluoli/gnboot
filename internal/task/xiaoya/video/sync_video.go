package video

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/jellyfindto"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/xiaoyadto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"math"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type XiaoyaVideoTask struct {
	episodeRepo *repo.EpisodeRepo
	videoRepo   *repo.VideoRepo
	client      redis.UniversalClient
	c           *conf.Bootstrap
	configRepo  *repo.ConfigRepo
}

func NewXiaoyaVideoTask(episodeRepo *repo.EpisodeRepo, videoRepo *repo.VideoRepo, client redis.UniversalClient, c *conf.Bootstrap, configRepo *repo.ConfigRepo) *XiaoyaVideoTask {
	return &XiaoyaVideoTask{
		episodeRepo: episodeRepo,
		videoRepo:   videoRepo,
		client:      client,
		c:           c,
		configRepo:  configRepo,
	}
}

func (t *XiaoyaVideoTask) Process(task *sdomain.Task) error {
	ctx := task.Ctx
	xiaoyaBoxIpStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_XiaoYaBoxIp)
	if err != nil {
		return err
	}
	xiaoyaBoxIps := strings.Split(xiaoyaBoxIpStr, ",")
	jellyfinBoxIpStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_JellyfinBoxIp)
	if err != nil {
		return err
	}
	jellyfinBoxIps := strings.Split(jellyfinBoxIpStr, ",")
	scanPathStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_XiaoyaVideoSyncCategory)
	if err != nil {
		return err
	}
	scanPaths := strings.Split(scanPathStr, ",")
	pageSize := int32(100)
	for _, scanPath := range scanPaths {
		mappingPath := strings.Split(scanPath, "::")
		log.Infof("开始同步xiaoya视频: syncURL=%s, initialPath=%s", constant.XiaoYaVideoList, mappingPath[0])
		err := t.deepLoopXiaoYaPath(ctx, xiaoyaBoxIps[0], mappingPath[0], "", pageSize)
		if err != nil {
			log.Errorf("同步xiaoya视频失败: %v", err)
			return err
		}
		log.Infof("结束同步xiaoya视频: syncURL=%s, initialPath=%s", constant.XiaoYaVideoList, mappingPath[0])

		log.Infof("开始同步jellyfin刮削: syncURL=%s, pathid=%s, type=%s", jellyfinBoxIps[0], mappingPath[1], mappingPath[2])
		err = t.deepLoopJellyfinPath(ctx, jellyfinBoxIps[0], mappingPath[1], mappingPath[2], pageSize)
		if err != nil {
			log.Errorf("同步jellyfin刮削失败: %v", err)
			return err
		}
		log.Infof("结束同步jellyfin刮削: syncURL=%s, pathid=%s, type=%s", jellyfinBoxIps[0], mappingPath[1], mappingPath[2])
	}
	return nil
}

func (t *XiaoyaVideoTask) deepLoopXiaoYaPath(ctx context.Context, domain, currentPath, parentPath string, pageSize int32) error {
	page := int32(1)
	for {
		result, err := httpclient_util.DoPost[xiaoyadto.VideoListReq, xiaoyadto.XiaoyaResult[xiaoyadto.VideoListResp]](ctx, domain+constant.XiaoYaVideoList, "", &xiaoyadto.VideoListReq{
			Path:     currentPath,
			Password: "",
			Page:     page,
			PerPage:  pageSize,
			Refresh:  false,
		})

		if err != nil {
			return fmt.Errorf("xiaoya请求路径 %s 第 %d 页失败: %v", currentPath, page, err)
		}

		if result == nil || result.Code != 200 || result.Data == nil {
			return fmt.Errorf("xiaoya路径 %s 第 %d 页返回结果无效", currentPath, page)
		}

		for _, content := range result.Data.Content {
			log.Infof("xiaoya处理内容: 路径=%s, 父级=%s, 名称=%s, 是目录=%v", currentPath, parentPath, content.Name, content.IsDir)
			if content.IsDir {
				newPath := path.Join(currentPath, content.Name)
				err := t.deepLoopXiaoYaPath(ctx, domain, newPath, content.Name, pageSize)
				if err != nil {
					log.Errorf("递归查询小雅失败: %v", err)
					return err
				}
			} else {
				existEpisode, err := t.episodeRepo.QueryByPathAndName(ctx, currentPath, content.Name) // 查询时使用currentPath
				if err != nil {
					log.Errorf("查询episode失败: %v，跳过当前文件", err)
					continue
				}
				if len(existEpisode) == 0 {
					episode := &model.Episode{
						XiaoyaPath:   &currentPath,
						EpisodeTitle: content.Name,
						Size:         fmt.Sprintf("%d", content.Size),
						Platform:     lo.ToPtr(constant.Platform),
						IsValid:      true,
						CreateTime:   time.Now(),
						UpdateTime:   time.Now(),
					}
					if err := t.episodeRepo.Create(ctx, nil, episode); err != nil {
						log.Errorf("写入episode失败（path=%s, title=%s）: %v", currentPath, content.Name, err)
					} else {
						log.Infof("成功写入episode: path=%s, title=%s", currentPath, content.Name)
					}
				} else {
					log.Infof("episode已存在（跳过）: path=%s, title=%s", currentPath, content.Name)
				}
			}
		}

		if int64(page)*int64(pageSize) >= result.Data.Total {
			break
		}
		page++
	}

	return nil
}

func (t *XiaoyaVideoTask) deepLoopJellyfinPath(ctx context.Context, domain, parentId string, mediaType string, pageSize int32) error {
	startIndex := int32(0)
	jellyfinDefaultToken, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_JellyfinDefaultToken)
	if err != nil {
		return err
	}
	jellyfinDefaultUserId, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_JellyfinDefaultUserId)
	if err != nil {
		return err
	}
	for {
		headerMap := make(map[string]string)
		headerMap["x-emby-authorization"] = jellyfinDefaultToken
		syncListURL := fmt.Sprintf(domain+constant.JellyfinVideoList, jellyfinDefaultUserId, mediaType, startIndex, parentId, pageSize)
		videoListResp, err := httpclient_util.DoGet[jellyfindto.VideoListResp](ctx, syncListURL, headerMap)
		if err != nil {
			return fmt.Errorf("jellyfin请求parentId %s 游标 %d 开始的返回结果失败: %v", parentId, startIndex, err)
		}
		if videoListResp == nil {
			return fmt.Errorf("jellyfin parentId %s 游标 %d 开始的返回结果无效", parentId, startIndex)
		}
		for _, content := range videoListResp.Items {
			log.Infof("jellyfin处理内容: 父级=%s, 名称=%s, 是目录=%v", parentId, content.Name, content.IsFolder)
			syncDetailURL := fmt.Sprintf(domain+constant.JellyfinVideoDetail, jellyfinDefaultUserId, content.Id)
			videoDetailResp, err := httpclient_util.DoGet[jellyfindto.VideoDetailResp](ctx, syncDetailURL, headerMap)
			if err != nil {
				return fmt.Errorf("jellyfin parentId %s file %s 返回结果失败: %v", parentId, content.Id, err)
			}
			if videoDetailResp == nil {
				return fmt.Errorf("jellyfin parentId %s file %s 返回结果无效", parentId, content.Id)
			}
			filePath, fileName := splitJellyfinPath(videoDetailResp.MediaSources[0].Path)
			existEpisode, err := t.episodeRepo.QueryByPathAndName(ctx, filePath, fileName)
			if err != nil {
				log.Errorf("查询episode失败: %v，跳过当前文件", err)
				continue
			}
			if len(existEpisode) == 0 {
				log.Errorf("查询episode为空: filePath %s fileName %s", filePath, fileName)
				continue
			}
			if existEpisode[0].VideoId == 0 {
				video, err := t.videoRepo.GetByJellyfinId(ctx, videoDetailResp.Id)
				if err != nil {
					log.Errorf("查询video失败: %v", err)
					return err
				}
				err = gen.Use(t.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
					if video == nil {
						video = &model.Video{
							Title:              videoDetailResp.Name,
							VideoType:          mediaType,
							VoteRate:           lo.ToPtr(float32(math.Round(videoDetailResp.GoodRating*10) / 10)),
							Region:             lo.ToPtr(strings.Join(t.replaceRegions(ctx, videoDetailResp.Regions), ",")),
							Description:        lo.ToPtr(videoDetailResp.Overview),
							PublishDay:         lo.ToPtr(time_util.FormatStrToYYYYMMDD(videoDetailResp.PremiereDate)),
							Thumbnail:          lo.ToPtr(fmt.Sprintf(constant.PrimaryThumbnail, videoDetailResp.Id)),
							Genres:             lo.ToPtr(strings.Join(t.replaceGenres(ctx, videoDetailResp.Genres), ",")),
							JellyfinID:         videoDetailResp.Id,
							JellyfinCreateTime: time_util.ParseUtcTime(videoDetailResp.DateCreated),
						}
						err := t.videoRepo.Create(ctx, nil, video)
						if err != nil {
							log.Errorf("创建video失败: %v", err)
							return err
						}
					}
					existEpisode[0].VideoId = video.ID
					err = t.episodeRepo.Updates(ctx, nil, existEpisode[0])
					if err != nil {
						log.Errorf("更新episode失败: %v", err)
						return err
					}
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}

			}
		}
		if int64(startIndex+pageSize) >= videoListResp.TotalRecordCount {
			break
		}
		startIndex = startIndex + pageSize
	}

	return nil
}

// splitJellyfinPath http://xiaoya.host:5678/d/电影/4K系列/Marvel漫威系列/漫威宇宙/蜘蛛侠：英雄无归.Spider-Man.No.Way.Home.2021.UHD.BluRay.2160p.x265.10bit.DoVi.2Audios.mUHD-FRDS.mkv
func splitJellyfinPath(url string) (string, string) {
	domain := strings.Split(url, "/")[2] // "xiaoya.host:5678"

	pathStart := strings.Index(url, domain+"/d") + len(domain+"/d")
	pathEnd := strings.LastIndex(url, "/")
	path := url[pathStart:pathEnd] // "/电影/4K系列/Marvel漫威系列/漫威宇宙"

	// 3. 提取文件名
	filename := filepath.Base(url) // "蜘蛛侠：英雄无归.Spider-Man.No.Way.Home.2021.UHD.BluRay.2160p.x265.10bit.DoVi.2Audios.mUHD-FRDS.mkv"

	return path, filename
}

func (t *XiaoyaVideoTask) replaceRegions(ctx context.Context, regions []string) []string {
	return lo.Map(regions, func(item string, index int) string {
		newRegion, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_RegionMapping, item)
		if err != nil {
			return item
		}
		if newRegion == "" {
			return item
		}
		return newRegion
	})
}

func (t *XiaoyaVideoTask) replaceGenres(ctx context.Context, genres []string) []string {
	temp := lo.Map(genres, func(item string, index int) string {
		newGenre, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_GenreMapping, item)
		if err != nil {
			return ""
		}
		if newGenre == "" {
			return ""
		}
		return newGenre
	})
	return lo.Filter(temp, func(item string, index int) bool {
		return item != ""
	})
}
