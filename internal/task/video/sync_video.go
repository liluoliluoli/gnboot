package video

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/region_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/jellyfindto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type JfVideoTask struct {
	episodeRepo                *repo.EpisodeRepo
	videoRepo                  *repo.VideoRepo
	client                     redis.UniversalClient
	c                          *conf.Bootstrap
	configRepo                 *repo.ConfigRepo
	actorRepo                  *repo.ActorRepo
	videoActorMappingRepo      *repo.VideoActorMappingRepo
	episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
}

func NewJfVideoTask(episodeRepo *repo.EpisodeRepo, videoRepo *repo.VideoRepo, client redis.UniversalClient, c *conf.Bootstrap, configRepo *repo.ConfigRepo, actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo, episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo) *JfVideoTask {
	return &JfVideoTask{
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

func (t *JfVideoTask) Process(task *sdomain.Task) error {
	ctx := task.Ctx
	err := t.DoProcess(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *JfVideoTask) DoProcess(ctx context.Context) error {
	jellyfinBoxIpStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_JellyfinBoxIp)
	if err != nil {
		return err
	}
	jellyfinBoxIps := strings.Split(jellyfinBoxIpStr, ",")
	scanPathStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_JellyfinVideoSyncCategory)
	if err != nil {
		return err
	}
	scanPathIds := strings.Split(scanPathStr, ",")
	for _, scanPathId := range scanPathIds {
		lastSyncTime, err := t.videoRepo.QueryLastSyncTimeByJfId(ctx, scanPathId)
		if err != nil {
			return err
		}
		log.Infof("开始同步jellyfin刮削: syncURL=%s, pathid=%s, 上次同步时间=%v", jellyfinBoxIps[0], scanPathId, lastSyncTime)
		err = t.deepLoopListJellyfinPath(ctx, jellyfinBoxIps[0], scanPathId, lastSyncTime, scanPathId)
		if err != nil {
			log.Errorf("同步jellyfin刮削失败: %v", err)
			return err
		}
		log.Infof("结束同步jellyfin刮削: syncURL=%s, pathid=%s, type=%s", jellyfinBoxIps[0], scanPathId)
	}
	return nil
}

func (t *JfVideoTask) deepLoopListJellyfinPath(ctx context.Context, domain, parentId string, lastSyncTime *time.Time, rootPathId string) error {
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
		for index, content := range videoListResp.Items {
			err = t.deepLoopDetailJellyfinPath(ctx, domain, parentId, content.Id, content.Type, nil, nil, jellyfinDefaultUserId, headerMap, index, lastSyncTime, rootPathId)
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

func (t *JfVideoTask) deepLoopDetailJellyfinPath(ctx context.Context, domain, parentId string, id string, mediaType string, season *jellyfindto.VideoDetailResp, series *jellyfindto.VideoDetailResp, jellyfinDefaultUserId string, headerMap map[string]string, index int, lastSyncTime *time.Time, rootPathId string) error {
	syncDetailURL := fmt.Sprintf(domain+constant.JellyfinVideoDetail, jellyfinDefaultUserId, id)
	videoDetailResp, err := httpclient_util.DoGet[jellyfindto.VideoDetailResp](ctx, syncDetailURL, headerMap)
	if err != nil {
		return fmt.Errorf("jellyfin videoDetailResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
	}
	if videoDetailResp == nil {
		return fmt.Errorf("jellyfin videoDetailResp parentId %s file %s 返回结果无效", parentId, id)
	}
	if videoDetailResp.IsFolder {
		if mediaType == constant.JfFolder {
			return t.deepLoopListJellyfinPath(ctx, domain, videoDetailResp.Id, lastSyncTime, rootPathId)
		}
		if mediaType == constant.JfSeries {
			syncSeasonListURL := fmt.Sprintf(domain+constant.JellyfinSeaonsList, id)
			seasonListResp, err := httpclient_util.DoGet[jellyfindto.SeasonListResp](ctx, syncSeasonListURL, headerMap)
			if err != nil {
				return fmt.Errorf("jellyfin seasonListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
			}
			if seasonListResp == nil {
				return fmt.Errorf("jellyfin seasonListResp parentId %s file %s 返回结果无效", parentId, id)
			}
			for newIndex, content := range seasonListResp.Items {
				err := t.deepLoopDetailJellyfinPath(ctx, domain, id, content.Id, "Season", nil, videoDetailResp, jellyfinDefaultUserId, headerMap, newIndex, lastSyncTime, rootPathId)
				if err != nil {
					return err
				}
			}
		}
		if mediaType == constant.JfSeason {
			syncEpisodeListURL := fmt.Sprintf(domain+constant.JellyfinEpisodesList, parentId, id)
			episodeListResp, err := httpclient_util.DoGet[jellyfindto.EpisodeListResp](ctx, syncEpisodeListURL, headerMap)
			if err != nil {
				return fmt.Errorf("jellyfin episodeListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
			}
			if episodeListResp == nil {
				return fmt.Errorf("jellyfin episodeListResp parentId %s file %s 返回结果无效", parentId, id)
			}
			for newIndex, content := range episodeListResp.Items {
				err := t.deepLoopDetailJellyfinPath(ctx, domain, id, content.Id, "Episode", videoDetailResp, series, jellyfinDefaultUserId, headerMap, newIndex, lastSyncTime, rootPathId)
				if err != nil {
					return err
				}
			}
		}
	} else {
		_, valid := lo.Find(constant.SupportVideoTypes, func(item string) bool {
			if strings.HasSuffix(videoDetailResp.MediaSources[0].Path, item) {
				return true
			}
			return false
		})
		if !valid {
			return nil
		}
		dateCreated := time_util.ParseUtcTime(videoDetailResp.DateCreated)
		if lastSyncTime != nil && dateCreated.Before(lo.FromPtr(lastSyncTime)) {
			log.Infof("最后更新时间小于库记录，跳过处理【video】: title=%s", videoDetailResp.Name)
			return nil
		}

		jellyfinId := ""
		if videoDetailResp.Type == constant.JfMovie {
			jellyfinId = fmt.Sprintf("%s-%s", videoDetailResp.Id, constant.JfMovie)
		}
		if videoDetailResp.Type == constant.JfEpisode {
			if videoDetailResp.SeriesId != "" {
				jellyfinId = fmt.Sprintf("%s|%s", videoDetailResp.SeriesId, constant.JfSeries)
			}
			if videoDetailResp.SeasonId != "" {
				jellyfinId = fmt.Sprintf("%s|%s", videoDetailResp.SeasonId, constant.JfSeason)
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
						title = videoDetailResp.SeriesName + "：" + videoDetailResp.SeasonName
					}
				}
				overview := ""
				if videoDetailResp.Type == "Movie" {
					overview = videoDetailResp.Overview
				}
				if videoDetailResp.Type == "Episode" {
					overview = videoDetailResp.Overview
					if videoDetailResp.SeasonId != "" && season.Overview != "" {
						overview = season.Overview
					}
				}
				genres := make([]string, 0)
				if videoDetailResp.Type == "Movie" {
					genres = videoDetailResp.Genres
				}
				if videoDetailResp.Type == "Episode" {
					genres = videoDetailResp.Genres
					if series != nil && len(series.Genres) > 0 {
						genres = series.Genres
					}
					if videoDetailResp.SeasonId != "" && len(season.Genres) != 0 {
						genres = season.Genres
					}
				}
				var goodRatting float64
				if videoDetailResp.Type == "Movie" {
					goodRatting = videoDetailResp.GoodRating
				}
				if videoDetailResp.Type == "Episode" {
					goodRatting = videoDetailResp.GoodRating
					if series != nil && series.GoodRating != 0 {
						goodRatting = series.GoodRating
					}
					if videoDetailResp.SeasonId != "" && season.GoodRating != 0 {
						goodRatting = season.GoodRating
					}
				}
				var totalEpisode int32
				if videoDetailResp.Type == "Movie" {
					totalEpisode = 1
				}
				if videoDetailResp.Type == "Episode" {
					if series != nil {
						totalEpisode = series.ChildCount
					}
					if videoDetailResp.SeasonId != "" && season.ChildCount != 0 {
						totalEpisode = season.ChildCount
					}
				}
				//国家处理
				var externalUrls []*jellyfindto.ExternalUrl
				if videoDetailResp.Type == "Movie" {
					externalUrls = videoDetailResp.ExternalUrls
				}
				if videoDetailResp.Type == "Episode" {
					if series != nil {
						externalUrls = append(externalUrls, series.ExternalUrls...)
					}
					if videoDetailResp.SeasonId != "" && season.ChildCount != 0 {
						externalUrls = append(externalUrls, season.ExternalUrls...)
					}
				}
				region, externalUrls := t.getRegion(ctx, externalUrls)
				externalUrlStr := ""
				if len(externalUrls) > 0 {
					externalUrlStr, _ = json_util.MarshalString(externalUrls)
				}

				video = &model.Video{
					Title:              strings.TrimSpace(title),
					VideoType:          t.replaceVideType(ctx, videoDetailResp.MediaSources[0].Path),
					VoteRate:           lo.ToPtr(float32(math.Round(goodRatting*10) / 10)),
					Region:             lo.ToPtr(region),
					Description:        lo.ToPtr(strings.TrimSpace(overview)),
					PublishDay:         lo.ToPtr(time_util.FormatStrToYYYYMMDD(videoDetailResp.PremiereDate)),
					Thumbnail:          lo.ToPtr(t.getValidThumbnail(ctx, domain, []string{videoDetailResp.Id, videoDetailResp.SeasonId, videoDetailResp.SeriesId})),
					Genres:             lo.ToPtr(strings.Join(t.replaceGenres(ctx, genres), ",")),
					JellyfinID:         jellyfinId,
					JellyfinCreateTime: dateCreated,
					TotalEpisode:       lo.ToPtr(totalEpisode),
					Ext:                lo.ToPtr(externalUrlStr),
					WatchCount:         0,
					JellyfinRootPathID: lo.ToPtr(rootPathId),
				}
				err = t.videoRepo.Create(ctx, tx, video)
				if err != nil {
					log.Errorf("创建video失败: %v", err)
					return err
				}
				log.Infof("成功写入【video】: title=%s", title)
			}
			filePath, fileName := splitJellyfinPath(videoDetailResp.MediaSources[0].Path)
			episode, err := t.episodeRepo.QueryByPathAndName(ctx, filePath, fileName)
			if err != nil {
				return err
			}
			if episode == nil {
				episode = &model.Episode{
					XiaoyaPath:   &filePath,
					EpisodeTitle: fileName,
					Platform:     lo.ToPtr(constant.Platform),
					IsValid:      true,
					CreateTime:   time.Now(),
					UpdateTime:   time.Now(),
					VideoID:      video.ID,
					Episode:      int32(index + 1),
					JellyfinID:   lo.ToPtr(videoDetailResp.Id),
					DisplayTitle: lo.ToPtr(lo.Ternary(videoDetailResp.Name != "", videoDetailResp.Name, fileName)),
				}
				if err := t.episodeRepo.Create(ctx, nil, episode); err != nil {
					log.Errorf("写入episode失败（path=%s, title=%s）: %v", filePath, fileName, err)
				} else {
					log.Infof("成功写入【episode】: path=%s, title=%s", filePath, fileName)
				}
			}
			//演员处理
			characters := make([]*jellyfindto.People, 0)
			if videoDetailResp.Type == "Movie" {
				characters = videoDetailResp.Characters
			}
			if videoDetailResp.Type == "Episode" {
				characters = videoDetailResp.Characters
				if series != nil && len(series.Characters) > 0 {
					characters = series.Characters
				}
				if videoDetailResp.SeasonId != "" && len(season.Characters) != 0 {
					characters = season.Characters
				}
			}
			videoActorMappings := make([]*model.VideoActorMapping, 0)
			characterGroup := lo.GroupBy(characters, func(item *jellyfindto.People) string {
				return item.Type
			})
			for _, character := range characterGroup["Director"] {
				actor := &model.Actor{
					Name:      character.Name,
					Thumbnail: lo.ToPtr(fmt.Sprintf(constant.PrimaryThumbnail, character.Id)),
				}
				err := t.actorRepo.Create(ctx, tx, actor)
				if err != nil {
					return err
				}
				videoActorMappings = append(videoActorMappings, &model.VideoActorMapping{
					VideoID:    video.ID,
					ActorID:    actor.ID,
					Character:  lo.ToPtr(character.Role),
					IsDirector: true,
				})
				break
			}
			for i, character := range characterGroup["Actor"] {
				if i > 5 {
					break
				}
				actor := &model.Actor{
					Name:      character.Name,
					Thumbnail: lo.ToPtr(fmt.Sprintf(constant.PrimaryThumbnail, character.Id)),
				}
				err := t.actorRepo.Create(ctx, tx, actor)
				if err != nil {
					return err
				}
				videoActorMappings = append(videoActorMappings, &model.VideoActorMapping{
					VideoID:    video.ID,
					ActorID:    actor.ID,
					Character:  lo.ToPtr(character.Role),
					IsDirector: false,
				})
				break
			}
			err = t.videoActorMappingRepo.BatchCreate(ctx, tx, videoActorMappings)
			if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
				return err
			}
			log.Infof("成功写入【videoActorMappingRepo】")
			//字幕处理
			result, err := httpclient_util.DoGet[jellyfindto.PlaybackInfo](ctx, domain+fmt.Sprintf(constant.JellyfinPlayInfo, videoDetailResp.Id, jellyfinDefaultUserId), headerMap)
			if err != nil {
				return err
			}
			if result == nil {
				return nil
			}
			err = t.episodeSubtitleMappingRepo.Delete(ctx, tx, episode.ID)
			if err != nil {
				return err
			}
			for _, mediaSource := range result.MediaSources {
				for _, mediaStream := range mediaSource.MediaStreams {
					if mediaStream.Type != "Subtitle" || mediaStream.DeliveryUrl == "" {
						continue
					}
					err = t.episodeSubtitleMappingRepo.Create(ctx, tx, &model.EpisodeSubtitleMapping{
						EpisodeID: episode.ID,
						URL:       mediaStream.DeliveryUrl,
						Title:     mediaStream.DisplayTitle,
						Language:  mediaStream.Language,
						MimeType:  mediaStream.Codec,
					})
					if err != nil {
						return err
					}
					log.Infof("成功写入【subtitle】: url=%s", mediaStream.DeliveryUrl)
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// splitJellyfinPath http://xiaoya.host:5678/d/电影/4K系列/Marvel漫威系列/漫威宇宙/蜘蛛侠：英雄无归.Spider-Man.No.Way.Home.2021.UHD.BluRay.2160p.x265.10bit.DoVi.2Audios.mUHD-FRDS.mkv
func splitJellyfinPath(url string) (string, string) {
	var path string
	if strings.HasPrefix(url, "/") {
		pathStart := 0
		pathEnd := strings.LastIndex(url, "/")
		path = url[pathStart:pathEnd]
	} else {
		domain := strings.Split(url, "/")[2] // "xiaoya.host:5678"
		pathStart := strings.Index(url, domain+"/d") + len(domain+"/d")
		pathEnd := strings.LastIndex(url, "/")
		path = url[pathStart:pathEnd] // "/电影/4K系列/Marvel漫威系列/漫威宇宙"
	}
	filename := filepath.Base(url) // "蜘蛛侠：英雄无归.Spider-Man.No.Way.Home.2021.UHD.BluRay.2160p.x265.10bit.DoVi.2Audios.mUHD-FRDS.mkv"
	return path, filename
}

func (t *JfVideoTask) replaceRegions(ctx context.Context, regions []string) []string {
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

func (t *JfVideoTask) replaceGenres(ctx context.Context, genres []string) []string {
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

func (t *JfVideoTask) replaceVideType(ctx context.Context, path string) string {
	configMap, _ := t.configRepo.GetConfigMapByKey(ctx, constant.Key_PathVideoTypeMapping)
	if configMap != nil {
		for k, v := range configMap {
			if strings.Contains(path, k) {
				return v
			}
		}
	}
	return "Unknown"
}

func (t *JfVideoTask) getValidThumbnail(ctx context.Context, domain string, ids []string) string {
	for _, id := range ids {
		if id == "" {
			continue
		}
		url := fmt.Sprintf(constant.PrimaryThumbnail, id)
		valid, err := httpclient_util.CheckImageUrl(domain + url)
		if err != nil {
			continue
		}
		if valid {
			return url
		}
	}
	return constant.DefaultThumbnail
}

func (t *JfVideoTask) getRegion(ctx context.Context, externalUrls []*jellyfindto.ExternalUrl) (string, []*jellyfindto.ExternalUrl) {
	if len(externalUrls) == 0 {
		return "", externalUrls
	}
	uniqExternalUrls := lo.UniqBy(externalUrls, func(item *jellyfindto.ExternalUrl) string {
		return item.Name
	})
	sort.Slice(uniqExternalUrls, func(i, j int) bool {
		if constant.SortMap[uniqExternalUrls[i].Name] == 0 || constant.SortMap[uniqExternalUrls[j].Name] == 0 {
			return true
		}
		return constant.SortMap[uniqExternalUrls[i].Name] < constant.SortMap[uniqExternalUrls[j].Name]
	})
	for _, externalUrl := range uniqExternalUrls {
		html, err := httpclient_util.DoHtml(ctx, externalUrl.Url)
		if err != nil {
			continue
		}
		regular := constant.RegularMap[externalUrl.Name]
		if regular != nil {
			matches := regular.FindStringSubmatch(html)
			if len(matches) > 1 {
				region := region_util.GetCnNameByName(ctx, matches[1])
				if region != "" {
					externalUrl.Used = true
					return region, uniqExternalUrls
				}
			}
		}
	}
	return "", uniqExternalUrls
}
