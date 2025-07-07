package httpclient_util

import (
	"context"
	"fmt"
	"regexp"
	"testing"
)

func TestRegular(t *testing.T) {
	ctx := context.Background()
	thumbnail, err := getTmdbInfo(ctx, "https://www.themoviedb.org/tv/38854")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("count:", thumbnail) // 输出: 15
}

func getTmdbInfo(ctx context.Context, url string) (string, error) { //region,genre,thumbnail,ratting,actors
	html, err := DoHtml(ctx, url)
	if err != nil {
		return "", nil
	}
	// 正则提取 window.__DATA__ 的 JSON
	re := regexp.MustCompile(`<img class="poster w-full"[^>]*srcset="[^"]*?\s+1x,\s*([^"\s]+)\s+2x"`)
	match := re.FindStringSubmatch(html)
	if len(match) < 2 {
		return "", nil
	}
	return match[1], nil
}
