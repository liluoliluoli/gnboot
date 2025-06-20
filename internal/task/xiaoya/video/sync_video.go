package video

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/jellyfindto"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/xiaoyadto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"path"
	"path/filepath"
	"strings"

	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type XiaoyaVideoTask struct {
	episodeRepo *repo.EpisodeRepo
	videoRepo   *repo.VideoRepo
	client      redis.UniversalClient
	c           *conf.Bootstrap
}

func NewXiaoyaVideoTask(episodeRepo *repo.EpisodeRepo, videoRepo *repo.VideoRepo, client redis.UniversalClient, c *conf.Bootstrap) *XiaoyaVideoTask {
	return &XiaoyaVideoTask{
		episodeRepo: episodeRepo,
		videoRepo:   videoRepo,
		client:      client,
		c:           c,
	}
}

func (t *XiaoyaVideoTask) Process(task *sdomain.Task) error {
	ctx := task.Ctx
	boxIps, err := json_util.Unmarshal[[]map[string]string](gerror.HandleRedisStringNotFound(t.client.Get(ctx, constant.RK_BoxIps).Val()))
	if err != nil {
		return err
	}
	scanPaths := strings.Split(t.c.XiaoyaVideoSyncCategory, ",")
	pageSize := int32(100)
	for _, scanPath := range scanPaths {
		mappingPath := strings.Split(scanPath, "::")
		log.Infof("开始同步xiaoya视频: syncURL=%s, initialPath=%s", constant.XiaoYaVideoList, mappingPath[0])
		err := t.deepLoopXiaoYaPath(ctx, boxIps[0][constant.Key_XiaoYaBoxIp], mappingPath[0], "", pageSize)
		if err != nil {
			return err
		}
		log.Infof("结束同步xiaoya视频: syncURL=%s, initialPath=%s", constant.XiaoYaVideoList, mappingPath[0])

		log.Infof("开始同步jellyfin刮削: syncURL=%s, pathid=%s, type=%s", boxIps[0][constant.Key_JellyfinBoxIp], mappingPath[1], mappingPath[2])
		err = t.deepLoopJellyfinPath(ctx, boxIps[0][constant.Key_JellyfinBoxIp], mappingPath[1], mappingPath[2], pageSize)
		if err != nil {
			return err
		}
		log.Infof("结束同步jellyfin刮削: syncURL=%s, pathid=%s, type=%s", boxIps[0][constant.Key_JellyfinBoxIp], mappingPath[1], mappingPath[2])
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
			// 记录当前处理的内容信息（包含父级路径）
			log.Infof("xiaoya处理内容: 路径=%s, 父级=%s, 名称=%s, 是目录=%v", currentPath, parentPath, content.Name, content.IsDir)
			if content.IsDir {
				newPath := path.Join(currentPath, content.Name)
				err := t.deepLoopXiaoYaPath(ctx, domain+constant.XiaoYaVideoList, newPath, content.Name, pageSize)
				if err != nil {
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
	for {
		headerMap := make(map[string]string)
		headerMap["x-emby-authorization"] = t.c.JellyfinDefaultToken
		syncListURL := fmt.Sprintf(domain+constant.JellyfinVideoList, t.c.JellyfinDefaultUserId, mediaType, startIndex, parentId, pageSize)
		videoListResp, err := httpclient_util.DoGet[jellyfindto.VideoListResp](ctx, syncListURL, headerMap)
		if err != nil {
			return fmt.Errorf("jellyfin请求parentId %s 游标 %d 开始的返回结果失败: %v", parentId, startIndex, err)
		}
		if videoListResp == nil {
			return fmt.Errorf("jellyfin parentId %s 游标 %d 开始的返回结果无效", parentId, startIndex)
		}
		for _, content := range videoListResp.Items {
			log.Infof("jellyfin处理内容: 父级=%s, 名称=%s, 是目录=%v", parentId, content.Name, content.IsFolder)
			syncDetailURL := fmt.Sprintf(domain+constant.JellyfinVideoDetail, t.c.JellyfinDefaultUserId, content.Id)
			videoDetailResp, err := httpclient_util.DoGet[jellyfindto.VideoDetailResp](ctx, syncDetailURL, headerMap)
			if err != nil {
				return fmt.Errorf("jellyfin parentId %s file %s 返回结果失败: %v", parentId, content.Id, err)
			}
			if videoDetailResp == nil {
				return fmt.Errorf("jellyfin parentId %s file %s 返回结果无效", parentId, content.Id)
			}
			filePath, fileName := splitJellyfinPath(videoDetailResp.MediaSources[0].Path)
			existEpisode, err := t.episodeRepo.QueryByPathAndName(ctx, filePath, fileName)
			if err != nil || len(existEpisode) == 0 {
				log.Errorf("查询episode失败: %v，跳过当前文件", err)
				continue
			}
			if existEpisode[0].VideoId == 0 {
				video := &sdomain.Video{
					Title:        videoDetailResp.Name,
					VideoType:    mediaType,
					VoteRate:     float32(videoDetailResp.GoodRating),
					Region:       strings.Join(videoDetailResp.Regions, ","),
					Description:  videoDetailResp.Overview,
					PublishMonth: time_util.FormatYYYYMM(videoDetailResp.PremiereDate),
					Thumbnail:    fmt.Sprintf(constant.PrimaryThumbnail, videoDetailResp.Id),
					Genres:       videoDetailResp.Genres,
				}
				err := t.videoRepo.Create(ctx, nil, video)
				if err != nil {
					log.Errorf("创建video失败: %v", err)
					return err
				}
				existEpisode[0].VideoId = video.ID
				err = t.episodeRepo.Updates(ctx, nil, existEpisode[0])
				if err != nil {
					log.Errorf("更新episode失败: %v", err)
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
