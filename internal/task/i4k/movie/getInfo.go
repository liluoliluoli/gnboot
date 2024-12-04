package movie

import (
	"encoding/json"
	"github.com/go-cinch/common/log"
	"github.com/liluoliluoli/gnboot/internal/task/utils"
	"strconv"
)

// i4k host
var host = "http://43.143.112.172:8168"

// 视频列表
var movieListURL = host + "/4k/getlist.php?calss=%E7%94%B5%E5%BD%B1&area=&year=&type=&pg="

// 视频详情页
var movieDetailURL = host + "/aliyun/getlist.php?fileid=null&marker=null&token=null&from=&id="

// 视频播放
//var playURL = host + "/aliyun/getlist.php?uid=zwbzwbxz&ukey=1234512345&fileid=null&marker=null&token=null&from=&id="

func main() {
	sum := 0
	for sum <= 10 {
		TaskList(strconv.Itoa(sum))
		sum += sum
	}
}

func TaskList(page string) {

	var movieList MovieList
	var listRe = "{\"list\": " + utils.FastGetWithDo(movieListURL+page) + "}"
	log.Info("movie list is: %s", listRe)

	// 解组 JSON 数据
	err := json.Unmarshal([]byte(listRe), &movieList)
	if err != nil {
		log.Error("Error unmarshalling JSON: %v", err)
	}

	// 打印解组后的数据
	for _, movie := range movieList.List {
		log.Info("ID: %s, Title: %s, Year: %s, Remarks: %s, Pic: %s\n",
			movie.ID, movie.Title, movie.Year, movie.Remarks, movie.Pic)
	}

}

// 列表结构体
type ReMovie struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Pic     string `json:"pic"`
	Year    string `json:"year"`
	Remarks string `json:"remarks"`
}

// MovieList 是包含电影列表的容器
type MovieList struct {
	List []ReMovie `json:"list"`
}

// 详情页结构体
type File struct {
	ShareID string `json:"share_id"` // 分享ID
	DriveID string `json:"drive_id"` // 云盘ID
	FileID  string `json:"file_id"`  // 文件ID
	Name    string `json:"name"`     // 文件名
	Type    string `json:"model"`    // 文件类型
}

type Detail struct {
	From   string `json:"from"`   // 来源
	Total  int    `json:"total"`  // 总数
	FileID string `json:"fileid"` // 文件ID
	Marker string `json:"marker"` // 标记
	Token  string `json:"token"`  // 令牌
	List   []File `json:"list"`   // 文件列表
}

// 原画播放结构体&转码结构体
type DownloadInfo struct {
	Type     string `json:"model"`    // 下载类型
	URL      string `json:"url"`      // 下载URL
	Referer  string `json:"referer"`  // 引用页URL
	Language string `json:"language"` // 语言信息
	Subtitle string `json:"subtitle"` // 字幕信息
	Delfile  int    `json:"delfile"`  // 删除文件标识
	DeleCopy bool   `json:"delecopy"` // 是否删除副本
}
