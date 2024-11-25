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
	MovieFilterType_genre   MovieFilterType = "genre"
	MovieFilterType_studio  MovieFilterType = "studio"
	MovieFilterType_keyword MovieFilterType = "keyword"
	MovieFilterType_actor   MovieFilterType = "actor"
)
