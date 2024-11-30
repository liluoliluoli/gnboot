package sdomain

import (
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type VideoKeywordMapping struct {
	ID        int64  `json:"id"`
	VideoId   int64  `json:"videoId"`
	VideoType string `json:"type"`
	KeywordId int64  `json:"keywordId"`
}

func (d *VideoKeywordMapping) ConvertFromRepo(m *model.VideoKeywordMapping) *VideoKeywordMapping {
	return &VideoKeywordMapping{
		ID:        m.ID,
		VideoId:   m.VideoID,
		VideoType: m.VideoType,
		KeywordId: m.KeywordID,
	}
}
