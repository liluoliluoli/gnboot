package sdomain

import (
	"encoding/json"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
	"github.com/samber/lo"
	"time"
)

type User struct {
	ID                  int64     `json:"id"`
	UserName            string    `json:"userName"`
	Password            string    `json:"password"`
	AliToken            string    `json:"aliToken"`
	AliTokenExpiredTime time.Time `json:"aliTokenExpiredTime"`
	WatchCount          int32     `json:"watchCount"`
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
		AliToken:            lo.FromPtr(m.AliToken),
		AliTokenExpiredTime: lo.FromPtr(m.AliTokenExpiredTime),
		WatchCount:          lo.FromPtr(m.WatchCount),
	}
}
