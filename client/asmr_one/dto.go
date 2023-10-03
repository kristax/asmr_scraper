package asmr_one

type RateCountDetail struct {
	ReviewPoint int `json:"review_point"`
	Count       int `json:"count"`
	Ratio       int `json:"ratio"`
}

type Rank struct {
	Term     string `json:"term"`
	Category string `json:"category"`
	Rank     int    `json:"rank"`
	RankDate string `json:"rank_date"`
}

type Vas struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	Id   int              `json:"id"`
	I18N map[string]*I18N `json:"i18n"`
	Name string           `json:"name"`
}

type I18N struct {
	Name string `json:"name"`
}

type Circle struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type WorkInfoResponse struct {
	Id                int                `json:"id"`
	Title             string             `json:"title"`
	CircleId          int                `json:"circle_id"`
	Name              string             `json:"name"`
	Nsfw              bool               `json:"nsfw"`
	Release           string             `json:"release"`
	DlCount           int                `json:"dl_count"`
	Price             int                `json:"price"`
	ReviewCount       int                `json:"review_count"`
	RateCount         int                `json:"rate_count"`
	RateAverage2Dp    float64            `json:"rate_average_2dp"`
	RateCountDetail   []*RateCountDetail `json:"rate_count_detail"`
	HasSubtitle       bool               `json:"has_subtitle"`
	CreateDate        string             `json:"create_date"`
	Vas               []*Vas             `json:"vas"`
	Tags              []*Tag             `json:"tags"`
	OriginalWorkno    interface{}        `json:"original_workno"`
	Circle            *Circle            `json:"circle"`
	SamCoverUrl       string             `json:"samCoverUrl"`
	ThumbnailCoverUrl string             `json:"thumbnailCoverUrl"`
	MainCoverUrl      string             `json:"mainCoverUrl"`
}
