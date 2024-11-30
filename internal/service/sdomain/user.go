package sdomain

import (
	"encoding/json"
	"github.com/liluoliluoli/gnboot/internal/repo/model"
)

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName"`
	Nick     string `json:"nick"`
}

func (d *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *User) ConvertFromRepo(m *model.User) *User {
	return &User{
		ID:       m.ID,
		UserName: m.UserName,
		Nick:     m.Nick,
	}
}
