package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type GenreProvider struct {
	genre.UnimplementedGenreRemoteServiceServer
	genre *service.GenreService
}

func NewGenreProvider(genre *service.GenreService) *GenreProvider {
	return &GenreProvider{genre: genre}
}

func (s *GenreProvider) FindGenre(ctx context.Context, req *genre.FindGenreRequest) (*genre.FindGenreResp, error) {
	res, err := s.genre.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return &genre.FindGenreResp{
		Genres: lo.Map(res, func(item *sdomain.Genre, index int) *genre.GenreResp {
			return &genre.GenreResp{
				Name: item.Name,
				Id:   item.ID,
			}
		}),
	}, nil

}
