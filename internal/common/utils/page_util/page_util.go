package page_util

import (
	"github.com/liluoliluoli/gnboot/api"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

func ToDomainPage(p *api.Page) *sdomain.Page {
	if p == nil {
		return &sdomain.Page{
			CurrentPage: 1,
			PageSize:    20,
		}
	}
	return &sdomain.Page{
		CurrentPage: p.CurrentPage,
		PageSize:    p.PageSize,
	}
}

func ToAdaptorPage(p *sdomain.Page) *api.Page {
	return &api.Page{
		CurrentPage: p.CurrentPage,
		PageSize:    p.PageSize,
		TotalPage:   lo.Ternary(int32(p.Count)%p.PageSize == 0, int32(p.Count)/p.PageSize, int32(p.Count)/p.PageSize+1),
		Count:       int32(p.Count),
	}
}
