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
	Status         string         `json:"Status"`
	GenreItems     []*GenreItem   `json:"GenreItems"`
}

type GenreItem struct {
	Id   int64  `json:"Id"`
	Name string `json:"Name"`
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

type PlaybackInfoReq struct {
	DeviceProfile *DeviceProfile `json:"DeviceProfile"`
}

type DeviceProfile struct {
	MaxStaticBitrate                 int64                 `json:"MaxStaticBitrate"`
	MaxStreamingBitrate              int64                 `json:"MaxStreamingBitrate"`
	MusicStreamingTranscodingBitrate int64                 `json:"MusicStreamingTranscodingBitrate"`
	DirectPlayProfiles               []*DirectPlayProfile  `json:"DirectPlayProfiles"`
	TranscodingProfiles              []*TranscodingProfile `json:"TranscodingProfiles"`
	CodecProfiles                    []*CodecProfile       `json:"CodecProfiles"`
	SubtitleProfiles                 []*SubtitleProfile    `json:"SubtitleProfiles"`
	ResponseProfiles                 []*ResponseProfile    `json:"ResponseProfiles"`
}

type DirectPlayProfile struct {
	Container  *string `json:"Container"`
	Type       *string `json:"Type"`
	VideoCodec *string `json:"VideoCodec"`
	AudioCodec *string `json:"AudioCodec"`
}

type TranscodingProfile struct {
	Container           *string `json:"Container"`
	Type                *string `json:"Type"`
	AudioCodec          *string `json:"AudioCodec"`
	Context             *string `json:"Context"`
	Protocol            *string `json:"Protocol"`
	MaxAudioChannels    *string `json:"MaxAudioChannels"`
	MinSegments         *string `json:"MinSegments"`
	BreakOnNonKeyFrames *bool   `json:"BreakOnNonKeyFrames"`
}

type CodecProfile struct {
	Type       *string      `json:"Type"`
	Codec      *string      `json:"Codec"`
	Conditions []*Condition `json:"Conditions"`
}

type Condition struct {
	Condition  *string `json:"Condition"`
	Property   *string `json:"Property"`
	Value      *string `json:"Value"`
	IsRequired *string `json:"IsRequired"`
}

type SubtitleProfile struct {
	Format *string `json:"Format"`
	Method *string `json:"Method"`
}

type ResponseProfile struct {
	Type      *string `json:"Type"`
	Container *string `json:"Container"`
	MimeType  *string `json:"MimeType"`
}
