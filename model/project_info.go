package model

import (
	"asmr_scraper/client/jellyfin"
	"encoding/json"
	"fmt"
	"github.com/go-kid/ioc/util/fas"
	"github.com/samber/lo"
	"time"
)

type ProjectInfoData interface {
	ToProjectInfo(code, path string, item *jellyfin.ItemInfoResponse, subItems []*jellyfin.ItemInfoResponse) (*ProjectInfo, error)
}

type ProjectInfo struct {
	ItemId          string
	Code            string
	Name            string
	Name2           string
	SortName        string
	Index           string
	ParentIndex     string
	Tags            []string
	ReleaseDate     time.Time
	CreateDate      time.Time
	Artist          []string
	People          []*People
	Rating          float64
	Group           []string
	Nsfw            bool
	Price           int
	Sales           int
	Overview        string
	PrimaryImageUrl string
	ProviderIds     json.RawMessage
	ItemsInfo       []*ProjectInfo
}

type People struct {
	Name     string
	Type     PeopleType
	Role     string
	Gender   string
	HomePage string
}

type PeopleType string

const (
	TypeActor     PeopleType = "Actor"
	TypeComposer  PeopleType = "Composer"
	TypeDirector  PeopleType = "Director"
	TypeGuestStar PeopleType = "GuestStar"
	TypeProducer  PeopleType = "Producer"
	TypeWriter    PeopleType = "Writer"
)

func (p *ProjectInfo) ToJellyfinUpdateItemReq() *jellyfin.UpdateItemRequest {
	artists := lo.Map(p.Artist, func(item string, _ int) *jellyfin.Subject {
		return &jellyfin.Subject{Name: item}
	})
	groups := lo.Map(p.Group, func(item string, _ int) *jellyfin.Subject {
		return &jellyfin.Subject{Name: item}
	})
	return &jellyfin.UpdateItemRequest{
		Id:                      p.ItemId,
		Name:                    p.Name,
		OriginalTitle:           p.Name2,
		ForcedSortName:          lo.If(p.SortName != "", p.SortName).Else(p.Code),
		CommunityRating:         fmt.Sprintf("%.2f", p.Rating),
		CriticRating:            "",
		IndexNumber:             p.Index,
		AirsBeforeSeasonNumber:  "",
		AirsAfterSeasonNumber:   "",
		AirsBeforeEpisodeNumber: "",
		ParentIndexNumber:       p.ParentIndex,
		DisplayOrder:            "",
		Album:                   p.Code,
		AlbumArtists:            artists,
		ArtistItems:             append(artists, fas.TernaryOp(len(groups) == 0, nil, groups)...),
		Overview:                p.Overview,
		Status:                  "",
		AirDays:                 []any{},
		AirTime:                 "",
		Genres:                  p.Tags,
		Tags:                    []string{fas.TernaryOp(p.Nsfw, "R18", "全年龄")},
		Studios:                 groups,
		PremiereDate:            p.ReleaseDate,
		DateCreated:             p.CreateDate,
		EndDate:                 "",
		ProductionYear:          fmt.Sprintf("%d", p.ReleaseDate.Year()),
		AspectRatio:             "",
		Video3DFormat:           "",
		OfficialRating:          fas.TernaryOp(p.Nsfw, "XXX", "APPROVED"),
		CustomRating:            "",
		People: lo.Map(p.People, func(item *People, index int) *jellyfin.People {
			return &jellyfin.People{
				Name: item.Name,
				Type: string(item.Type),
				Role: item.Role,
			}
		}),
		LockData:                     false,
		LockedFields:                 []any{},
		ProviderIds:                  lo.If(p.ProviderIds != nil, p.ProviderIds).Else(json.RawMessage("{}")),
		PreferredMetadataLanguage:    "",
		PreferredMetadataCountryCode: "",
		Taglines:                     []any{},
	}
}
