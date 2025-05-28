package video

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/model"

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

	// 尝试获取 xiaoya_url 配置项，如果获取失败则使用默认值
	defaultXiaoyaURL := "http://127.0.0.1:5678"
	xiaoyaURL := viper.GetString("xiaoya_url")
	if xiaoyaURL == "" {
		xiaoyaURL = defaultXiaoyaURL
	}
	baseURL := xiaoyaURL + "/api/fs/list"
	initialPath := "/电影"
	password := ""
	perPage := 40

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
func (t *XiaoyaVideoTask) fetchTotalPages(ctx context.Context, baseURL, path string, password string, perPage int) (int, error) {
	resp, err := t.requestAPI(ctx, baseURL, path, password, 1, perPage)
	if err != nil {
		return 0, err
	}
	log.Infof("获取总页数: total=%d, perPage=%d", resp.Data.Total, perPage)
	totalPages := (resp.Data.Total + perPage - 1) / perPage // 向上取整计算总页数
	return totalPages, nil
}

// processPage 处理单页数据
func (t *XiaoyaVideoTask) processPage(ctx context.Context, baseURL, path string, password string, page, perPage int) error {
	resp, err := t.requestAPI(ctx, baseURL, path, password, page, perPage)
	if err != nil {
		return err
	}

	for _, item := range resp.Data.Content {
		//如果是目录就递归请求，否则就写入到数据库中
		if !item.IsDir {
			newPath := fmt.Sprintf("%s/%s", path, item.Name)
			log.Infof("递归处理目录: %s", newPath)
			err := t.processPage(ctx, baseURL, newPath, password, 1, perPage)
			if err != nil {
				log.Errorf("递归处理目录%s失败: %v", newPath, err)
			}
		} else {
			// 查询是否已存在记录（需补充正确的查询条件）
			existingEpisode, err := t.episodeRepo.QueryByPathAndName(ctx, path, item.Name)
			if err != nil {
				log.Errorf("查询episode失败: %v", err)
				continue // 跳过本次循环避免重复插入
			}
			if existingEpisode == nil {
				log.Infof("执行写入episode: path=%s, title=%s", path, item.Name)
				episode := &model.Episode{
					XiaoyaPath:   &path,
					EpisodeTitle: item.Name,
					Size:         strconv.Itoa(item.Size),
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
	// 修改：使用 map[string]interface{} 允许混合类型
	params := map[string]interface{}{
		"path":     path,
		"password": password,
		"page":     page,    // 直接传递 int 类型，无需转换为字符串
		"per_page": perPage, // 直接传递 int 类型
		"refresh":  false,
	}

	log.Infof("请求参数: baseURL=%s, path=%s, password=%s, page=%d, perPage=%d", baseURL, path, password, page, perPage)

	paramJSON, err := json.Marshal(params)
	if err != nil {
		log.Errorf("参数JSON序列化失败: %v", err)
		return nil, err
	}

	// 新增：设置HTTP请求超时（使用ctx的超时控制）
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, strings.NewReader(string(paramJSON)))
	if err != nil {
		log.Errorf("创建请求失败: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8") // 明确设置Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("请求API失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取完整响应体用于调试
	respBody, err := io.ReadAll(resp.Body) // 新增：读取完整响应内容
	if err != nil {
		log.Errorf("读取响应体失败: %v", err)
		return nil, err
	}
	log.Infof("API原始响应: %s", string(respBody)) // 记录完整响应

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		log.Errorf("响应反序列化失败: 原始响应=%s, 错误=%v", string(respBody), err) // 记录原始响应+错误
		return nil, err
	}

	if apiResp.Code != 200 {
		// 记录完整错误信息（包括响应体）
		log.Errorf("接口返回错误: 状态码=%d, 消息=%s, 原始响应=%s", apiResp.Code, apiResp.Msg, string(respBody))
		return nil, fmt.Errorf("接口返回错误: %s（状态码=%d）", apiResp.Msg, apiResp.Code)
	}

	log.Infof("API成功响应: total=%d, 内容数量=%d", apiResp.Data.Total, len(apiResp.Data.Content))
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
			Size  int    `json:"size"`
		} `json:"content"`
	} `json:"data"`
}
