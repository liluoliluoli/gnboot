package video

import (
	"context"
	"errors"
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
	"net/url"
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
	PalyInfoReq                *embydto.PlaybackInfoReq
}

func NewEmbyVideoTask(episodeRepo *repo.EpisodeRepo, videoRepo *repo.VideoRepo, client redis.UniversalClient, c *conf.Bootstrap, configRepo *repo.ConfigRepo, actorRepo *repo.ActorRepo, videoActorMappingRepo *repo.VideoActorMappingRepo, episodeSubtitleMappingRepo *repo.EpisodeSubtitleMappingRepo) *EmbyVideoTask {
	task := &EmbyVideoTask{
		episodeRepo:                episodeRepo,
		videoRepo:                  videoRepo,
		client:                     client,
		c:                          c,
		configRepo:                 configRepo,
		actorRepo:                  actorRepo,
		videoActorMappingRepo:      videoActorMappingRepo,
		episodeSubtitleMappingRepo: episodeSubtitleMappingRepo,
	}
	embyPlayinfoReq, err := json_util.Unmarshal[*embydto.PlaybackInfoReq]("{\"DeviceProfile\":{\"MaxStaticBitrate\":140000000,\"MaxStreamingBitrate\":140000000,\"MusicStreamingTranscodingBitrate\":192000,\"DirectPlayProfiles\":[{\"Container\":\"mp4,m4v\",\"Type\":\"Video\",\"VideoCodec\":\"h264,hevc,av1,vp8,vp9\",\"AudioCodec\":\"mp3,aac,opus,flac,vorbis\"},{\"Container\":\"mkv\",\"Type\":\"Video\",\"VideoCodec\":\"h264,hevc,av1,vp8,vp9\",\"AudioCodec\":\"mp3,aac,opus,flac,vorbis\"},{\"Container\":\"flv\",\"Type\":\"Video\",\"VideoCodec\":\"h264\",\"AudioCodec\":\"aac,mp3\"},{\"Container\":\"3gp\",\"Type\":\"Video\",\"VideoCodec\":\"\",\"AudioCodec\":\"mp3,aac,opus,flac,vorbis\"},{\"Container\":\"mov\",\"Type\":\"Video\",\"VideoCodec\":\"h264\",\"AudioCodec\":\"mp3,aac,opus,flac,vorbis\"},{\"Container\":\"opus\",\"Type\":\"Audio\"},{\"Container\":\"mp3\",\"Type\":\"Audio\",\"AudioCodec\":\"mp3\"},{\"Container\":\"mp2,mp3\",\"Type\":\"Audio\",\"AudioCodec\":\"mp2\"},{\"Container\":\"m4a\",\"AudioCodec\":\"aac\",\"Type\":\"Audio\"},{\"Container\":\"mp4\",\"AudioCodec\":\"aac\",\"Type\":\"Audio\"},{\"Container\":\"flac\",\"Type\":\"Audio\"},{\"Container\":\"webma,webm\",\"Type\":\"Audio\"},{\"Container\":\"wav\",\"Type\":\"Audio\",\"AudioCodec\":\"PCM_S16LE,PCM_S24LE\"},{\"Container\":\"ogg\",\"Type\":\"Audio\"},{\"Container\":\"webm\",\"Type\":\"Video\",\"AudioCodec\":\"vorbis,opus\",\"VideoCodec\":\"av1,VP8,VP9\"}],\"TranscodingProfiles\":[{\"Container\":\"aac\",\"Type\":\"Audio\",\"AudioCodec\":\"aac\",\"Context\":\"Streaming\",\"Protocol\":\"hls\",\"MaxAudioChannels\":\"2\",\"MinSegments\":\"1\",\"BreakOnNonKeyFrames\":true},{\"Container\":\"aac\",\"Type\":\"Audio\",\"AudioCodec\":\"aac\",\"Context\":\"Streaming\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"mp3\",\"Type\":\"Audio\",\"AudioCodec\":\"mp3\",\"Context\":\"Streaming\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"opus\",\"Type\":\"Audio\",\"AudioCodec\":\"opus\",\"Context\":\"Streaming\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"wav\",\"Type\":\"Audio\",\"AudioCodec\":\"wav\",\"Context\":\"Streaming\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"opus\",\"Type\":\"Audio\",\"AudioCodec\":\"opus\",\"Context\":\"Static\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"mp3\",\"Type\":\"Audio\",\"AudioCodec\":\"mp3\",\"Context\":\"Static\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"aac\",\"Type\":\"Audio\",\"AudioCodec\":\"aac\",\"Context\":\"Static\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"wav\",\"Type\":\"Audio\",\"AudioCodec\":\"wav\",\"Context\":\"Static\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"mkv\",\"Type\":\"Video\",\"AudioCodec\":\"mp3,aac,opus,flac,vorbis\",\"VideoCodec\":\"h264,hevc,av1,vp8,vp9\",\"Context\":\"Static\",\"MaxAudioChannels\":\"2\",\"CopyTimestamps\":true},{\"Container\":\"ts\",\"Type\":\"Video\",\"AudioCodec\":\"mp3,aac\",\"VideoCodec\":\"hevc,h264,av1\",\"Context\":\"Streaming\",\"Protocol\":\"hls\",\"MaxAudioChannels\":\"2\",\"MinSegments\":\"1\",\"BreakOnNonKeyFrames\":true,\"ManifestSubtitles\":\"vtt\"},{\"Container\":\"webm\",\"Type\":\"Video\",\"AudioCodec\":\"vorbis\",\"VideoCodec\":\"vpx\",\"Context\":\"Streaming\",\"Protocol\":\"http\",\"MaxAudioChannels\":\"2\"},{\"Container\":\"mp4\",\"Type\":\"Video\",\"AudioCodec\":\"mp3,aac,opus,flac,vorbis\",\"VideoCodec\":\"h264\",\"Context\":\"Static\",\"Protocol\":\"http\"}],\"ContainerProfiles\":[],\"CodecProfiles\":[{\"Type\":\"VideoAudio\",\"Codec\":\"aac\",\"Conditions\":[{\"Condition\":\"Equals\",\"Property\":\"IsSecondaryAudio\",\"Value\":\"false\",\"IsRequired\":\"false\"}]},{\"Type\":\"VideoAudio\",\"Conditions\":[{\"Condition\":\"Equals\",\"Property\":\"IsSecondaryAudio\",\"Value\":\"false\",\"IsRequired\":\"false\"}]},{\"Type\":\"Video\",\"Codec\":\"h264\",\"Conditions\":[{\"Condition\":\"EqualsAny\",\"Property\":\"VideoProfile\",\"Value\":\"high|main|baseline|constrained baseline|high 10\",\"IsRequired\":false},{\"Condition\":\"LessThanEqual\",\"Property\":\"VideoLevel\",\"Value\":\"62\",\"IsRequired\":false}]},{\"Type\":\"Video\",\"Codec\":\"hevc\",\"Conditions\":[{\"Condition\":\"EqualsAny\",\"Property\":\"VideoCodecTag\",\"Value\":\"hvc1|hev1|hevc|hdmv\",\"IsRequired\":false}]}],\"SubtitleProfiles\":[{\"Format\":\"vtt\",\"Method\":\"Hls\"},{\"Format\":\"eia_608\",\"Method\":\"VideoSideData\",\"Protocol\":\"hls\"},{\"Format\":\"eia_708\",\"Method\":\"VideoSideData\",\"Protocol\":\"hls\"},{\"Format\":\"vtt\",\"Method\":\"External\"},{\"Format\":\"ass\",\"Method\":\"External\"},{\"Format\":\"ssa\",\"Method\":\"External\"}],\"ResponseProfiles\":[{\"Type\":\"Video\",\"Container\":\"m4v\",\"MimeType\":\"video/mp4\"}]}}")
	if err != nil {

	}
	task.PalyInfoReq = embyPlayinfoReq
	return task
}

