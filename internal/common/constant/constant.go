package constant

const (
	GN_OPERATOR_CONTEXT = "gn_operator_context"
)

type VideoType = string

const (
	VideoType_movie   VideoType = "movie"
	VideoType_series  VideoType = "series"
	VideoType_season  VideoType = "season"
	VideoType_episode VideoType = "episode"
)

type MovieFilterType = string

const (
	FilterType_genre   MovieFilterType = "genre"
	FilterType_studio  MovieFilterType = "studio"
	FilterType_keyword MovieFilterType = "keyword"
	FilterType_actor   MovieFilterType = "actor"
)
