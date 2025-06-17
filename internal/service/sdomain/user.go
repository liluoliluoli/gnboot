package sdomain

import (
	"encoding/json"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"time"
)

type User struct {
	ID                  int64      `json:"id"`
	UserName            string     `json:"userName"`
	Password            string     `json:"password"`
	PackageType         string     `json:"packageType"`
	PackageExpiredTime  *time.Time `json:"packageExpiredTime"`
	SessionToken        *string    `json:"sessionToken"`
	AliToken            *string    `json:"aliToken"`
	AliTokenExpiredTime *time.Time `json:"aliTokenExpiredTime"`
	NoticeTitle         string     `json:"noticeTitle"`
	NoticeContent       string     `json:"noticeContent"`
}

func (d *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *User) ConvertFromRepo(m *model.User) *User {
	return &User{
		ID:                  m.ID,
		UserName:            m.UserName,
		Password:            m.Password,
		PackageType:         m.PackageType,
		PackageExpiredTime:  m.PackageExpiredTime,
		SessionToken:        m.SessionToken,
		AliToken:            m.AliToken,
		AliTokenExpiredTime: m.AliTokenExpiredTime,
	}
}

func (d *User) ConvertToRepo() *model.User {
	return &model.User{
		ID:                  d.ID,
		UserName:            d.UserName,
		Password:            d.Password,
		PackageType:         d.PackageType,
		SessionToken:        d.SessionToken,
		AliToken:            d.AliToken,
		AliTokenExpiredTime: d.AliTokenExpiredTime,
	}
}