func (t *EmbyVideoTask) Process(task *sdomain.Task) error {
	ctx := task.Ctx
	err := t.LatestSync(ctx, "", 200)
	if err != nil {
		return err
	}
	return nil
}

func (t *EmbyVideoTask) LatestSync(ctx context.Context, scanPathStr string, findLatestCount int32) error {
	ctx = context.Background()
	boxIpStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_EmbyBoxIp)
	if err != nil {
		return err
	}
	boxIps := strings.Split(boxIpStr, ",")
	var scanPaths = make([]string, 0)
	if scanPathStr != "" {
		scanPaths = strings.Split(scanPathStr, "|")
	} else {
		scanPathStr, err = t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyVideoSyncCategory)
		if err != nil {
			return err
		}
		scanPaths = strings.Split(scanPathStr, "|")
	}

	for _, scanPath := range scanPaths {
		log.Warnf("开始增量同步emby刮削: pathid=%s", scanPath)
		scanPathAndType := strings.Split(scanPath, ":")
		err = t.deepLoopLatestListEmbyPath(ctx, boxIps[0], scanPathAndType[0], "", scanPathAndType[0], lo.TernaryF(len(scanPathAndType) >= 2, func() string {
			return scanPathAndType[1]
		}, func() string {
			return ""
		}), lo.TernaryF(len(scanPathAndType) >= 3, func() string {
			return scanPathAndType[2]
		}, func() string {
			return "false"
		}), findLatestCount)
		if err != nil {
			log.Errorf("增量同步emby刮削失败: %v", err)
			return err
		}
		log.Warnf("结束增量同步emby刮削: pathid=%s", scanPathAndType[0])
	}
	return nil
}

