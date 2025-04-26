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
