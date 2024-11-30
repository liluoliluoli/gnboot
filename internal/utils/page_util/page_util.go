package page_util

import (
	"github.com/liluoliluoli/gnboot/api"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
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
		TotalPage:   p.TotalPage,
		Count:       p.Count,
	}
}
