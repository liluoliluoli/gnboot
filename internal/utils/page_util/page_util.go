package page_util

import (
	"github.com/go-cinch/common/page"
	"github.com/go-cinch/common/proto/params"
)

func ToDomainPage(p *params.Page) *page.Page {
	if p == nil {
		return &page.Page{
			Num:  1,
			Size: 20,
		}
	}
	return &page.Page{
		Num:  p.Num,
		Size: p.Size,
	}
}

func ToAdaptorPage(p page.Page) *params.Page {
	return &params.Page{
		Num:   p.Num,
		Size:  p.Size,
		Total: p.Total,
	}
}
