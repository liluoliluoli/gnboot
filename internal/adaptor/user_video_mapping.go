package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/user_video_mapping"
	"github.com/liluoliluoli/gnboot/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserVideoMappingProvider struct {
	user_video_mapping.UnimplementedUserVideoMappingRemoteServiceServer
	user *service.UserService
}

func NewUserVideoMappingProvider(user *service.UserService) *UserVideoMappingProvider {
	return &UserVideoMappingProvider{user: user}
}

func (s *UserVideoMappingProvider) UpdateFavorite(ctx context.Context, req *user_video_mapping.UpdateFavoriteRequest) (*emptypb.Empty, error) {
	err := s.user.UpdateFavorite(ctx, 1, req.Id, req.Type, req.Favorite)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserVideoMappingProvider) UpdatePlayedStatus(ctx context.Context, req *user_video_mapping.UpdatePlayedStatusRequest) (*emptypb.Empty, error) {
	err := s.user.UpdatePlayStatus(ctx, 1, req.Id, req.Type, req.Position)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
