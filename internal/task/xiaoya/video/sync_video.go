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
	episodeRepo                *repo.EpisodeRepo
	videoRepo                  *repo.VideoRepo
	client                     redis.UniversalClient
	c                          *conf.Bootstrap
	configRepo                 *repo.ConfigRepo
	actorRepo                  *repo.ActorRepo
	videoActorMappingRepo      *repo.VideoActorMappingRepo
	episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
}

func NewXiaoyaVideoTask(episodeRepo *repo.EpisodeRepo, videoRepo *repo.VideoRepo, client redis.UniversalClient, c *conf.Bootstrap, configRepo *repo.ConfigRepo, actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo, episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo) *XiaoyaVideoTask {
	return &XiaoyaVideoTask{
		episodeRepo:                episodeRepo,
		videoRepo:                  videoRepo,
		client:                     client,
		c:                          c,
		configRepo:                 configRepo,
		actorRepo:                  actorRepo,
		videoActorMappingRepo:      videoActorMappingRepo,
		episodeSubtitleMappingRepo: episodeSubtitleMappingRepo,
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
	for _, scanPath := range scanPaths {
		mappingPath := strings.Split(scanPath, "::")
		log.Infof("开始同步xiaoya视频: syncURL=%s, initialPath=%s", constant.XiaoYaVideoList, mappingPath[0])
		err := t.deepLoopXiaoYaPath(ctx, xiaoyaBoxIps[0], mappingPath[0], "")
		if err != nil {
			log.Errorf("同步xiaoya视频失败: %v", err)
			return err
		}
		log.Infof("结束同步xiaoya视频: syncURL=%s, initialPath=%s", constant.XiaoYaVideoList, mappingPath[0])

		log.Infof("开始同步jellyfin刮削: syncURL=%s, pathid=%s", jellyfinBoxIps[0], mappingPath[1])
		err = t.deepLoopListJellyfinPath(ctx, jellyfinBoxIps[0], mappingPath[1])
		if err != nil {
			log.Errorf("同步jellyfin刮削失败: %v", err)
			return err
		}
		log.Infof("结束同步jellyfin刮削: syncURL=%s, pathid=%s, type=%s", jellyfinBoxIps[0], mappingPath[1])
	}
	return nil
}

func (t *XiaoyaVideoTask) deepLoopXiaoYaPath(ctx context.Context, domain, currentPath, parentPath string) error {
	page := int32(1)
	headerMap := make(map[string]string)
	headerMap["Authorization"] = constant.XiaoYaToken
	for {
		result, err := httpclient_util.DoPost[xiaoyadto.VideoListReq, xiaoyadto.XiaoyaResult[xiaoyadto.VideoListResp]](ctx, domain+constant.XiaoYaVideoList, &xiaoyadto.VideoListReq{
			Path:     currentPath,
			Password: "",
			Page:     page,
			PerPage:  constant.PageSize,
			Refresh:  false,
		}, headerMap)

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
				err := t.deepLoopXiaoYaPath(ctx, domain, newPath, content.Name)
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

		if int64(page)*int64(constant.PageSize) >= result.Data.Total {
			break
		}
		page++
	}

	return nil
}

func (t *XiaoyaVideoTask) deepLoopListJellyfinPath(ctx context.Context, domain, parentId string) error {
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
		syncListURL := fmt.Sprintf(domain+constant.JellyfinVideoList, jellyfinDefaultUserId, startIndex, parentId, constant.PageSize)
		videoListResp, err := httpclient_util.DoGet[jellyfindto.VideoListResp](ctx, syncListURL, headerMap)
		if err != nil {
			return fmt.Errorf("jellyfin请求parentId %s 游标 %d 开始的返回结果失败: %v", parentId, startIndex, err)
		}
		if videoListResp == nil {
			return fmt.Errorf("jellyfin parentId %s 游标 %d 开始的返回结果无效", parentId, startIndex)
		}
		for _, content := range videoListResp.Items {
			err = t.deepLoopDetailJellyfinPath(ctx, domain, parentId, content.Id, content.Type, nil, jellyfinDefaultUserId, headerMap)
			if err != nil {
				return err
			}
		}
		if int64(startIndex+constant.PageSize) >= videoListResp.TotalRecordCount {
			break
		}
		startIndex = startIndex + constant.PageSize
	}
	return nil
}

