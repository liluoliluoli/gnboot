package sdomain

import (
	appversiondto "github.com/liluoliluoli/gnboot/api/appversion"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type AppVersion struct {
	ID          int64     `json:"id"`
	VersionCode string    `json:"versionCode"`
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

func (d *AppVersion) ConvertToDto() *appversiondto.AppVersion {
	return &appversiondto.AppVersion{
		Id:            int32(d.ID),
		VersionCode:   d.VersionCode,
		VersionName:   d.VersionName,
		ForceUpdate:   d.ForceUpdate,
		ApkUrl:        d.ApkUrl,
		PublishedTime: int32(d.PublishTime.Unix()),
		Remark:        d.Remark,
	}
}