func (t *EmbyVideoTask) deepLoopLatestListEmbyPath(ctx context.Context, domain, parentId string, parentName string, rootPathId string, includeItemType string, recursive string, findLatestCount int32) error {
	startIndex := int32(0)
	embyDefaultToken, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyDefaultToken)
	if err != nil {
		return err
	}
	embyDefaultUserId, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyDefaultUserId)
	if err != nil {
		return err
	}
	headerMap := make(map[string]string)
	headerMap["X-Emby-Token"] = embyDefaultToken
	syncListURL := fmt.Sprintf(domain+constant.EmbyLatestVideoList, embyDefaultUserId, parentId, findLatestCount)
	if includeItemType != "" {
		syncListURL = syncListURL + "&IncludeItemTypes=" + includeItemType
	}
	videoItemsResp, err := httpclient_util.DoGet[[]*embydto.VideoItem](ctx, syncListURL, headerMap)
	if err != nil {
		return fmt.Errorf("emby请求parentId %s 游标 %d 开始的返回结果失败: %v", parentId, startIndex, err)
	}
	if videoItemsResp == nil {
		return fmt.Errorf("emby parentId %s 游标 %d 开始的返回结果无效", parentId, startIndex)
	}
	for index, content := range lo.FromPtr(videoItemsResp) {
		err = t.deepLoopDetailEmbyPath(ctx, domain, parentId, parentName, len(lo.FromPtr(videoItemsResp)), content.Id, content.Type, nil, nil, embyDefaultUserId, headerMap, index, rootPathId, recursive)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *EmbyVideoTask) FullSync(ctx context.Context, scanPathStr string) error {
	ctx = context.Background()
	boxIpStr, err := t.configRepo.GetConfigBySubKey(ctx, constant.Key_BoxIpMapping, constant.SubKey_EmbyBoxIp)
	if err != nil {
		return err
	}
	boxIps := strings.Split(boxIpStr, ",")
	var scanPaths = make([]string, 0)
	if scanPathStr != "" {
		scanPaths = strings.Split(scanPathStr, "|")
	} else {
		scanPathStr, err = t.configRepo.GetConfigBySubKey(ctx, constant.Key_VideoSyncMapping, constant.SubKey_EmbyVideoSyncCategory)
		if err != nil {
			return err
		}
		scanPaths = strings.Split(scanPathStr, "|")
	}

	for _, scanPath := range scanPaths {
		log.Warnf("开始全量同步emby刮削: pathid=%s", scanPath)
		scanPathAndType := strings.Split(scanPath, ":")
		err = t.deepLoopListEmbyPath(ctx, boxIps[0], scanPathAndType[0], "", scanPathAndType[0], lo.TernaryF(len(scanPathAndType) >= 2, func() string {
			return scanPathAndType[1]
		}, func() string {
			return ""
		}), lo.TernaryF(len(scanPathAndType) >= 3, func() string {
			return scanPathAndType[2]
		}, func() string {
			return "false"
		}))
		if err != nil {
			log.Errorf("全量同步emby刮削失败: %v", err)
			return err
		}
		log.Warnf("结束全量同步emby刮削: pathid=%s", scanPathAndType[0])
	}
	return nil
}

