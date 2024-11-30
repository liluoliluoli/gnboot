package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/api/series"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/api/subtitle"
	"github.com/liluoliluoli/gnboot/internal/common/utils/page_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type SeriesProvider struct {
	series.UnimplementedSeriesRemoteServiceServer
	series *service.SeriesService
}

func NewSeriesProvider(series *service.SeriesService) *SeriesProvider {
	return &SeriesProvider{series: series}
}

func (s *SeriesProvider) GetSeries(ctx context.Context, req *series.GetSeriesRequest) (*series.SeriesResp, error) {
	res, err := s.series.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return res.ConvertToDto(), nil
}

func (s *SeriesProvider) FindSeries(ctx context.Context, req *series.FindSeriesRequest) (*series.SearchSeriesResp, error) {
	condition := &sdomain.SearchSeries{
		Page:   page_util.ToDomainPage(req.Page),
		Search: lo.FromPtr(req.Search),
		Sort: &sdomain.Sort{
			Filter: lo.TernaryF(req.Sort != nil, func() string {
				return lo.FromPtr(req.Sort.Filter)
			}, func() string {
				return ""
			}),
			Type: lo.TernaryF(req.Sort != nil, func() string {
				return lo.FromPtr(req.Sort.Type)
			}, func() string {
				return ""
			}),
			Direction: lo.TernaryF(req.Sort != nil, func() string {
				return lo.FromPtr(req.Sort.Direction)
			}, func() string {
				return ""
			}),
		},
	}
	res, err := s.series.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &series.SearchSeriesResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Series, index int) *series.SeriesResp {
			return &series.SeriesResp{
				Id:            item.ID,
				OriginalTitle: item.OriginalTitle,
				Status:        item.Status,
				VoteAverage:   item.VoteAverage,
				VoteCount:     item.VoteCount,
				Country:       item.Country,
				Trailer:       item.Trailer,
				Url:           item.URL,
				Downloaded:    item.Downloaded,
				FileSize:      item.FileSize,
				Filename:      item.Filename,
				Ext:           item.Ext,
				Genres: lo.Map(item.Genres, func(item *sdomain.Genre, index int) *genre.GenreResp {
					return item.ConvertToDto()
				}),
				Studios: lo.Map(item.Studios, func(item *sdomain.Studio, index int) *studio.StudioResp {
					return item.ConvertToDto()
				}),
				Keywords: lo.Map(item.Keywords, func(item *sdomain.Keyword, index int) *keyword.KeywordResp {
					return item.ConvertToDto()
				}),
				Actors: lo.Map(item.Actors, func(item *sdomain.Actor, index int) *actor.ActorResp {
					return item.ConvertToDto()
				}),
				Subtitles: lo.Map(item.Subtitles, func(item *sdomain.VideoSubtitleMapping, index int) *subtitle.SubtitleResp {
					return item.ConvertToDto()
				}),
			}
		}),
	}, nil
}

func (s *SeriesProvider) FilterSeries(ctx context.Context, req *series.FilterSeriesRequest) (*series.SearchSeriesResp, error) {
	condition := &sdomain.SearchSeries{
		Page: page_util.ToDomainPage(req.Page),
		Id:   req.Id,
		Type: req.Type,
	}
	res, err := s.series.Page(ctx, condition)
	if err != nil {
		return nil, err
	}
	return &series.SearchSeriesResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Series, index int) *series.SeriesResp {
			return &series.SeriesResp{
				Id:            item.ID,
				OriginalTitle: item.OriginalTitle,
				Status:        item.Status,
				VoteAverage:   item.VoteAverage,
				VoteCount:     item.VoteCount,
				Country:       item.Country,
				Trailer:       item.Trailer,
				Url:           item.URL,
				Downloaded:    item.Downloaded,
				FileSize:      item.FileSize,
				Filename:      item.Filename,
				Ext:           item.Ext,
				Genres: lo.Map(item.Genres, func(item *sdomain.Genre, index int) *genre.GenreResp {
					return item.ConvertToDto()
				}),
				Studios: lo.Map(item.Studios, func(item *sdomain.Studio, index int) *studio.StudioResp {
					return item.ConvertToDto()
				}),
				Keywords: lo.Map(item.Keywords, func(item *sdomain.Keyword, index int) *keyword.KeywordResp {
					return item.ConvertToDto()
				}),
				Actors: lo.Map(item.Actors, func(item *sdomain.Actor, index int) *actor.ActorResp {
					return item.ConvertToDto()
				}),
				Subtitles: lo.Map(item.Subtitles, func(item *sdomain.VideoSubtitleMapping, index int) *subtitle.SubtitleResp {
					return item.ConvertToDto()
				}),
			}
		}),
	}, nil
}