func (t *XiaoyaVideoTask) deepLoopDetailJellyfinPath(ctx context.Context, domain, parentId string, id string, mediaType string, season *jellyfindto.VideoDetailResp, jellyfinDefaultUserId string, headerMap map[string]string) error {
	syncDetailURL := fmt.Sprintf(domain+constant.JellyfinVideoDetail, jellyfinDefaultUserId, id)
	videoDetailResp, err := httpclient_util.DoGet[jellyfindto.VideoDetailResp](ctx, syncDetailURL, headerMap)
	if err != nil {
		return fmt.Errorf("jellyfin videoDetailResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
	}
	if videoDetailResp == nil {
		return fmt.Errorf("jellyfin videoDetailResp parentId %s file %s 返回结果无效", parentId, id)
	}
	if videoDetailResp.IsFolder {
		if mediaType == "Folder" {
			return t.deepLoopListJellyfinPath(ctx, domain, videoDetailResp.Id)
		}
		if mediaType == "Series" {
			syncSeasonListURL := fmt.Sprintf(domain+constant.JellyfinSeaonsList, id)
			seasonListResp, err := httpclient_util.DoGet[jellyfindto.SeasonListResp](ctx, syncSeasonListURL, headerMap)
			if err != nil {
				return fmt.Errorf("jellyfin seasonListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
			}
			if seasonListResp == nil {
				return fmt.Errorf("jellyfin seasonListResp parentId %s file %s 返回结果无效", parentId, id)
			}
			for _, content := range seasonListResp.Items {
				err := t.deepLoopDetailJellyfinPath(ctx, domain, id, content.Id, "Season", nil, jellyfinDefaultUserId, headerMap)
				if err != nil {
					return err
				}
			}
		}
		if mediaType == "Season" {
			syncEpisodeListURL := fmt.Sprintf(domain+constant.JellyfinEpisodesList, id)
			episodeListResp, err := httpclient_util.DoGet[jellyfindto.EpisodeListResp](ctx, syncEpisodeListURL, headerMap)
			if err != nil {
				return fmt.Errorf("jellyfin episodeListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
			}
			if episodeListResp == nil {
				return fmt.Errorf("jellyfin episodeListResp parentId %s file %s 返回结果无效", parentId, id)
			}
			for _, content := range episodeListResp.Items {
				err := t.deepLoopDetailJellyfinPath(ctx, domain, id, content.Id, "Episode", videoDetailResp, jellyfinDefaultUserId, headerMap)
				if err != nil {
					return err
				}
			}
		}
	} else {
		filePath, fileName := splitJellyfinPath(videoDetailResp.MediaSources[0].Path)
		existEpisode, err := t.episodeRepo.QueryByPathAndName(ctx, filePath, fileName)
		if err != nil {
			log.Errorf("查询episode失败: %v，跳过当前文件", err)
		}
		if len(existEpisode) == 0 {
			log.Errorf("查询episode为空: filePath %s fileName %s", filePath, fileName)
		}
		if existEpisode[0].VideoId == 0 {
			jellyfinId := ""
			if videoDetailResp.Type == "Movie" {
				jellyfinId = videoDetailResp.Id
			}
			if videoDetailResp.Type == "Episode" {
				if videoDetailResp.SeriesId != "" {
					jellyfinId = videoDetailResp.SeriesId
				}
				if videoDetailResp.SeasonId != "" {
					jellyfinId = videoDetailResp.SeasonId
				}
			}
			video, err := t.videoRepo.GetByJellyfinId(ctx, jellyfinId)
			if err != nil {
				log.Errorf("查询video失败: %v", err)
				return err
			}
			err = gen.Use(t.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
				if video == nil {
					title := ""
					if videoDetailResp.Type == "Movie" {
						title = videoDetailResp.Name
					}
					if videoDetailResp.Type == "Episode" {
						if videoDetailResp.SeriesId != "" {
							title = videoDetailResp.SeriesName
						}
						if videoDetailResp.SeasonId != "" {
							title = videoDetailResp.SeasonName
						}
					}
					overview := ""
					if videoDetailResp.Type == "Movie" {
						overview = videoDetailResp.Overview
					}
					if videoDetailResp.Type == "Episode" {
						overview = videoDetailResp.Overview
						if videoDetailResp.SeasonId != "" {
							overview = season.Overview
						}
					}
					genres := make([]string, 0)
					if videoDetailResp.Type == "Movie" {
						genres = videoDetailResp.Genres
					}
					if videoDetailResp.Type == "Episode" {
						genres = videoDetailResp.Genres
						if videoDetailResp.SeasonId != "" {
							genres = season.Genres
						}
					}

					video = &model.Video{
						Title:              title,
						VideoType:          videoDetailResp.Type,
						VoteRate:           lo.ToPtr(float32(math.Round(videoDetailResp.GoodRating*10) / 10)),
						Region:             lo.ToPtr(strings.Join(t.replaceRegions(ctx, videoDetailResp.Regions), ",")),
						Description:        lo.ToPtr(overview),
						PublishDay:         lo.ToPtr(time_util.FormatStrToYYYYMMDD(videoDetailResp.PremiereDate)),
						Thumbnail:          lo.ToPtr(strings.Join([]string{videoDetailResp.Id, videoDetailResp.SeasonId, videoDetailResp.SeriesId}, ",")),
						Genres:             lo.ToPtr(strings.Join(t.replaceGenres(ctx, genres), ",")),
						JellyfinID:         jellyfinId,
						JellyfinCreateTime: time_util.ParseUtcTime(videoDetailResp.DateCreated),
					}
					err := t.videoRepo.Create(ctx, tx, video)
					if err != nil {
						log.Errorf("创建video失败: %v", err)
						return err
					}
				}
				existEpisode[0].VideoId = video.ID
				err = t.episodeRepo.Updates(ctx, tx, existEpisode[0])
				if err != nil {
					log.Errorf("更新episode失败: %v", err)
					return err
				}
				//演员处理
				director, _ := lo.Find(videoDetailResp.Characters, func(item *jellyfindto.People) bool {
					return item.Type == "Director"
				})
				if director != nil {
					actor := &model.Actor{
						Name: director.Name,
					}
					err := t.actorRepo.Create(ctx, tx, actor)
					if err != nil {
						return err
					}
					err = t.videoActorMappingRepo.Create(ctx, tx, &model.VideoActorMapping{
						VideoID:    video.ID,
						ActorID:    actor.ID,
						Character:  lo.ToPtr(director.Role),
						IsDirector: true,
					})
					if err != nil {
						return err
					}
				}
				actors := lo.Filter(videoDetailResp.Characters, func(item *jellyfindto.People, index int) bool {
					return item.Type == "Actor"
				})
				for i, character := range actors {
					if i > 5 {
						break
					}
					actor := &model.Actor{
						Name: character.Name,
					}
					err = t.actorRepo.Create(ctx, tx, actor)
					if err != nil {
						return err
					}
					err = t.videoActorMappingRepo.Create(ctx, tx, &model.VideoActorMapping{
						VideoID:    video.ID,
						ActorID:    actor.ID,
						Character:  lo.ToPtr(character.Role),
						IsDirector: false,
					})
					if err != nil {
						return err
					}
				}
				//字幕处理
				result, err := httpclient_util.DoPost[any, jellyfindto.PlaybackInfo](ctx, domain+fmt.Sprintf(constant.JellyfinPlayInfo, videoDetailResp.Id, jellyfinDefaultUserId), nil, headerMap)
				if err != nil {
					return err
				}
				if result == nil {
					return nil
				}
				for _, mediaSource := range result.MediaSources {
					for _, mediaStream := range mediaSource.MediaStreams {
						if mediaStream.Type != "Subtitle" {
							continue
						}
						err := t.episodeSubtitleMappingRepo.Delete(ctx, tx, existEpisode[0].ID)
						if err != nil {
							return err
						}
						err = t.episodeSubtitleMappingRepo.Create(ctx, tx, &model.EpisodeSubtitleMapping{
							EpisodeID: existEpisode[0].ID,
							URL:       mediaStream.DeliveryUrl,
							Title:     mediaStream.DisplayTitle,
							Language:  mediaStream.Language,
							MimeType:  mediaStream.Codec,
						})
						if err != nil {
							return err
						}
					}
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
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
