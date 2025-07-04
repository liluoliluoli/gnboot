package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api"
	"github.com/liluoliluoli/gnboot/api/video"
	"github.com/liluoliluoli/gnboot/internal/common/utils/page_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VideoProvider struct {
	video.UnimplementedVideoRemoteServiceServer
	video   *service.VideoService
	user    *service.UserService
	episode *service.EpisodeService
}

func NewVideoProvider(video *service.VideoService, user *service.UserService, episode *service.EpisodeService) *VideoProvider {
	return &VideoProvider{
		video:   video,
		user:    user,
		episode: episode,
	}
}

func (s *VideoProvider) CreateMovie(ctx context.Context, req *video.CreateVideoRequest) (*emptypb.Empty, error) {
	err := s.video.Create(ctx, (&sdomain.Video{}).ConvertFromDto(req))
	return &emptypb.Empty{}, err
}

func (s *VideoProvider) GetVideo(ctx context.Context, req *video.GetVideoRequest) (*video.Video, error) {
	rs, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.video.Get(ctx, int64(req.Id), rs.ID)
	if err != nil {
		return nil, err
	}
	go s.episode.TransferStoreNextEpisodeToAliyun(ctx, res.ID, res.LastPlayedEpisodeId)
	return res.ConvertToDto(), nil
}

func (s *VideoProvider) SearchVideo(ctx context.Context, req *video.SearchVideoRequest) (*video.SearchVideoResp, error) {
	condition := &sdomain.VideoSearch{
		Page:      page_util.ToDomainPage(req.Page),
		Type:      lo.FromPtr(req.VideoType),
		Search:    lo.FromPtr(req.Search),
		Sort:      lo.FromPtr(req.Sort),
		Genre:     lo.FromPtr(req.Genre),
		Region:    lo.FromPtr(req.Region),
		Year:      lo.FromPtr(req.Year),
		IsHistory: req.IsHistory,
	}
	rs, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.video.Page(ctx, condition, rs.ID)
	if err != nil {
		return nil, err
	}
	return &video.SearchVideoResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Video, index int) *video.Video {
			return item.ConvertToDto()
		}),
	}, nil
}

func (s *VideoProvider) PageFavorites(ctx context.Context, req *video.PageFavoritesRequest) (*video.SearchVideoResp, error) {
	rs, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.video.PageFavorites(ctx, rs.ID, page_util.ToDomainPage(req.Page))
	if err != nil {
		return nil, err
	}
	return &video.SearchVideoResp{
		Page: page_util.ToAdaptorPage(res.Page),
		List: lo.Map(res.List, func(item *sdomain.Video, index int) *video.Video {
			return item.ConvertToDto()
		}),
	}, nil
}

func (s *VideoProvider) UpdateMovie(ctx context.Context, req *video.UpdateVideoRequest) (*emptypb.Empty, error) {
	err := s.video.Update(ctx, (&sdomain.UpdateVideo{}).ConvertFromDto(req))
	return &emptypb.Empty{}, err
}

func (s *VideoProvider) DeleteMovie(ctx context.Context, req *api.IdsRequest) (*emptypb.Empty, error) {
	err := s.video.Delete(ctx, lo.Map(strings.Split(req.Ids, ","), func(item string, index int) int64 {
		id, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return 0
		}
		return id
	})...)
	return &emptypb.Empty{}, err
}
