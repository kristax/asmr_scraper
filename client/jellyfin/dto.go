package jellyfin

import (
	"encoding/json"
	"time"
)

type ViewsResponse struct {
	Items            []*ViewItem `json:"Items"`
	TotalRecordCount int         `json:"TotalRecordCount"`
	StartIndex       int         `json:"StartIndex"`
}

type ViewItem struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
}

type ItemsResponse struct {
	Items            []*Items `json:"Items"`
	TotalRecordCount int      `json:"TotalRecordCount"`
	StartIndex       int      `json:"StartIndex"`
}

type Items struct {
	Name              string        `json:"Name"`
	ServerId          string        `json:"ServerId"`
	Id                string        `json:"Id"`
	SortName          string        `json:"SortName"`
	ChannelId         interface{}   `json:"ChannelId"`
	RunTimeTicks      int64         `json:"RunTimeTicks"`
	IsFolder          bool          `json:"IsFolder"`
	Type              string        `json:"Type"`
	UserData          *UserData     `json:"UserData"`
	Artists           []string      `json:"Artists"`
	ArtistItems       []*Subject    `json:"ArtistItems"`
	AlbumArtists      []*Subject    `json:"AlbumArtists"`
	ImageTags         struct{}      `json:"ImageTags"`
	BackdropImageTags []interface{} `json:"BackdropImageTags"`
	ImageBlurHashes   struct{}      `json:"ImageBlurHashes"`
	LocationType      string        `json:"LocationType"`
	PremiereDate      time.Time     `json:"PremiereDate,omitempty"`
	ProductionYear    int           `json:"ProductionYear,omitempty"`
	AlbumArtist       string        `json:"AlbumArtist,omitempty"`
}

type UserData struct {
	PlaybackPositionTicks int    `json:"PlaybackPositionTicks"`
	PlayCount             int    `json:"PlayCount"`
	IsFavorite            bool   `json:"IsFavorite"`
	Played                bool   `json:"Played"`
	Key                   string `json:"Key"`
}

type Subject struct {
	Name string `json:"Name,omitempty"`
	Id   string `json:"Id,omitempty"`
}

type ItemInfoResponse struct {
	Name                         string          `json:"Name"`
	OriginalTitle                string          `json:"OriginalTitle"`
	ServerId                     string          `json:"ServerId"`
	Id                           string          `json:"Id"`
	Etag                         string          `json:"Etag"`
	DateCreated                  time.Time       `json:"DateCreated"`
	CanDelete                    bool            `json:"CanDelete"`
	CanDownload                  bool            `json:"CanDownload"`
	PreferredMetadataLanguage    string          `json:"PreferredMetadataLanguage"`
	PreferredMetadataCountryCode string          `json:"PreferredMetadataCountryCode"`
	SortName                     string          `json:"SortName"`
	ForcedSortName               string          `json:"ForcedSortName"`
	PremiereDate                 time.Time       `json:"PremiereDate"`
	ExternalUrls                 []interface{}   `json:"ExternalUrls"`
	Path                         string          `json:"Path"`
	EnableMediaSourceDisplay     bool            `json:"EnableMediaSourceDisplay"`
	CustomRating                 string          `json:"CustomRating"`
	ChannelId                    interface{}     `json:"ChannelId"`
	Overview                     string          `json:"Overview"`
	Taglines                     []interface{}   `json:"Taglines"`
	Genres                       []string        `json:"Genres"`
	CommunityRating              float64         `json:"CommunityRating"`
	CumulativeRunTimeTicks       int64           `json:"CumulativeRunTimeTicks"`
	RunTimeTicks                 int64           `json:"RunTimeTicks"`
	PlayAccess                   string          `json:"PlayAccess"`
	ProductionYear               int             `json:"ProductionYear"`
	RemoteTrailers               []interface{}   `json:"RemoteTrailers"`
	ProviderIds                  json.RawMessage `json:"ProviderIds"`
	IsFolder                     bool            `json:"IsFolder"`
	ParentId                     string          `json:"ParentId"`
	Type                         string          `json:"Type"`
	People                       []interface{}   `json:"People"`
	Studios                      []*Subject      `json:"Studios"`
	GenreItems                   []*Subject      `json:"GenreItems"`
	LocalTrailerCount            int             `json:"LocalTrailerCount"`
	UserData                     *UserData       `json:"UserData"`
	RecursiveItemCount           int             `json:"RecursiveItemCount"`
	ChildCount                   int             `json:"ChildCount"`
	SpecialFeatureCount          int             `json:"SpecialFeatureCount"`
	DisplayPreferencesId         string          `json:"DisplayPreferencesId"`
	Tags                         []string        `json:"Tags"`
	PrimaryImageAspectRatio      float64         `json:"PrimaryImageAspectRatio"`
	Artists                      []interface{}   `json:"Artists"`
	ArtistItems                  []interface{}   `json:"ArtistItems"`
	AlbumArtists                 []interface{}   `json:"AlbumArtists"`
	AlbumArtist                  string          `json:"AlbumArtist"`
	ImageTags                    struct {
		Primary string `json:"Primary"`
	} `json:"ImageTags"`
	BackdropImageTags []interface{} `json:"BackdropImageTags"`
	ImageBlurHashes   struct {
		Primary map[string]string `json:"Primary"`
	} `json:"ImageBlurHashes"`
	LocationType string        `json:"LocationType"`
	LockedFields []interface{} `json:"LockedFields"`
	LockData     bool          `json:"LockData"`
}

