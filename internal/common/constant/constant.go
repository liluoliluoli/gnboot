package constant

import "regexp"

const (
	GN_OPERATOR_CONTEXT = "gn_operator_context"
	SYS_PWD             = "SDDSIOPOPPP"
)

type Sort = string

const (
	SortByPublish Sort = "publish"
	SortByHot     Sort = "hot"
	SortByRate    Sort = "rate"
)

const RK_UserTokenPrefix = "ut_%s"
const RK_UserWatchCountPrefix = "uwc_%s"
const RK_UserPackagePrefix = "upk_%s"
const RK_Notice = "notice"
const HK_NoticeTitle = "title"
const HK_NoticeContent = "content"
const HK_NoticeDonateImageUrl = "donateImageUrl"
const HK_NoticeDonateDesc = "donateDesc"

const RK_Configs = "configs"

type Key = string

const (
	Key_BoxIpMapping         SubKey = "boxIpMapping"
	Key_GenreMapping         SubKey = "genreMapping"
	Key_RegionMapping        SubKey = "regionMapping"
	Key_VideoSyncMapping     SubKey = "videoSyncMapping"
	Key_PathVideoTypeMapping SubKey = "pathVideoTypeMapping"
)

type SubKey = string

const (
	SubKey_XiaoYaBoxIp               SubKey = "xiaoYaBoxIp"
	SubKey_JellyfinBoxIp             SubKey = "jellyfinBoxIp"
	SubKey_JellyfinVideoSyncCategory SubKey = "jellyfinVideoSyncCategory"
	SubKey_JellyfinDefaultUserId     SubKey = "jellyfinDefaultUserId"
	SubKey_JellyfinDefaultToken      SubKey = "jellyfinDefaultToken"
	SubKey_AliOpenapiDomain          SubKey = "aliOpenapiDomain"
	SubKey_AliOpenapiToken           SubKey = "aliOpenapiToken"
	SubKey_UseAliOpenapi             SubKey = "useAliOpenapi"
	SubKey_EmbyBoxIp                 SubKey = "embyBoxIp"
	SubKey_EmbyVideoSyncCategory     SubKey = "embyVideoSyncCategory"
	SubKey_EmbyDefaultToken          SubKey = "embyDefaultToken"
	SubKey_EmbyDefaultUserId         SubKey = "embyDefaultUserId"
)

const CTX_UserName = "CTX_UserName"
const CTX_SessionToken = "CTX_SessionToken"
const CTX_ClientIp = "CTX_ClientIp"

const MaxWatchCountByDay = 200
const TrailVideoId = 1
const HD_EPISODE_MIN_SIZE = 1024 * 1024 * 1024 * 1.5
const HD_MOVIE_MIN_SIZE = 1024 * 1024 * 1024 * 5
const QHD_MOVIE_MIN_SIZE = 1024 * 1024 * 1024 * 10

type PackageType = string

const (
	None PackageType = "none" //可观看标清视频和测试目录高清视频
	Year PackageType = "year" //可观看高清视频
)

type Radio = string

const (
	LD  Radio = "LD"
	SD  Radio = "SD"
	HD  Radio = "HD"
	QHD Radio = "QHD"
)

type Provider = string

const (
	IMDb       Provider = "IMDb"
	TheMovieDb Provider = "TheMovieDb"
	Trakt      Provider = "Trakt"
	Douban     Provider = "Douban"
)

var (
	XiaoYaToken string                       = ""
	ConfigMap   map[string]map[string]string = make(map[string]map[string]string)
	RegularMap  map[string]*regexp.Regexp    = map[string]*regexp.Regexp{
		IMDb:       regexp.MustCompile(`tt_dt_cnt">(.*?)</a>`),
		TheMovieDb: nil,
		Trakt:      regexp.MustCompile(`<label>Country</label>(.*?)\s*</li>`),
		Douban:     regexp.MustCompile(`<span class="pl">制片国家/地区:</span>\s*(.*?)<br/>`),
	}
	SortMap map[string]int32 = map[string]int32{
		Douban:     1,
		TheMovieDb: 2,
		IMDb:       3,
		Trakt:      4,
	}
	RegularChinese    *regexp.Regexp    = regexp.MustCompile(`[\p{Han}]+`)
	RegionMap         map[string]string = make(map[string]string)
	SupportVideoTypes []string          = []string{".mkv", ".mp4", ".rmvb", ".avi"}
)

const (
	XiaoYaLoginPath              = "/api/auth/login/hash"
	XiaoYaTransferStorePath      = "/api/fs/get"
	XiaoYaM3u8Path               = "/api/fs/other"
	AliyunM3u8Path               = "/adrive/v1.0/openFile/getVideoPreviewPlayInfo"
	XiaoYaVideoList              = "/api/fs/list"
	JellyfinVideoList            = "/Users/%s/Items?SortBy=DateCreated&SortOrder=Descending&Recursive=false&StartIndex=%d&ParentId=%s&Limit=%d"
	JellyfinVideoDetail          = "/Users/%s/Items/%s"
	JellyfinSeaonsList           = "/Shows/%s/Seasons"
	JellyfinEpisodesList         = "/Shows/%s/Episodes?seasonId=%s"
	JellyfinPlayInfo             = "/Items/%s/PlaybackInfo?UserId=%s"
	PrimaryThumbnail             = "/Items/%s/Images/Primary"
	GetCountryDetail             = "https://restcountries.com/v3.1/alpha/%s"
	DefaultThumbnail             = ""
	AliyunM3u8EarlyExpireSeconds = 10 * 60                                                            //提前失效分钟
	XiaoYaLoginName              = "admin"                                                            //xiaoya登录账号
	XiaoYaLoginPassword          = "6fcb57cd10b2c11d765dcf16148d99130afd895082af83725ee8bb181b1d2b0f" //xiaoya登录密码
	Platform                     = "aliyun"
	PageSize                     = 100

	EmbyVideoList        = "/emby/Users/%s/Items?StartIndex=%d&ParentId=%s&Recursive=true&Limit=%d"
	EmbyVideoDetail      = "/emby/Users/%s/Items/%s"
	EmbySeaonsList       = "/emby/Shows/%s/Seasons?Limit=%d"
	EmbyEpisodesList     = "/emby/Shows/%s/Episodes?SeasonId=%s&Limit=%d"
	EmbyPrimaryThumbnail = "/emby/Items/%s/Images/Primary"
	EmbyPlayInfo         = "/emby/Items/%s/PlaybackInfo"
)

type JfVideoType = string

const (
	JfMovie   JfVideoType = "Movie"
	JfEpisode JfVideoType = "Episode"
	JfSeason  JfVideoType = "Season"
	JfSeries  JfVideoType = "Series"
	JfFolder  JfVideoType = "Folder"
	JFBoxSet  JfVideoType = "BoxSet"
)
