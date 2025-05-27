package video

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/repo/gen"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/spf13/viper"
)

// XiaoyaVideoTask 定义xiaoya视频同步任务
type XiaoyaVideoTask struct {
	episodeRepo gen.IEpisodeDo
	log         *log.Helper
}

// NewXiaoyaVideoTask 初始化任务
func NewXiaoyaVideoTask(episodeRepo gen.IEpisodeDo, logger log.Logger) *XiaoyaVideoTask {
	return &XiaoyaVideoTask{
		episodeRepo: episodeRepo,
		log:         log.NewHelper(logger),
	}
}

// Execute 执行同步任务主逻辑
func (t *XiaoyaVideoTask) Execute(ctx context.Context) error {
	baseURL := viper.GetString("xiaoya_url") + "/api/fs/list"
	initialPath := "/电影"
	password := ""
	perPage := 100

	// 首次请求获取总页数
	total, err := t.fetchTotalPages(ctx, baseURL, initialPath, password, perPage)
	if err != nil {
		t.log.Errorf("获取总页数失败: %v", err)
		return err
	}

	// 分页循环请求
	for page := 1; page <= total; page++ {
		err := t.processPage(ctx, baseURL, initialPath, password, page, perPage)
		if err != nil {
			t.log.Errorf("处理第%d页失败: %v", page, err)
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
			err := t.processPage(ctx, baseURL, newPath, password, 1, perPage)
			if err != nil {
				t.log.Errorf("递归处理目录%s失败: %v", newPath, err)
			}
		} else {
			// 查询是否已存在记录（需补充正确的查询条件）
			existingEpisode, err := t.episodeRepo.Where(
				gen.Episode.XiaoyaPath.Eq(path),
				gen.Episode.EpisodeTitle.Eq(item.Name),
			).First()
			if err != nil {
				t.log.Errorf("查询episode失败: %v", err)
				continue // 跳过本次循环避免重复插入
			}
			if existingEpisode == nil {
				episode := &model.Episode{
					XiaoyaPath:   &path,
					EpisodeTitle: item.Name,
					Size:         item.Size,
					IsValid:      true,
				}
				if err := t.episodeRepo.Create(episode); err != nil {
					t.log.Errorf("写入episode失败: %v", err)
				} else {
					t.log.Infof("成功写入episode: path=%s, title=%s", path, item.Name) // 新增成功日志
				}
			} else {
				t.log.Infof("episode已存在: path=%s, title=%s", path, item.Name) // 新增存在日志
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
