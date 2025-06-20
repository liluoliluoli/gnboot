package video

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/json_util"
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/integration/dto"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"path"
	"strings"

	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type XiaoyaVideoTask struct {
	episodeRepo *repo.EpisodeRepo
	client      redis.UniversalClient
	c           *conf.Bootstrap
}

func NewXiaoyaVideoTask(episodeRepo *repo.EpisodeRepo, client redis.UniversalClient, c *conf.Bootstrap) *XiaoyaVideoTask {
	return &XiaoyaVideoTask{
		episodeRepo: episodeRepo,
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
	syncURL := boxIps[0][constant.Key_XiaoYaBoxIp] + constant.XiaoYaVideoList
	scanPaths := strings.Split(t.c.XiaoyaVideoSyncCategory, ",")
	pageSize := int32(100)
	for _, scanPath := range scanPaths {
		log.Infof("开始同步xiaoya视频: syncURL=%s, initialPath=%s", syncURL, scanPath)
		err := t.deepLoopPath(ctx, syncURL, scanPath, "", pageSize)
		if err != nil {
			return err
		}
		log.Infof("结束同步xiaoya视频: syncURL=%s, initialPath=%s", syncURL, scanPath)
	}
	return nil
}

func (t *XiaoyaVideoTask) deepLoopPath(ctx context.Context, syncURL, currentPath, parentPath string, pageSize int32) error {
	page := int32(1)
	for {
		result, err := httpclient_util.DoPost[dto.VideoListReq, dto.XiaoyaResult[dto.VideoListResp]](ctx, syncURL, "", &dto.VideoListReq{
			Path:     currentPath,
			Password: "",
			Page:     page,
			PerPage:  pageSize,
			Refresh:  false,
		})

		if err != nil {
			return fmt.Errorf("请求路径 %s 第 %d 页失败: %v", currentPath, page, err)
		}

		if result == nil || result.Code != 200 || result.Data == nil {
			return fmt.Errorf("路径 %s 第 %d 页返回结果无效", currentPath, page)
		}

		for _, content := range result.Data.Content {
			// 记录当前处理的内容信息（包含父级路径）
			log.Infof("处理内容: 路径=%s, 父级=%s, 名称=%s, 是目录=%v", currentPath, parentPath, content.Name, content.IsDir)
			if content.IsDir {
				newPath := path.Join(currentPath, content.Name)
				err := t.deepLoopPath(ctx, syncURL, newPath, content.Name, pageSize)
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

		// 判断是否还有下一页
		if int64(page)*int64(pageSize) >= result.Data.Total {
			break
		}
		page++
	}

	return nil
}
