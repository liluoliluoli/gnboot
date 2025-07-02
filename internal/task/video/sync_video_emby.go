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
	"github.com/liluoliluoli/gnboot/internal/integration/dto/embydto"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/xiaoyadto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type EmbyVideoTask struct {
	episodeRepo                *repo.EpisodeRepo
	videoRepo                  *repo.VideoRepo
	client                     redis.UniversalClient
	c                          *conf.Bootstrap
	configRepo                 *repo.ConfigRepo
	actorRepo                  *repo.ActorRepo
	videoActorMappingRepo      *repo.VideoActorMappingRepo
	episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo
}

func NewEmbyVideoTask(episodeRepo *repo.EpisodeRepo, videoRepo *repo.VideoRepo, client redis.UniversalClient, c *conf.Bootstrap, configRepo *repo.ConfigRepo, actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo, episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo) *EmbyVideoTask {
	return &EmbyVideoTask{
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

func (t *EmbyVideoTask) Process(task *sdomain.Task) error {
	ctx := task.Ctx
	err := t.FullSync(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *EmbyVideoTask) FullSync(ctx context.Context) error {
	boxIpStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_EmbyBoxIp)
	if err != nil {
		return err
	}
	boxIps := strings.Split(boxIpStr, ",")
	scanPathStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyVideoSyncCategory)
	if err != nil {
		return err
	}
	scanPathIds := strings.Split(scanPathStr, ",")
	for _, scanPathId := range scanPathIds {
		log.Infof("开始全量同步emby刮削: syncURL=%s, pathid=%s", boxIps[0], scanPathId)
		err = t.deepLoopListEmbyPath(ctx, boxIps[0], scanPathId, scanPathId)
		if err != nil {
			log.Errorf("全量同步emby刮削失败: %v", err)
			return err
		}
		log.Infof("结束全量同步emby刮削: syncURL=%s, pathid=%s", boxIps[0], scanPathId)
	}
	return nil
}

func (t *EmbyVideoTask) deepLoopListEmbyPath(ctx context.Context, domain, parentId string, rootPathId string) error {
	startIndex := int32(0)
	embyDefaultToken, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyDefaultToken)
	if err != nil {
		return err
	}
	embyDefaultUserId, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyDefaultUserId)
	if err != nil {
		return err
	}
	for {
		headerMap := make(map[string]string)
		headerMap["X-Emby-Token"] = embyDefaultToken
		syncListURL := fmt.Sprintf(domain+constant.EmbyVideoList, embyDefaultUserId, startIndex, parentId, constant.PageSize)
		videoListResp, err := httpclient_util.DoGet[embydto.VideoListResp](ctx, syncListURL, headerMap)
		if err != nil {
			return fmt.Errorf("emby请求parentId %s 游标 %d 开始的返回结果失败: %v", parentId, startIndex, err)
		}
		if videoListResp == nil {
			return fmt.Errorf("emby parentId %s 游标 %d 开始的返回结果无效", parentId, startIndex)
		}
		for index, content := range videoListResp.Items {
			err = t.deepLoopDetailEmbyPath(ctx, domain, parentId, content.Id, content.Type, nil, nil, embyDefaultUserId, headerMap, index, rootPathId)
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

func (t *EmbyVideoTask) deepLoopDetailEmbyPath(ctx context.Context, domain, parentId string, id string, mediaType string, season *embydto.VideoDetailResp, series *embydto.VideoDetailResp, embyDefaultUserId string, headerMap map[string]string, index int, rootPathId string) error {
	//查出详情
	syncDetailURL := fmt.Sprintf(domain+constant.EmbyVideoDetail, embyDefaultUserId, id)
	videoDetailResp, err := httpclient_util.DoGet[embydto.VideoDetailResp](ctx, syncDetailURL, headerMap)
	if err != nil {
		return fmt.Errorf("jellyfin videoDetailResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
	}
	if videoDetailResp == nil {
		return fmt.Errorf("jellyfin videoDetailResp parentId %s file %s 返回结果无效", parentId, id)
	}

	//按类型处理
	if mediaType == constant.JfFolder {
		return t.deepLoopListEmbyPath(ctx, domain, id, rootPathId)
	}
	if mediaType == constant.JfSeries {
		syncSeasonListURL := fmt.Sprintf(domain+constant.EmbySeaonsList, id, 10000)
		seasonListResp, err := httpclient_util.DoGet[embydto.SeasonListResp](ctx, syncSeasonListURL, headerMap)
		if err != nil {
			return fmt.Errorf("emby seasonListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
		}
		if seasonListResp == nil {
			return fmt.Errorf("emby seasonListResp parentId %s file %s 返回结果无效", parentId, id)
		}
		for newIndex, content := range seasonListResp.Items {
			err := t.deepLoopDetailEmbyPath(ctx, domain, id, content.Id, constant.JfSeason, nil, videoDetailResp, embyDefaultUserId, headerMap, newIndex, rootPathId)
			if err != nil {
				return err
			}
		}
	}
	if mediaType == constant.JfSeason {
		syncEpisodeListURL := fmt.Sprintf(domain+constant.EmbyEpisodesList, parentId, id, 10000)
		episodeListResp, err := httpclient_util.DoGet[embydto.EpisodeListResp](ctx, syncEpisodeListURL, headerMap)
		if err != nil {
			return fmt.Errorf("emby episodeListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
		}
		if episodeListResp == nil {
			return fmt.Errorf("emby episodeListResp parentId %s file %s 返回结果无效", parentId, id)
		}
		for newIndex, content := range episodeListResp.Items {
			err := t.deepLoopDetailEmbyPath(ctx, domain, id, content.Id, constant.JfEpisode, videoDetailResp, series, embyDefaultUserId, headerMap, newIndex, rootPathId)
			if err != nil {
				return err
			}
		}
	}

	if mediaType == constant.JfMovie || mediaType == constant.JfEpisode {
		_, valid := lo.Find(constant.SupportVideoTypes, func(item string) bool {
			if strings.HasSuffix(videoDetailResp.MediaSources[0].Path, item) {
				return true
			}
			return false
		})
		if !valid {
			return nil
		}

		filePath, fileName := splitEmbyPath(videoDetailResp.MediaSources[0].Path)
		episode, err := t.episodeRepo.QueryByPathAndName(ctx, filePath, fileName)
		if err != nil {
			return err
		}
		if episode != nil { //先前的任务已经处理过了
			return nil
		}

		//ratio处理
		sizeMap := make(map[string]int64)
		xiaoyaDomain, _ := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_XiaoYaBoxIp)
		if xiaoyaDomain != "" {
			xiaoyaListURL := fmt.Sprintf(xiaoyaDomain + constant.XiaoYaVideoList)
			xiaoyaVideoListResp, _ := httpclient_util.DoPost[xiaoyadto.VideoListReq, xiaoyadto.XiaoyaResult[xiaoyadto.VideoListResp]](ctx, xiaoyaListURL, &xiaoyadto.VideoListReq{
				Path:     filePath,
				Password: "",
				PerPage:  2000,
				Refresh:  false,
				Page:     1,
			}, nil)
			if xiaoyaVideoListResp != nil && xiaoyaVideoListResp.Data != nil && xiaoyaVideoListResp.Data.Content != nil {
				sizeMap = lo.SliceToMap(xiaoyaVideoListResp.Data.Content, func(item *xiaoyadto.VideoContent) (string, int64) {
					return item.Name, item.Size
				})
			}
		}
		episodeSize := sizeMap[fileName]
		if episodeSize == 0 { //xiaoya不存在该数据
			return nil
		}
		ratio := ""
		if videoDetailResp.Type == constant.JfMovie {
			ratio = constant.SD
			if episodeSize > constant.HD_MOVIE_MIN_SIZE {
				ratio = constant.HD
			}
			if episodeSize > constant.QHD_MOVIE_MIN_SIZE {
				ratio = constant.QHD
			}
		}
		if videoDetailResp.Type == constant.JfEpisode {
			ratio = constant.SD
			if episodeSize > constant.HD_EPISODE_MIN_SIZE {
				ratio = constant.HD
			}
		}

		embyId := ""
		if videoDetailResp.Type == constant.JfMovie {
			embyId = fmt.Sprintf("%s-%s", videoDetailResp.Id, constant.JfMovie)
		}
		if videoDetailResp.Type == constant.JfEpisode {
			if videoDetailResp.SeriesId != "" {
				embyId = fmt.Sprintf("%s|%s", videoDetailResp.SeriesId, constant.JfSeries)
			}
			if videoDetailResp.SeasonId != "" {
				embyId = fmt.Sprintf("%s|%s", videoDetailResp.SeasonId, constant.JfSeason)
			}
		}
		embyCreateTime := time_util.ParseUtcTime(videoDetailResp.DateCreated)
		video, err := t.videoRepo.GetByJellyfinId(ctx, embyId)
		if err != nil {
			log.Errorf("查询video失败: %v", err)
			return err
		}
		err = gen.Use(t.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
			if video == nil {
				title := ""
				if videoDetailResp.Type == constant.JfMovie {
					title = videoDetailResp.Name
				}
				if videoDetailResp.Type == constant.JfEpisode {
					if videoDetailResp.SeriesId != "" {
						title = videoDetailResp.SeriesName
					}
					if videoDetailResp.SeasonId != "" {
						title = videoDetailResp.SeriesName + "：" + videoDetailResp.SeasonName
					}
				}
				overview := ""
				if videoDetailResp.Type == constant.JfMovie {
					overview = videoDetailResp.Overview
				}
				if videoDetailResp.Type == constant.JfEpisode {
					overview = videoDetailResp.Overview
					if series != nil && series.Overview != "" {
						overview = series.Overview
					}
					if season != nil && season.Overview != "" {
						overview = season.Overview
					}
				}
				genres := make([]string, 0)
				if videoDetailResp.Type == constant.JfMovie {
					genres = videoDetailResp.Genres
				}
				if videoDetailResp.Type == constant.JfEpisode {
					genres = videoDetailResp.Genres
					if series != nil && len(series.Genres) > 0 {
						genres = series.Genres
					}
					if season != nil && len(season.Genres) != 0 {
						genres = season.Genres
					}
				}
				var goodRatting float64
				if videoDetailResp.Type == constant.JfMovie {
					goodRatting = videoDetailResp.GoodRating
				}
				if videoDetailResp.Type == constant.JfEpisode {
					goodRatting = videoDetailResp.GoodRating
					if series != nil && series.GoodRating != 0 {
						goodRatting = series.GoodRating
					}
					if season != nil && season.GoodRating != 0 {
						goodRatting = season.GoodRating
					}
				}
				var totalEpisode int32
				if videoDetailResp.Type == constant.JfMovie {
					totalEpisode = 1
				}
				if videoDetailResp.Type == constant.JfEpisode {
					if series != nil {
						totalEpisode = series.ChildCount
					}
					if season != nil && season.ChildCount != 0 {
						totalEpisode = season.ChildCount
					}
				}
				//国家处理
				var externalUrls []*embydto.ExternalUrl
				if videoDetailResp.Type == constant.JfMovie {
					externalUrls = videoDetailResp.ExternalUrls
				}
				if videoDetailResp.Type == constant.JfEpisode {
					if series != nil {
						externalUrls = append(externalUrls, series.ExternalUrls...)
					}
					if season != nil {
						externalUrls = append(externalUrls, season.ExternalUrls...)
					}
				}
				region, externalUrls := t.getRegion(ctx, externalUrls, videoDetailResp.MediaSources[0].Path)
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
					JellyfinID:         embyId,
					JellyfinCreateTime: embyCreateTime,
					TotalEpisode:       lo.ToPtr(totalEpisode),
					Ext:                lo.ToPtr(externalUrlStr),
					WatchCount:         0,
					JellyfinRootPathID: lo.ToPtr(rootPathId),
					Ratio:              lo.ToPtr(ratio),
				}
				err = t.videoRepo.Create(ctx, tx, video)
				if err != nil {
					log.Errorf("创建video失败: %v", err)
					return err
				}
				log.Infof("成功写入【video】: title=%s", title)
			}

			episode = &model.Episode{
				XiaoyaPath:    &filePath,
				EpisodeTitle:  fileName,
				Platform:      lo.ToPtr(constant.Platform),
				IsValid:       true,
				CreateTime:    time.Now(),
				UpdateTime:    time.Now(),
				VideoID:       video.ID,
				Episode:       int32(index + 1),
				JellyfinID:    lo.ToPtr(videoDetailResp.Id),
				DisplayTitle:  lo.ToPtr(lo.Ternary(videoDetailResp.Name != "", videoDetailResp.Name, fileName)),
				Size:          lo.ToPtr(episodeSize),
				Ratio:         lo.ToPtr(ratio),
				JfCreateTime:  lo.ToPtr(embyCreateTime),
				JfPublishTime: lo.ToPtr(time_util.ParseUtcTime(videoDetailResp.DateCreated)),
				JfRootPathID:  lo.ToPtr(rootPathId),
			}
			if err := t.episodeRepo.Create(ctx, nil, episode); err != nil {
				log.Errorf("写入episode失败（path=%s, title=%s）: %v", filePath, fileName, err)
			} else {
				log.Infof("成功写入【episode】: path=%s, title=%s", filePath, fileName)
			}
			//演员处理
			characters := make([]*embydto.People, 0)
			if videoDetailResp.Type == constant.JfMovie {
				characters = videoDetailResp.Characters
			}
			if videoDetailResp.Type == constant.JfEpisode {
				characters = videoDetailResp.Characters
				if series != nil && len(series.Characters) > 0 {
					characters = series.Characters
				}
				if season != nil && len(season.Characters) != 0 {
					characters = season.Characters
				}
			}
			videoActorMappings := make([]*model.VideoActorMapping, 0)
			characterGroup := lo.GroupBy(characters, func(item *embydto.People) string {
				return item.Type
			})
			for _, character := range characterGroup["Director"] {
				actor := &model.Actor{
					Name:      character.Name,
					Thumbnail: lo.ToPtr(fmt.Sprintf(constant.EmbyPrimaryThumbnail, character.Id)),
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
					Thumbnail: lo.ToPtr(fmt.Sprintf(constant.EmbyPrimaryThumbnail, character.Id)),
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
			}
			err = t.videoActorMappingRepo.BatchCreate(ctx, tx, videoActorMappings)
			if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
				return err
			}
			log.Infof("成功写入【videoActorMappingRepo】")
			//字幕处理
			result, err := httpclient_util.DoGet[embydto.PlaybackInfo](ctx, domain+fmt.Sprintf(constant.EmbyPlayInfo, videoDetailResp.Id), headerMap)
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
func splitEmbyPath(url string) (string, string) {
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

func (t *EmbyVideoTask) replaceRegions(ctx context.Context, regions []string) []string {
	regionMap, _ := t.configRepo.GetConfigMapByKey(ctx, constant.Key_RegionMapping)
	targets := lo.Map(regions, func(item string, index int) string {
		for source, target := range regionMap {
			if strings.Contains(item, source) {
				return target
			}
		}
		return ""
	})
	return lo.Filter(targets, func(item string, index int) bool {
		return item != ""
	})
}

func (t *EmbyVideoTask) replaceGenres(ctx context.Context, genres []string) []string {
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

func (t *EmbyVideoTask) replaceVideType(ctx context.Context, path string) string {
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

func (t *EmbyVideoTask) getValidThumbnail(ctx context.Context, domain string, ids []string) string {
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

func (t *EmbyVideoTask) getRegion(ctx context.Context, externalUrls []*embydto.ExternalUrl, xiaoyaPath string) (string, []*embydto.ExternalUrl) {
	newRegions := t.replaceRegions(ctx, []string{xiaoyaPath})
	if len(newRegions) != 0 {
		return newRegions[0], externalUrls
	}
	if len(externalUrls) == 0 {
		return "", externalUrls
	}
	uniqExternalUrls := lo.UniqBy(externalUrls, func(item *embydto.ExternalUrl) string {
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
					externalUrl.UsedCountry = true
					return region, uniqExternalUrls
				}
			}
		}
	}
	return "", uniqExternalUrls
}
