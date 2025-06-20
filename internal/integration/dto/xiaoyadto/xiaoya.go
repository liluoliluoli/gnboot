package xiaoyadto

import "time"

type TransferStoreReq struct {
	Path     string `json:"path"`
	Password string `json:"password"`
}

type TransferStoreResp struct {
	Id       string     `json:"id"`
	Path     string     `json:"path"`
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	IsDir    bool       `json:"is_dir"`
	Modified *time.Time `json:"modified"`
	Created  *time.Time `json:"created"`
	Sign     string     `json:"sign"`
	Thumb    string     `json:"thumb"`
	RawUrl   string     `json:"raw_url"`
	Provider string     `json:"provider"`
}

type M3u8Req struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	Method   string `json:"method"`
}

type M3u8Resp struct {
	Category                    string                `json:"category"`
	DriveId                     string                `json:"drive_id"`
	FileId                      string                `json:"file_id"`
	MetaNameInvestigationStatus int32                 `json:"meta_name_investigation_status"`
	MetaNamePunishFlag          int32                 `json:"meta_name_punish_flag"`
	PunishFlag                  int32                 `json:"punish_flag"`
	VideoPreviewPlayInfo        *VideoPreviewPlayInfo `json:"video_preview_play_info"`
}

type VideoPreviewPlayInfo struct {
	Category                string                `json:"category"`
	LiveTranscodingTaskList []LiveTranscodingTask `json:"live_transcoding_task_list"`
	Meta                    *Meta                 `json:"meta"`
}

type LiveTranscodingTask struct {
	Stage          string `json:"stage"`
	Status         string `json:"status"`
	TemplateHeight int64  `json:"template_height"`
	TemplateId     string `json:"template_id"`
	TemplateName   string `json:"template_name"`
	TemplateWidth  int64  `json:"template_width"`
	Url            string `json:"url"`
}

type Meta struct {
	Duration float64 `json:"duration"`
	Height   int64   `json:"height"`
	Width    int64   `json:"width"`
}

type XiaoyaResult[T any] struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OtpCode  string `json:"otp_code"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type VideoListReq struct {
	Path     string `json:"path"`
	Password string `json:"password"`
	Page     int32  `json:"page"`
	PerPage  int32  `json:"per_page"`
	Refresh  bool   `json:"refresh"`
}

type VideoListResp struct {
	Content  []*VideoContent `json:"content"`
	Total    int64           `json:"total"`
	Readme   string          `json:"readme"`
	Write    bool            `json:"write"`
	Provider string          `json:"provider"`
}

type VideoContent struct {
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	IsDir    bool       `json:"is_dir"`
	Modified *time.Time `json:"modified"`
	Sign     string     `json:"sign"`
	Thumb    string     `json:"thumb"`
	Type     int32      `json:"type"`
}
