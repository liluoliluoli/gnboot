package embydto

type VideoItem struct {
	Name           string `json:"Name"`
	Id             string `json:"Id"`
	PremiereDate   string `json:"PremiereDate"`
	ProductionYear int64  `json:"ProductionYear"`
	IsFolder       bool   `json:"IsFolder"`
	Type           string `json:"Type"`
}

type VideoListResp struct {
	TotalRecordCount int64        `json:"TotalRecordCount"`
	Items            []*VideoItem `json:"Items"`
}

type VideoDetailResp struct {
	Name           string         `json:"Name"`
	OriginalTitle  string         `json:"OriginalTitle"`
	Id             string         `json:"Id"`
	PremiereDate   string         `json:"PremiereDate"`
	MediaSources   []*MediaSource `json:"MediaSources"`
	BadRating      int64          `json:"CriticRating"`
	Regions        []string       `json:"ProductionLocations"`
	Overview       string         `json:"Overview"`
	Genres         []string       `json:"Genres"`
	GoodRating     float64        `json:"CommunityRating"`
	ProductionYear int64          `json:"ProductionYear"`
	ParentId       string         `json:"ParentId"`
	Type           string         `json:"Type"` //Movie
	Characters     []*People      `json:"People"`
	MediaType      string         `json:"MediaType"` //Video
	DateCreated    string         `json:"DateCreated"`
	IsFolder       bool           `json:"IsFolder"`
	SeriesId       string         `json:"SeriesId"`
	SeriesName     string         `json:"SeriesName"`
	SeasonId       string         `json:"SeasonId"`
	SeasonName     string         `json:"SeasonName"`
	ChildCount     int32          `json:"ChildCount"`
	ExternalUrls   []*ExternalUrl `json:"ExternalUrls"`
}

type MediaSource struct {
	Path         string         `json:"Path"`
	Name         string         `json:"Name"`
	Protocol     string         `json:"Protocol"`
	MediaStreams []*MediaStream `json:"MediaStreams"`
	Container    string         `json:"Container"`
	Size         int64          `json:"Size"`
	Bitrate      int64          `json:"Bitrate"`
	Duration     int64          `json:"RunTimeTicks"`
}

type MediaStream struct {
	Codec                string `json:"Codec"`
	Language             string `json:"Language"`
	TimeBase             string `json:"TimeBase"`
	DisplayTitle         string `json:"DisplayTitle"`
	IsForced             bool   `json:"IsForced"`
	Type                 string `json:"Type"`
	Index                int64  `json:"Index"`
	Score                int64  `json:"Score"`
	IsExternal           bool   `json:"IsExternal"`
	IsTextSubtitleStream bool   `json:"IsTextSubtitleStream"`
	Path                 string `json:"Path"`
	DeliveryUrl          string `json:"DeliveryUrl"`
}

type People struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
	Role string `json:"Role"`
	Type string `json:"Type"` //Actor,Director
}

type ExternalUrl struct {
	Name        string `json:"Name"`
	Url         string `json:"Url"`
	UsedCountry bool   `json:"UsedCountry"`
}

type SeasonListResp struct {
	TotalRecordCount int64         `json:"TotalRecordCount"`
	Items            []*SeasonItem `json:"Items"`
}

type SeasonItem struct {
	Name       string `json:"Name"`
	Id         string `json:"Id"`
	IsFolder   bool   `json:"IsFolder"`
	Type       string `json:"Type"`
	SeriesName string `json:"SeriesName"`
	SeriesId   string `json:"SeriesId"`
}

type EpisodeListResp struct {
	TotalRecordCount int64          `json:"TotalRecordCount"`
	Items            []*EpisodeItem `json:"Items"`
}

type EpisodeItem struct {
	Name       string `json:"Name"`
	Id         string `json:"Id"`
	IsFolder   bool   `json:"IsFolder"`
	Type       string `json:"Type"`
	SeriesName string `json:"SeriesName"`
	SeriesId   string `json:"SeriesId"`
	SeasonId   string `json:"SeasonId"`
	SeasonName string `json:"SeasonName"`
}

type PlaybackInfo struct {
	MediaSources []*MediaSource `json:"MediaSources"`
}
