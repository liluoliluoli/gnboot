package alidto

type GetVideoPreviewPlayInfoReq struct {
	DriveId         string `json:"drive_id"`
	FileId          string `json:"file_id"`
	Category        string `json:"category"`
	UrlExpireSec    int64  `json:"url_expire_sec"`
	GetSubtitleInfo bool   `json:"get_subtitle_info"`
}

type VideoPreviewPlayInfoResp struct {
	DomainId             string                `json:"domain_id"`
	DriveId              string                `json:"drive_id"`
	FileId               string                `json:"file_id"`
	VideoPreviewPlayInfo *VideoPreviewPlayInfo `json:"video_preview_play_info"`
}

type VideoPreviewPlayInfo struct {
	Category                        string                         `json:"category"`
	LiveTranscodingTaskList         []*LiveTranscodingTask         `json:"live_transcoding_task_list"`
	LiveTranscodingSubtitleTaskList []*LiveTranscodingSubtitleTask `json:"live_transcoding_subtitle_task_list"`
}

type LiveTranscodingTask struct {
	TemplateId  string `json:"template_id"`
	Status      string `json:"status"` //finished running failed
	Url         string `json:"url"`
	Description string `json:"description"`
}

type LiveTranscodingSubtitleTask struct {
	Language string `json:"language"`
	Status   string `json:"status"`
	Url      string `json:"url"`
}