func (t *EmbyVideoTask) deepLoopListEmbyPath(ctx context.Context, domain, parentId string, parentName string, rootPathId string, includeItemType string, recursive string) error {
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
		syncListURL := fmt.Sprintf(domain+constant.EmbyVideoList, embyDefaultUserId, startIndex, parentId, constant.PageSize, recursive)
		if includeItemType != "" {
			syncListURL = syncListURL + "&IncludeItemTypes=" + includeItemType
		}
		videoListResp, err := httpclient_util.DoGet[embydto.VideoListResp](ctx, syncListURL, headerMap)
		if err != nil {
			return fmt.Errorf("emby请求parentId %s 游标 %d 开始的返回结果失败: %v", parentId, startIndex, err)
		}
		if videoListResp == nil {
			return fmt.Errorf("emby parentId %s 游标 %d 开始的返回结果无效", parentId, startIndex)
		}
		for index, content := range videoListResp.Items {
			err = t.deepLoopDetailEmbyPath(ctx, domain, parentId, parentName, len(videoListResp.Items), content.Id, content.Type, nil, nil, embyDefaultUserId, headerMap, index, rootPathId, recursive)
			if err != nil {
				return err
			}
		}
		if int64(startIndex+constant.PageSize) >= videoListResp.TotalRecordCount {
			log.Warnf("==============跳出循环: startIndex=%d, total:%d", startIndex+constant.PageSize, videoListResp.TotalRecordCount)
			break
		}
		startIndex = startIndex + constant.PageSize
		log.Warnf("==============继续执行: startIndex=%d, total:%d", startIndex, videoListResp.TotalRecordCount)
	}
	return nil
}

