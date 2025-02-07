package model

import (
	"asmr_scraper/client/jellyfin"
	"fmt"
	"github.com/go-kid/ioc/util/fas"
	"github.com/samber/lo"
	"time"
)

type ProjectInfoData interface {
	ToProjectInfo(code, path string, item *jellyfin.ItemInfoResponse) (*ProjectInfo, error)
}

type ProjectInfo struct {
	ItemId          string
	Code            string
	Name            string
	Name2           string
	Tags            []string
	ReleaseDate     time.Time
	CreateDate      time.Time
	Artist          []string
	People          []*People
	Rating          float64
	Group           string
	Nsfw            bool
	Price           int
	Sales           int
	Overview        string
	PrimaryImageUrl string
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
	return &jellyfin.UpdateItemRequest{
		Id:                      p.ItemId,
		Name:                    p.Name,
		OriginalTitle:           p.Name2,
		ForcedSortName:          p.Code,
		CommunityRating:         fmt.Sprintf("%.2f", p.Rating),
		CriticRating:            "",
		IndexNumber:             nil,
		AirsBeforeSeasonNumber:  "",
		AirsAfterSeasonNumber:   "",
		AirsBeforeEpisodeNumber: "",
		ParentIndexNumber:       nil,
		DisplayOrder:            "",
		Album:                   p.Code,
		AlbumArtists:            artists,
		ArtistItems:             append(artists, fas.TernaryOp(p.Group == "", nil, []*jellyfin.Subject{{Name: p.Group}})...),
		Overview:                p.Overview,
		Status:                  "",
		AirDays:                 []any{},
		AirTime:                 "",
		Genres:                  p.Tags,
		Tags:                    []string{fas.TernaryOp(p.Nsfw, "R18", "全年龄")},
		Studios:                 fas.TernaryOp(p.Group == "", nil, []*jellyfin.Subject{{Name: p.Group}}),
		PremiereDate:            p.ReleaseDate,
		DateCreated:             p.CreateDate,
		EndDate:                 nil,
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
		ProviderIds:                  &jellyfin.ProviderIds{},
		PreferredMetadataLanguage:    "",
		PreferredMetadataCountryCode: "",
		Taglines:                     []any{},
	}
}
