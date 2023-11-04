package model

import (
	"asmr_scraper/client/jellyfin"
	"fmt"
	"github.com/go-kid/ioc/util/fas"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
	"time"
)

type ProjectInfo struct {
	ItemId          string
	Code            string
	Path            string
	Name            string
	Name2           string
	Tags            []string
	ReleaseDate     time.Time
	CreateDate      time.Time
	Artists         []string
	Rating          float64
	Group           string
	Nsfw            bool
	Price           int
	Sales           int
	Overview        string
	PrimaryImageUrl string
}

func (p *ProjectInfo) ToJellyfinUpdateItemReq() *jellyfin.UpdateItemRequest {
	base := filepath.Base(p.Path)
	return &jellyfin.UpdateItemRequest{
		Id: p.ItemId,
		Name: func() string {
			builder := strings.Builder{}
			builder.WriteString(p.Code)
			builder.WriteString(" ")
			builder.WriteString(p.Name)
			if base != p.Code {
				builder.WriteString(fmt.Sprintf(" [%s]", base))
			}
			builder.WriteString(fmt.Sprintf(" Actors: %s", strings.Join(p.Artists, ",")))
			return builder.String()
		}(),
		OriginalTitle:           fas.TernaryOp(p.Name2 == "", p.Path, p.Name2),
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
		AlbumArtists: lo.Map[string, *jellyfin.Subject](p.Artists, func(item string, _ int) *jellyfin.Subject {
			return &jellyfin.Subject{Name: item}
		}),
		ArtistItems: []*jellyfin.Subject{{Name: p.Group}},
		//Overview:    fmt.Sprintf(p.OverviewTemplate, p.Group, p.Price, p.Sales, p.Code),
		Overview: p.Overview,
		Status:   "",
		AirDays:  []any{},
		AirTime:  "",
		Genres:   p.Tags,
		Tags:     []string{fas.TernaryOp(p.Nsfw, "R18", "全年龄")},
		Studios: []*jellyfin.Subject{
			{
				Name: p.Group,
			},
		},
		PremiereDate:                 p.ReleaseDate,
		DateCreated:                  p.CreateDate,
		EndDate:                      nil,
		ProductionYear:               fmt.Sprintf("%d", p.ReleaseDate.Year()),
		AspectRatio:                  "",
		Video3DFormat:                "",
		OfficialRating:               fas.TernaryOp(p.Nsfw, "XXX", "APPROVED"),
		CustomRating:                 "",
		People:                       []any{},
		LockData:                     true,
		LockedFields:                 []any{},
		ProviderIds:                  &jellyfin.ProviderIds{},
		PreferredMetadataLanguage:    "",
		PreferredMetadataCountryCode: "",
		Taglines:                     []any{},
	}
}
