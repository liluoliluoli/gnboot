package video

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/liluoliluoli/gnboot/internal/service/sdomain" // 新增sdomain导入
	"github.com/spf13/viper"
)

// XiaoyaVideoTask 定义xiaoya视频同步任务
type XiaoyaVideoTask struct {
	episodeRepo *repo.EpisodeRepo
}

// NewXiaoyaVideoTask 初始化任务（保持原构造函数不变）
func NewXiaoyaVideoTask(episodeRepo *repo.EpisodeRepo) *XiaoyaVideoTask {
	return &XiaoyaVideoTask{
		episodeRepo: episodeRepo,
	}
}

// Process 执行同步任务主逻辑（原Execute方法重命名，并调整参数为*sdomain.Task）
func (t *XiaoyaVideoTask) Process(task *sdomain.Task) error {
	// 从任务上下文中获取原始ctx（根据项目规范可能需要调整）
	ctx := task.Ctx

	baseURL := viper.GetString("xiaoya_url") + "/api/fs/list"
	initialPath := "/电影"
	password := ""
	perPage := 100

	// 首次请求获取总页数
	log.Infof("开始同步xiaoya视频: baseURL=%s, initialPath=%s", baseURL, initialPath)
	total, err := t.fetchTotalPages(ctx, baseURL, initialPath, password, perPage)
	if err != nil {
		log.Errorf("获取总页数失败: %v", err)
		return err
	}

	// 分页循环请求
	for page := 1; page <= total; page++ {
		log.Infof("开始处理第%d页", page)
		err := t.processPage(ctx, baseURL, initialPath, password, page, perPage)
		if err != nil {
			log.Errorf("处理第%d页失败: %v", page, err)
			continue
		}
	}
	return nil
}

// fetchTotalPages 获取总页数
func (t *XiaoyaVideoTask) fetchTotalPages(ctx context.Context, baseURL, path, password string, perPage int) (int, error) {
	resp, err := t.requestAPI(ctx, baseURL, path, password, 1, perPage)
	if err != nil {
		return 0, err
	}
	log.Infof("获取总页数: total=%d, perPage=%d", resp.Data.Total, perPage)
	totalPages := (resp.Data.Total + perPage - 1) / perPage // 向上取整计算总页数
	return totalPages, nil
}

// processPage 处理单页数据
func (t *XiaoyaVideoTask) processPage(ctx context.Context, baseURL, path, password string, page, perPage int) error {
	resp, err := t.requestAPI(ctx, baseURL, path, password, page, perPage)
	if err != nil {
		return err
	}

	for _, item := range resp.Data.Content {
		if item.IsDir {
			// 目录：递归请求
			newPath := fmt.Sprintf("%s/%s", path, item.Name)
			log.Infof("递归处理目录: %s", newPath)
			err := t.processPage(ctx, baseURL, newPath, password, 1, perPage)
			if err != nil {
				log.Errorf("递归处理目录%s失败: %v", newPath, err)
			}
		} else {
			// 查询是否已存在记录（需补充正确的查询条件）
			existingEpisode, err := t.episodeRepo.Get(ctx, 1) // todo：这里在episodeRepo实现你的查询方法
			if err != nil {
				log.Errorf("查询episode失败: %v", err)
				continue // 跳过本次循环避免重复插入
			}
			if existingEpisode == nil {
				log.Infof("执行写入episode: path=%s, title=%s", path, item.Name)
				episode := &model.Episode{
					XiaoyaPath:   &path,
					EpisodeTitle: item.Name,
					Size:         item.Size,
					IsValid:      true,
				}
				if err := t.episodeRepo.Create(ctx, nil, episode); err != nil {
					log.Errorf("写入episode失败: %v", err)
				} else {
					log.Infof("成功写入episode: path=%s, title=%s", path, item.Name) // 新增成功日志
				}
			} else {
				log.Infof("episode已存在: path=%s, title=%s", path, item.Name) // 新增存在日志
			}
		}
	}
	return nil
}

// requestAPI 调用xiaoya接口
func (t *XiaoyaVideoTask) requestAPI(ctx context.Context, baseURL, path, password string, page, perPage int) (*APIResponse, error) {
	params := map[string]string{
		"path":     path,
		"password": password,
		"page":     strconv.Itoa(page),
		"per_page": strconv.Itoa(perPage),
		"refresh":  "false",
	}

	paramJSON, _ := json.Marshal(params)
	resp, err := http.Post(baseURL, "application/json", strings.NewReader(string(paramJSON)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if apiResp.Code != 200 {
		return nil, fmt.Errorf("接口返回错误: %s", apiResp.Msg)
	}
	return &apiResp, nil
}

// APIResponse 定义接口响应结构
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total   int `json:"total"`
		Content []struct {
			Name  string `json:"name"`
			IsDir bool   `json:"is_dir"`
			Size  string `json:"size"`
		} `json:"content"`
	} `json:"data"`
}
