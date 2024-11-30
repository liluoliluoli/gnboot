package page_util

import (
	"github.com/go-cinch/common/proto/params"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

func ToDomainPage(p *params.Page) *sdomain.Page {
	if p == nil {
		return &sdomain.Page{
			CurrentPage: 1,
			PageSize:    20,
		}
	}
	return &sdomain.Page{
		CurrentPage: p.Num,
		PageSize:    p.Size,
	}
}

func ToAdaptorPage(p *sdomain.Page) *params.Page {
	return &params.Page{
		Num:   p.CurrentPage,
		Size:  p.PageSize,
		Total: p.TotalPage,
	}
}
