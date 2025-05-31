package video

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liluoliluoli/gnboot/internal/repo"
	"github.com/liluoliluoli/gnboot/internal/repo/model"

	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
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

// Process 执行同步任务主逻辑（优化配置检查）
func (t *XiaoyaVideoTask) Process(task *sdomain.Task) error {
	ctx := task.Ctx

	// 优化：明确提示配置缺失，可根据需求调整为使用默认值
	xiaoyaURL := viper.GetString("xiaoya_url")
	if xiaoyaURL == "" {
		return fmt.Errorf("未配置 xiaoya_url，请检查配置文件")
	}

	baseURL := xiaoyaURL + "/api/fs/list"
	initialPath := "/电影"
	password := ""
	perPage := 40 // 根据业务需求调整分页大小

	log.Infof("开始同步xiaoya视频: baseURL=%s, initialPath=%s", baseURL, initialPath)
	total, err := t.fetchTotalPages(ctx, baseURL, initialPath, password, perPage)
	if err != nil {
		return fmt.Errorf("获取总页数失败: %w", err)
	}

	// 分页循环（添加分页范围校验）
	if total < 1 {
		log.Info("总页数为0，无需同步")
		return nil
	}
	for page := 1; page <= total; page++ {
		log.Infof("开始处理第%d页", page)
		if err := t.processPage(ctx, baseURL, initialPath, password, page, perPage, 0); err != nil {
			log.Errorf("处理第%d页失败: %v", page, err)
		}
	}
	return nil
}

// processPage 处理单页数据（优化递归逻辑和参数）
func (t *XiaoyaVideoTask) processPage(ctx context.Context, baseURL, currentPath string, password string, page, perPage int, depth int) error { // 关键修改：重命名参数为currentPath
	const maxDepth = 10
	if depth > maxDepth {
		log.Warnf("目录递归超过最大深度%d，终止: %s", maxDepth, currentPath) // 同步修改为currentPath
		return nil
	}

	resp, err := t.requestAPI(ctx, baseURL, currentPath, password, page, perPage) // 传递currentPath给requestAPI
	if err != nil {
		return fmt.Errorf("请求API失败: %w", err)
	}

	for _, item := range resp.Data.Content {
		if item.IsDir {
			newPath := path.Join(currentPath, item.Name) // 使用包名path调用Join，currentPath为当前路径变量
			log.Infof("递归处理目录（深度%d）: %s", depth+1, newPath)
			if err := t.processPage(ctx, baseURL, newPath, password, 1, perPage, depth+1); err != nil { // 传递newPath作为新的currentPath
				log.Errorf("递归处理目录%s失败: %v", newPath, err)
			}
		} else {
			// 查询是否已存在记录（优化判断条件）
			existingEpisode, err := t.episodeRepo.QueryByPathAndName(ctx, currentPath, item.Name) // 查询时使用currentPath
			if err != nil {
				log.Errorf("查询episode失败: %v，跳过当前文件", err)
				continue
			}
			if len(existingEpisode) == 0 {
				// 注意：此处currentPath是函数参数重命名后的变量，与包名无冲突
				episode := &model.Episode{
					XiaoyaPath:   &currentPath, // 指向修改后的变量名
					EpisodeTitle: item.Name,
					Size:         strconv.Itoa(item.Size),
					IsValid:      true,
				}
				if err := t.episodeRepo.Create(ctx, nil, episode); err != nil {
					log.Errorf("写入episode失败（path=%s, title=%s）: %v", currentPath, item.Name, err) // 使用currentPath
				} else {
					log.Infof("成功写入episode: path=%s, title=%s", currentPath, item.Name) // 使用currentPath
				}
			} else {
				log.Infof("episode已存在（跳过）: path=%s, title=%s", currentPath, item.Name) // 使用currentPath
			}
		}
	}
	return nil
}

// requestAPI 调用xiaoya接口（优化HTTP客户端配置）
func (t *XiaoyaVideoTask) requestAPI(ctx context.Context, baseURL, path, password string, page, perPage int) (*APIResponse, error) {
	params := map[string]interface{}{
		"path":     path,
		"password": password,
		"page":     page,
		"per_page": perPage,
		"refresh":  false,
	}

	paramJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("参数序列化失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, strings.NewReader(string(paramJSON)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// 优化：显式设置HTTP客户端超时（防止长时间挂起）
	client := &http.Client{
		Timeout: 10 * time.Second, // 10秒超时
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	log.Infof("API原始响应: %s", string(respBody))

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("反序列化失败（原始响应=%s）: %w", string(respBody), err)
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("接口错误（状态码=%d）: %s（原始响应=%s）", apiResp.Code, apiResp.Msg, string(respBody))
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

// fetchTotalPages 计算总页数（通过第一页数据获取总记录数）
func (t *XiaoyaVideoTask) fetchTotalPages(ctx context.Context, baseURL, path, password string, perPage int) (int, error) {
	// 调用第一页获取总记录数
	resp, err := t.requestAPI(ctx, baseURL, path, password, 1, perPage)
	if err != nil {
		return 0, fmt.Errorf("获取总记录数失败: %w", err)
	}

	// 计算总页数（总记录数 / 每页数量，向上取整）
	if resp.Data.Total <= 0 {
		return 0, nil
	}
	totalPages := (resp.Data.Total + perPage - 1) / perPage // 避免整除时少一页
	return totalPages, nil
}