func (t *EmbyVideoTask) deepLoopDetailEmbyPath(ctx context.Context, domain, parentId string, parentName string, childrenSize int, id string, mediaType string, season *embydto.VideoDetailResp, series *embydto.VideoDetailResp, embyDefaultUserId string, headerMap map[string]string, index int, rootPathId string, recursive string) error {
	//查出详情
	syncDetailURL := fmt.Sprintf(domain+constant.EmbyVideoDetail, embyDefaultUserId, id)
	videoDetailResp, err := httpclient_util.DoGet[embydto.VideoDetailResp](ctx, syncDetailURL, headerMap)
	if err != nil {
		return fmt.Errorf("emby videoDetailResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
	}
	if videoDetailResp == nil {
		return fmt.Errorf("emby videoDetailResp parentId %s file %s 返回结果无效", parentId, id)
	}

	//按类型处理
	if mediaType == constant.JfFolder || mediaType == constant.JFBoxSet {
		return t.deepLoopListEmbyPath(ctx, domain, id, videoDetailResp.Name, rootPathId, "", recursive)
	}
	if mediaType == constant.JfSeries {
		syncSeasonListURL := fmt.Sprintf(domain+constant.EmbySeaonsList, id, 10000, embyDefaultUserId)
		seasonListResp, err := httpclient_util.DoGet[embydto.SeasonListResp](ctx, syncSeasonListURL, headerMap)
		if err != nil {
			return fmt.Errorf("emby seasonListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
		}
		if seasonListResp == nil {
			return fmt.Errorf("emby seasonListResp parentId %s file %s 返回结果无效", parentId, id)
		}
		for newIndex, content := range seasonListResp.Items {
			err := t.deepLoopDetailEmbyPath(ctx, domain, id, videoDetailResp.Name, len(seasonListResp.Items), content.Id, constant.JfSeason, nil, videoDetailResp, embyDefaultUserId, headerMap, newIndex, rootPathId, recursive)
			if err != nil {
				return err
			}
		}
	}
	if mediaType == constant.JfSeason {
		syncEpisodeListURL := fmt.Sprintf(domain+constant.EmbyEpisodesList, parentId, id, 10000, embyDefaultUserId)
		episodeListResp, err := httpclient_util.DoGet[embydto.EpisodeListResp](ctx, syncEpisodeListURL, headerMap)
		if err != nil {
			return fmt.Errorf("emby episodeListResp parentId %s file %s 返回结果失败: %v", parentId, id, err)
		}
		if episodeListResp == nil {
			return fmt.Errorf("emby episodeListResp parentId %s file %s 返回结果无效", parentId, id)
		}
		for newIndex, content := range episodeListResp.Items {
			err := t.deepLoopDetailEmbyPath(ctx, domain, id, videoDetailResp.Name, len(episodeListResp.Items), content.Id, constant.JfEpisode, videoDetailResp, series, embyDefaultUserId, headerMap, newIndex, rootPathId, recursive)
			if err != nil {
				return err
			}
		}
	}

	if mediaType == constant.JfMovie || mediaType == constant.JfEpisode {
		sFilePath, sFileName, err := splitEmbyPath(videoDetailResp.MediaSources[0].Path)
		if err != nil {
			return nil
		}
		filePath, _ := url.PathUnescape(sFilePath)
		fileName, _ := url.PathUnescape(sFileName)
		videoType := t.replaceVideType(ctx, filePath)
		log.Infof("遍历到根节点: xiaoya路径=%s", filePath+fileName)
		_, valid := lo.Find(constant.SupportVideoTypes, func(item string) bool {
			if strings.HasSuffix(fileName, item) {
				return true
			}
			return false
		})
		if !valid {
			log.Warnf("跳过=====文件后缀不处理: xiaoya路径=%s", fileName)
			return nil
		}
		episode, err := t.episodeRepo.QueryByPathAndName(ctx, filePath, fileName)
		if err != nil {
			return err
		}
		if episode != nil {
			log.Warnf("跳过=====先前的任务已经处理过了: xiaoya路径=%s", filePath+fileName)
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
		if episodeSize == 0 {
			log.Warnf("xiaoya不存在该数据: xiaoya路径=%s", filePath+fileName)
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
		var parentVideoDetailResp *embydto.VideoDetailResp
		if videoDetailResp.Type == constant.JfMovie {
			if videoType == "Record" && parentId != "" {
				embyId = fmt.Sprintf("%s|%s", parentId, constant.JfFolder)
			} else if videoType == "Record" && parentId == "" {
				parentVideoDetailResp, err = httpclient_util.DoGet[embydto.VideoDetailResp](ctx, fmt.Sprintf(domain+constant.EmbyVideoDetail, embyDefaultUserId, videoDetailResp.ParentId), headerMap)
				if err != nil {
					return fmt.Errorf("emby parentVideoDetailResp file %s 返回结果失败: %v", videoDetailResp.ParentId, err)
				}
				if parentVideoDetailResp == nil {
					return fmt.Errorf("emby parentVideoDetailResp file %s 返回结果无效", videoDetailResp.ParentId)
				}
				if parentVideoDetailResp.Type == constant.JfFolder {
					embyId = fmt.Sprintf("%s|%s", videoDetailResp.ParentId, constant.JfFolder)
				} else {
					embyId = fmt.Sprintf("%s|%s", videoDetailResp.Id, constant.JfMovie)
				}
			} else {
				embyId = fmt.Sprintf("%s|%s", videoDetailResp.Id, constant.JfMovie)
			}
		}
		if videoDetailResp.Type == constant.JfEpisode {
			if videoDetailResp.SeriesId != "" {
				embyId = fmt.Sprintf("%s|%s", videoDetailResp.SeriesId, constant.JfSeries)
			}
			if videoDetailResp.SeasonId != "" {
				embyId = fmt.Sprintf("%s|%s", videoDetailResp.SeasonId, constant.JfSeason)
			}
		}
		var embyCreateTime time.Time
		if videoDetailResp.Type == constant.JfMovie {
			embyCreateTime = time_util.ParseUtcTime(videoDetailResp.DateCreated)
		}
		if videoDetailResp.Type == constant.JfEpisode {
			embyCreateTime = time_util.ParseUtcTime(videoDetailResp.DateCreated)
			if series != nil {
				embyCreateTime = time_util.ParseUtcTime(series.DateCreated)
			}
			if season != nil {
				embyCreateTime = time_util.ParseUtcTime(season.DateCreated)
			}
		}

		video, err := t.videoRepo.GetByJellyfinId(ctx, embyId)
		if err != nil {
			log.Errorf("查询video失败: %v", err)
			return err
		}
		err = gen.Use(t.videoRepo.Data.DB(ctx)).Transaction(func(tx *gen.Query) error {
			if video == nil {
				title := ""
				if videoDetailResp.Type == constant.JfMovie {
					if videoType == "Record" && parentName != "" {
						title = parentName
					} else if videoType == "Record" && parentName == "" && parentVideoDetailResp != nil && parentVideoDetailResp.Type == constant.JfFolder {
						title = parentVideoDetailResp.Name
					} else {
						title = videoDetailResp.Name
					}
				}
				if videoDetailResp.Type == constant.JfEpisode {
					if videoDetailResp.SeriesId != "" {
						title = videoDetailResp.SeriesName
					}
					if videoDetailResp.SeasonId != "" {
						title = lo.Ternary(series.ChildCount == 1, videoDetailResp.SeriesName, videoDetailResp.SeriesName+"："+videoDetailResp.SeasonName)
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
				overview = strings.TrimSpace(overview)
				runes := []rune(overview)
				overview = lo.TernaryF(len(runes) > 1000, func() string {
					return string(runes[:1000]) + "..."
				}, func() string {
					return overview
				})

				genres := make([]string, 0)
				if videoDetailResp.Type == constant.JfMovie {
					genres = videoDetailResp.Genres
					genres = append(genres, lo.Map(videoDetailResp.GenreItems, func(item *embydto.GenreItem, index int) string {
						return item.Name
					})...)
				}
				if videoDetailResp.Type == constant.JfEpisode {
					genres = videoDetailResp.Genres
					genres = append(genres, lo.Map(videoDetailResp.GenreItems, func(item *embydto.GenreItem, index int) string {
						return item.Name
					})...)
					if series != nil && len(series.Genres) > 0 {
						genres = series.Genres
						genres = append(genres, lo.Map(series.GenreItems, func(item *embydto.GenreItem, index int) string {
							return item.Name
						})...)
					}
					if season != nil && len(season.Genres) != 0 {
						genres = season.Genres
						genres = append(genres, lo.Map(season.GenreItems, func(item *embydto.GenreItem, index int) string {
							return item.Name
						})...)
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
					if videoType == "Record" && childrenSize != 0 {
						totalEpisode = int32(childrenSize)
					} else if videoType == "Record" && parentName == "" && parentVideoDetailResp != nil && parentVideoDetailResp.Type == constant.JfFolder {
						videoListResp, err := httpclient_util.DoGet[embydto.VideoListResp](ctx, fmt.Sprintf(domain+constant.EmbyVideoList, embyDefaultUserId, 0, parentVideoDetailResp.Id, constant.PageSize, false), headerMap)
						if err != nil {
							return fmt.Errorf("emby请求parentId %s 的size开始的返回结果失败: %v", parentId, err)
						}
						if videoListResp == nil {
							return fmt.Errorf("emby parentId %s 的size开始的返回结果无效", parentId)
						}
						totalEpisode = int32(len(videoListResp.Items) + 1)
					} else {
						totalEpisode = 1
					}
				}
				if videoDetailResp.Type == constant.JfEpisode {
					if series != nil {
						totalEpisode = series.ChildCount
					}
					if season != nil && season.ChildCount != 0 {
						totalEpisode = season.ChildCount
					}
					if series != nil && series.Status == "Continuing" {
						totalEpisode = -1
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
				region, externalUrls := t.getRegion(ctx, externalUrls, videoDetailResp.MediaSources[0].Path, videoDetailResp.Regions)
				externalUrlStr := ""
				if len(externalUrls) > 0 {
					externalUrlStr, _ = json_util.MarshalString(externalUrls)
				}

				video = &model.Video{
					Title:              strings.TrimSpace(title),
					VideoType:          videoType,
					VoteRate:           lo.ToPtr(float32(math.Round(goodRatting*10) / 10)),
					Region:             lo.ToPtr(region),
					Description:        lo.ToPtr(overview),
					PublishDay:         lo.ToPtr(lo.Ternary(videoDetailResp.PremiereDate != "", time_util.FormatStrToYYYYMMDD(videoDetailResp.PremiereDate), time_util.FormatStrToYYYYMMDD(videoDetailResp.DateCreated))),
					Thumbnail:          lo.ToPtr(t.getValidThumbnail(ctx, domain, []string{videoDetailResp.SeasonId, videoDetailResp.SeriesId, lo.Ternary(videoDetailResp.Type == constant.JfMovie && videoType == "Record", parentId, videoDetailResp.Id)})),
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
				log.Infof("成功写入【video】: xiaoya路径=%s", filePath+fileName)
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
				log.Infof("成功写入【episode】: xiaoya路径=%s", filePath+fileName)
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
					Character:  lo.ToPtr(lo.Ternary(len(character.Role) > 50, "", character.Role)),
					IsDirector: false,
				})
			}
			err = t.videoActorMappingRepo.BatchCreate(ctx, tx, videoActorMappings)
			if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
				return err
			}
			log.Infof("成功写入【videoActorMappingRepo】: xiaoya路径=%s", filePath+fileName)
			//字幕处理
			result, err := httpclient_util.DoPost[embydto.PlaybackInfoReq, embydto.PlaybackInfo](ctx, domain+fmt.Sprintf(constant.EmbyPlayInfo, videoDetailResp.Id), t.PalyInfoReq, headerMap)
			if err != nil {
				return err
			}
			if result == nil {
				return nil
			}
			err = t.episodeSubtitleMappingRepo.Delete(ctx, nil, episode.ID)
			if err != nil {
				return err
			}
			for _, mediaSource := range result.MediaSources {
				for _, mediaStream := range mediaSource.MediaStreams {
					if mediaStream.Type != "Subtitle" || mediaStream.DeliveryUrl == "" {
						continue
					}
					err = t.episodeSubtitleMappingRepo.Create(ctx, nil, &model.EpisodeSubtitleMapping{
						EpisodeID: episode.ID,
						URL:       mediaStream.DeliveryUrl,
						Title:     mediaStream.DisplayTitle,
						Language:  mediaStream.Language,
						MimeType:  mediaStream.Codec,
					})
					if err != nil {
						return err
					}
					log.Infof("成功写入【subtitle】: xiaoya路径=%s", filePath+fileName)
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
func splitEmbyPath(url string) (string, string, error) {
	var path string
	if strings.HasPrefix(url, "/") {
		pathStart := 0
		pathEnd := strings.LastIndex(url, "/")
		path = url[pathStart:pathEnd]
	} else {
		splits := strings.Split(url, "/d")
		if len(splits) != 2 {
			return "", "", errors.New("path not valid")
		}
		filePath := splits[1] // "/综艺/国内综艺/G%20国家宝藏/S04/xxx.mkv"
		pathEnd := strings.LastIndex(filePath, "/")
		path = filePath[0:pathEnd] // "综艺/国内综艺/G%20国家宝藏/S04"
	}
	filename := filepath.Base(url) // "xxx.mkv"
	if filename == "" {
		return "", "", errors.New("path not valid")
	}
	return path, filename, nil
}

func (t *EmbyVideoTask) replaceRegions(ctx context.Context, regions []string) []string {
	if len(regions) == 0 {
		return nil
	}
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
	genres = lo.UniqBy(genres, func(item string) string {
		return item
	})
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
		url := fmt.Sprintf(constant.EmbyPrimaryThumbnail, id)
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

func (t *EmbyVideoTask) getRegion(ctx context.Context, externalUrls []*embydto.ExternalUrl, xiaoyaPath string, regions []string) (string, []*embydto.ExternalUrl) {
	newRegions := t.replaceRegions(ctx, []string{xiaoyaPath})
	if len(newRegions) != 0 {
		return newRegions[0], externalUrls
	}
	if len(regions) != 0 {
		newRegions = t.replaceRegions(ctx, regions)
		if len(newRegions) != 0 {
			return newRegions[0], externalUrls
		}
		return regions[0], externalUrls
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
