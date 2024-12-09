package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/api/series"
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
	res, err := s.series.Get(ctx, int64(req.Id), 1)
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
	res, err := s.series.Page(ctx, condition, 1)
	if err != nil {
		return nil, err
	}
	return &series.SearchSeriesResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Series, index int) *series.SeriesResp {
			return item.ConvertToDto()
		}),
	}, nil
}

func (s *SeriesProvider) FilterSeries(ctx context.Context, req *series.FilterSeriesRequest) (*series.SearchSeriesResp, error) {
	condition := &sdomain.SearchSeries{
		Page: page_util.ToDomainPage(req.Page),
		Id:   int64(req.Id),
		Type: req.Type,
	}
	res, err := s.series.Page(ctx, condition, 1)
	if err != nil {
		return nil, err
	}
	return &series.SearchSeriesResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Series, index int) *series.SeriesResp {
			return item.ConvertToDto()
		}),
	}, nil
}

func (s *SeriesProvider) NextToPlaySeries(ctx context.Context, req *series.NextToPlaySeriesRequest) (*series.NextToPlaySeriesResp, error) {
	condition := &sdomain.SearchSeries{
		Page:             page_util.ToDomainPage(req.Page),
		FilterByNextPlay: true,
	}
	res, err := s.series.Page(ctx, condition, 1)
	if err != nil {
		return nil, err
	}
	return &series.NextToPlaySeriesResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Series, index int) *episode.EpisodeResp {
			return item.Seasons[0].Episodes[0].ConvertToDto()
		}),
	}, nil
}
