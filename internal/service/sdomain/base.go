package sdomain

import "github.com/go-cinch/common/page"

type PageResult[T any] struct {
	Page *page.Page `json:"page"`
	List []T        `json:"list"`
}