type UpdateItemRequest struct {
	Id                           string          `json:"Id"`
	Name                         string          `json:"Name"`
	OriginalTitle                string          `json:"OriginalTitle"`
	ForcedSortName               string          `json:"ForcedSortName"`
	CommunityRating              string          `json:"CommunityRating,omitempty"`
	CriticRating                 string          `json:"CriticRating,omitempty"`
	IndexNumber                  string          `json:"IndexNumber,omitempty"`
	AirsBeforeSeasonNumber       string          `json:"AirsBeforeSeasonNumber,omitempty"`
	AirsAfterSeasonNumber        string          `json:"AirsAfterSeasonNumber,omitempty"`
	AirsBeforeEpisodeNumber      string          `json:"AirsBeforeEpisodeNumber,omitempty"`
	ParentIndexNumber            string          `json:"ParentIndexNumber,omitempty"`
	DisplayOrder                 string          `json:"DisplayOrder,omitempty"`
	Album                        string          `json:"Album"`
	AlbumArtists                 []*Subject      `json:"AlbumArtists"`
	ArtistItems                  []*Subject      `json:"ArtistItems"`
	Overview                     string          `json:"Overview,omitempty"`
	Status                       string          `json:"Status,omitempty"`
	AirDays                      []interface{}   `json:"AirDays,omitempty"`
	AirTime                      string          `json:"AirTime,omitempty"`
	Genres                       []string        `json:"Genres"`
	Tags                         []string        `json:"Tags"`
	Studios                      []*Subject      `json:"Studios"`
	PremiereDate                 time.Time       `json:"PremiereDate"`
	DateCreated                  time.Time       `json:"DateCreated"`
	EndDate                      string          `json:"EndDate,omitempty"`
	ProductionYear               string          `json:"ProductionYear"`
	AspectRatio                  string          `json:"AspectRatio,omitempty"`
	Video3DFormat                string          `json:"Video3DFormat,omitempty"`
	OfficialRating               string          `json:"OfficialRating,omitempty"`
	CustomRating                 string          `json:"CustomRating,omitempty"`
	People                       []*People       `json:"People"`
	LockData                     bool            `json:"LockData,omitempty"`
	LockedFields                 []interface{}   `json:"LockedFields"`
	ProviderIds                  json.RawMessage `json:"ProviderIds"`
	PreferredMetadataLanguage    string          `json:"PreferredMetadataLanguage,omitempty"`
	PreferredMetadataCountryCode string          `json:"PreferredMetadataCountryCode,omitempty"`
	Taglines                     []interface{}   `json:"Taglines"`
}

type People struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
	Role string `json:"Role"`
}

type ProviderIds struct {
}
