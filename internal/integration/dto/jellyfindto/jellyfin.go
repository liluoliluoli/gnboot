package jellyfindto

type VideoItem struct {
	Name            string  `json:"Name"`
	Id              string  `json:"Id"`
	PremiereDate    string  `json:"PremiereDate"`
	CriticRating    int64   `json:"CriticRating"`
	OfficialRating  string  `json:"OfficialRating"`
	CommunityRating float64 `json:"CommunityRating"`
	ProductionYear  int64   `json:"ProductionYear"`
	IsFolder        bool    `json:"IsFolder"`
	Type            string  `json:"Type"`
	MediaType       string  `json:"MediaType"`
}

type VideoListResp struct {
	TotalRecordCount int64        `json:"TotalRecordCount"`
	StartIndex       int64        `json:"StartIndex"`
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
}

type MediaSource struct {
	Path     string `json:"Path"`
	Name     string `json:"Name"`
	Protocol string `json:"Protocol"`
}

type People struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
	Role string `json:"Role"`
	Type string `json:"Type"` //Actor,Director
}
