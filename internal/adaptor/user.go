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
	err := s.user.UpdateFavorite(ctx, 1, int64(req.Id), req.Type, req.Favorite)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserProvider) UpdatePlayedStatus(ctx context.Context, req *user.UpdatePlayedStatusRequest) (*emptypb.Empty, error) {
	err := s.user.UpdatePlayStatus(ctx, 1, int64(req.Id), req.Type, req.Position)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
