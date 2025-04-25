package sdomain

import (
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type AppVersion struct {
	ID          int64     `json:"id"`
	VersionCode int32     `json:"versionCode"`
	VersionName string    `json:"versionName"`
	ForceUpdate bool      `json:"forceUpdate"`
	ApkUrl      string    `json:"apkUrl"`
	PublishTime time.Time `json:"publishTime"`
	Remark      string    `json:"remark"`
}

func (d *AppVersion) ConvertFromRepo(m *model.AppVersion) *AppVersion {
	return &AppVersion{
		ID:          m.ID,
		VersionCode: m.VersionCode,
		VersionName: m.VersionName,
		ForceUpdate: m.ForceUpdate,
		ApkUrl:      m.ApkURL,
		PublishTime: m.PublishTime,
		Remark:      lo.FromPtr(m.Remark),
	}
}
