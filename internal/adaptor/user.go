package adaptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/liluoliluoli/gnboot/api/user"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/gerror"
	"github.com/liluoliluoli/gnboot/internal/common/utils/context_util"
	jwtutil "github.com/liluoliluoli/gnboot/internal/common/utils/jwt_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/security_util"
	"github.com/liluoliluoli/gnboot/internal/common/utils/time_util"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type UserProvider struct {
	user.UnimplementedUserRemoteServiceServer
	user   *service.UserService
	client redis.UniversalClient
	video  *service.VideoService
}

func NewUserProvider(user *service.UserService, client redis.UniversalClient, video *service.VideoService) *UserProvider {
	return &UserProvider{
		user:   user,
		client: client,
		video:  video,
	}
}

func (s *UserProvider) UpdateFavorite(ctx context.Context, req *user.UpdateFavoriteRequest) (*emptypb.Empty, error) {
	rs, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	err = s.user.UpdateFavorite(ctx, rs.ID, int64(req.VideoId), req.Favorite)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserProvider) UpdatePlayedStatus(ctx context.Context, req *user.UpdatePlayedStatusRequest) (*emptypb.Empty, error) {
	rs, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	for _, updatePlayedStatus := range req.UpdatePlayedStatusList {
		err = s.user.UpdatePlayStatus(ctx, rs.ID, int64(updatePlayedStatus.VideoId), int64(updatePlayedStatus.EpisodeId), int64(updatePlayedStatus.Position), int64(updatePlayedStatus.PlayTimestamp))
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *UserProvider) Create(ctx context.Context, req *user.CreateUserRequest) (*emptypb.Empty, error) {
	if req.Password != req.ConfirmPassword {
		return &emptypb.Empty{}, errors.New("两次密码不一致")
	}
	err := s.user.Create(ctx, req.UserName, security_util.SignByHMACSha256(req.Password, constant.SYS_PWD))
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserProvider) Login(ctx context.Context, req *user.LoginUserRequest) (*user.LoginUserResp, error) {
	rs, err := s.user.QueryByUserName(ctx, req.UserName)
	if err != nil {
		return nil, err
	}
	if rs == nil {
		return nil, errors.New("用户不存在")
	}
	if rs.Password != security_util.SignByHMACSha256(req.Password, constant.SYS_PWD) {
		return nil, errors.New("密码错误")
	}
	authorization, err := jwtutil.GenerateUserToken(&jwtutil.UserClaims{
		UserName: req.UserName,
	}, constant.SYS_PWD)
	if err != nil {
		return nil, err
	}
	err = s.client.Set(ctx, fmt.Sprintf(constant.RK_UserTokenPrefix, authorization), req.UserName, 0).Err()
	if err != nil {
		return nil, err
	}
	rs.SessionToken = lo.ToPtr(fmt.Sprintf(constant.RK_UserTokenPrefix, authorization))
	err = s.user.UpdateSessionToken(ctx, rs)
	if err != nil {
		return nil, err
	}
	return &user.LoginUserResp{
		Authorization: fmt.Sprintf(constant.RK_UserTokenPrefix, authorization),
	}, nil
}

func (s *UserProvider) Logout(ctx context.Context, req *user.LogoutUserRequest) (*user.LogoutUserResp, error) {
	rs, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	sessionToken, err := context_util.GetGenericContext[string](ctx, constant.CTX_SessionToken)
	if err != nil {
		return nil, err
	}
	err = s.client.Del(ctx, sessionToken).Err()
	if err != nil {
		return nil, err
	}
	rs.SessionToken = lo.ToPtr("-1")
	err = s.user.UpdateSessionToken(ctx, rs)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *UserProvider) GetCurrentWatchCount(ctx context.Context, req *user.GetCurrentWatchCountRequest) (*user.GetCurrentWatchCountResp, error) {
	userName, err := context_util.GetGenericContext[string](ctx, constant.CTX_UserName)
	if err != nil {
		return nil, err
	}
	currentWatchs, err := s.client.HGet(ctx, fmt.Sprintf(constant.RK_UserWatchCountPrefix, userName), time_util.FormatYYYYMMDD(time.Now())).Int()
	if gerror.HandleRedisNotFoundError(err) != nil {
		return nil, err
	}

	return &user.GetCurrentWatchCountResp{
		WatchCount: int32(currentWatchs),
	}, nil
}

func (s *UserProvider) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	currentUser, err := s.user.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	currentWatches, err := s.client.HGet(ctx, fmt.Sprintf(constant.RK_UserWatchCountPrefix, currentUser.UserName), time_util.FormatYYYYMMDD(time.Now())).Int()
	if gerror.HandleRedisNotFoundError(err) != nil {
		return nil, err
	}
	res, err := s.video.PageFavorites(ctx, currentUser.ID, &sdomain.Page{
		CurrentPage: 1,
		PageSize:    100,
	})
	noticeTitle := s.client.HGet(ctx, constant.RK_Notice, constant.HK_NoticeTitle).Val()
	noticeContent := s.client.HGet(ctx, constant.RK_Notice, constant.HK_NoticeContent).Val()
	return &user.User{
		WatchCount:     int32(currentWatches),
		RestWatchCount: int32(constant.MaxWatchCountByDay - currentWatches),
		UserName:       currentUser.UserName,
		FavoriteCount:  int32(res.Page.Count),
		PackageType:    currentUser.PackageType,
		PackageExpiredTime: lo.TernaryF(currentUser.PackageExpiredTime != nil, func() *int32 {
			return lo.ToPtr(int32(currentUser.PackageExpiredTime.Unix()))
		}, func() *int32 {
			return nil
		}),
		NoticeTitle:   gerror.HandleRedisStringNotFound(noticeTitle),
		NoticeContent: gerror.HandleRedisStringNotFound(noticeContent),
	}, nil
}

func (s *UserProvider) UpdateNotice(ctx context.Context, req *user.UpdateNoticeRequest) (*emptypb.Empty, error) {
	err := s.user.UpdateNotice(ctx, req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
