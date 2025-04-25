package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/user"
	"github.com/liluoliluoli/gnboot/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserProvider struct {
	user.UnimplementedUserRemoteServiceServer
	user *service.UserService
}

func NewUserProvider(user *service.UserService) *UserProvider {
	return &UserProvider{user: user}
}

func (s *UserProvider) UpdateFavorite(ctx context.Context, req *user.UpdateFavoriteRequest) (*emptypb.Empty, error) {
	err := s.user.UpdateFavorite(ctx, int64(req.UserId), int64(req.VideoId), req.Favorite)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserProvider) UpdatePlayedStatus(ctx context.Context, req *user.UpdatePlayedStatusRequest) (*emptypb.Empty, error) {
	err := s.user.UpdatePlayStatus(ctx, int64(req.UserId), int64(req.VideoId), int64(req.EpisodeId), int64(req.Position))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
