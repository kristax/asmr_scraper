package jellyfin

import "time"

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
	Name                         string        `json:"Name"`
	ServerId                     string        `json:"ServerId"`
	Id                           string        `json:"Id"`
	Etag                         string        `json:"Etag"`
	DateCreated                  time.Time     `json:"DateCreated"`
	CanDelete                    bool          `json:"CanDelete"`
	CanDownload                  bool          `json:"CanDownload"`
	PreferredMetadataLanguage    string        `json:"PreferredMetadataLanguage"`
	PreferredMetadataCountryCode string        `json:"PreferredMetadataCountryCode"`
	SortName                     string        `json:"SortName"`
	ForcedSortName               string        `json:"ForcedSortName"`
	PremiereDate                 time.Time     `json:"PremiereDate"`
	ExternalUrls                 []interface{} `json:"ExternalUrls"`
	Path                         string        `json:"Path"`
	EnableMediaSourceDisplay     bool          `json:"EnableMediaSourceDisplay"`
	CustomRating                 string        `json:"CustomRating"`
	ChannelId                    interface{}   `json:"ChannelId"`
	Overview                     string        `json:"Overview"`
	Taglines                     []interface{} `json:"Taglines"`
	Genres                       []string      `json:"Genres"`
	CommunityRating              float64       `json:"CommunityRating"`
	CumulativeRunTimeTicks       int64         `json:"CumulativeRunTimeTicks"`
	RunTimeTicks                 int64         `json:"RunTimeTicks"`
	PlayAccess                   string        `json:"PlayAccess"`
	ProductionYear               int           `json:"ProductionYear"`
	RemoteTrailers               []interface{} `json:"RemoteTrailers"`
	ProviderIds                  struct{}      `json:"ProviderIds"`
	IsFolder                     bool          `json:"IsFolder"`
	ParentId                     string        `json:"ParentId"`
	Type                         string        `json:"Type"`
	People                       []interface{} `json:"People"`
	Studios                      []*Subject    `json:"Studios"`
	GenreItems                   []*Subject    `json:"GenreItems"`
	LocalTrailerCount            int           `json:"LocalTrailerCount"`
	UserData                     *UserData     `json:"UserData"`
	RecursiveItemCount           int           `json:"RecursiveItemCount"`
	ChildCount                   int           `json:"ChildCount"`
	SpecialFeatureCount          int           `json:"SpecialFeatureCount"`
	DisplayPreferencesId         string        `json:"DisplayPreferencesId"`
	Tags                         []string      `json:"Tags"`
	PrimaryImageAspectRatio      float64       `json:"PrimaryImageAspectRatio"`
	Artists                      []interface{} `json:"Artists"`
	ArtistItems                  []interface{} `json:"ArtistItems"`
	AlbumArtists                 []interface{} `json:"AlbumArtists"`
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
	Id                           string        `json:"Id"`
	Name                         string        `json:"Name"`
	OriginalTitle                string        `json:"OriginalTitle"`
	ForcedSortName               string        `json:"ForcedSortName"`
	CommunityRating              string        `json:"CommunityRating"`
	CriticRating                 string        `json:"CriticRating"`
	IndexNumber                  interface{}   `json:"IndexNumber"`
	AirsBeforeSeasonNumber       string        `json:"AirsBeforeSeasonNumber"`
	AirsAfterSeasonNumber        string        `json:"AirsAfterSeasonNumber"`
	AirsBeforeEpisodeNumber      string        `json:"AirsBeforeEpisodeNumber"`
	ParentIndexNumber            interface{}   `json:"ParentIndexNumber"`
	DisplayOrder                 string        `json:"DisplayOrder"`
	Album                        string        `json:"Album"`
	AlbumArtists                 []*Subject    `json:"AlbumArtists"`
	ArtistItems                  []*Subject    `json:"ArtistItems"`
	Overview                     string        `json:"Overview"`
	Status                       string        `json:"Status"`
	AirDays                      []interface{} `json:"AirDays"`
	AirTime                      string        `json:"AirTime"`
	Genres                       []string      `json:"Genres"`
	Tags                         []string      `json:"Tags"`
	Studios                      []*Subject    `json:"Studios"`
	PremiereDate                 time.Time     `json:"PremiereDate"`
	DateCreated                  time.Time     `json:"DateCreated"`
	EndDate                      interface{}   `json:"EndDate"`
	ProductionYear               string        `json:"ProductionYear"`
	AspectRatio                  string        `json:"AspectRatio"`
	Video3DFormat                string        `json:"Video3DFormat"`
	OfficialRating               string        `json:"OfficialRating"`
	CustomRating                 string        `json:"CustomRating"`
	People                       []*People     `json:"People"`
	LockData                     bool          `json:"LockData"`
	LockedFields                 []interface{} `json:"LockedFields"`
	ProviderIds                  *ProviderIds  `json:"ProviderIds"`
	PreferredMetadataLanguage    string        `json:"PreferredMetadataLanguage"`
	PreferredMetadataCountryCode string        `json:"PreferredMetadataCountryCode"`
	Taglines                     []interface{} `json:"Taglines"`
}

type People struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
	Role string `json:"Role"`
}

type ProviderIds struct {
	Bangumi                 string `json:"Bangumi"`
	MusicBrainzAlbum        string `json:"MusicBrainzAlbum"`
	MusicBrainzArtist       string `json:"MusicBrainzArtist"`
	MusicBrainzReleaseGroup string `json:"MusicBrainzReleaseGroup"`
	AudioDbAlbum            string `json:"AudioDbAlbum"`
	AudioDbArtist           string `json:"AudioDbArtist"`
}
