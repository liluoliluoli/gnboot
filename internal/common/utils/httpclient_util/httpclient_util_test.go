package httpclient_util

import (
	"context"
	"fmt"
	"github.com/biter777/countries"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"testing"
)

func TestRegular(t *testing.T) {
	ctx := context.Background()
	html, _ := DoHtml(ctx, "https://www.imdb.com/title/tt18302114")
	regular := constant.RegularMap[constant.IMDb]
	matches := regular.FindStringSubmatch(html)
	if len(matches) > 1 {
		country := countries.ByName(matches[1])
		fmt.Println("find:" + country.Domain().String())
	}
	//rs := GetLocalNameByName(ctx, "United States")
	//fmt.Println("find:" + rs)
}
