package httpclient_util

import (
	"fmt"
	"testing"
)

func TestRegular(t *testing.T) {
	//ctx := context.Background()
	html, err := CheckImageUrl("https://2a43-45-196-216-126.ngrok-free.app/Items/8d52885cdf2977c47dad7b060e5dd5e9/Images/Primary")
	if err != nil {
		fmt.Println("find:" + err.Error())
	}
	fmt.Println("find:" + fmt.Sprintf("%v", html))
}
