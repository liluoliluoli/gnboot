package constant

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

const CTX_UserName = "CTX_UserName"
const CTX_SessionToken = "CTX_SessionToken"

const MaxWatchCountByDay = 200

type PackageType = string

const (
	None  PackageType = "none"
	Trial PackageType = "trial"
	Month PackageType = "month"
	Year  PackageType = "year"
)

var (
	XiaoYaToken string = ""
)

const (
	XiaoYaLoginPath               = "/api/auth/login/hash"
	XiaoYaTransferStorePath       = "/api/fs/get"
	XiaoYaM3u8Path                = "/api/fs/other"
	AliyunM3u8EarlyExpireMinutes  = 2 * 60                                                             //提前失效分钟
	AliyunM3u8ReallyExpireMinutes = 4 * 60                                                             //实际失效分钟
	XiaoYaLoginName               = "admin"                                                            //xiaoya登录账号
	XiaoYaLoginPassword           = "6fcb57cd10b2c11d765dcf16148d99130afd895082af83725ee8bb181b1d2b0f" //xiaoya登录密码
)
