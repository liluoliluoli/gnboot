package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/appversion"
	"github.com/liluoliluoli/gnboot/internal/service"
)

type AppVersionProvider struct {
	appversion.UnimplementedAppVersionRemoteServiceServer
	appVersion *service.AppVersionService
}

func NewAppVersionProvider(appVersion *service.AppVersionService) *AppVersionProvider {
	return &AppVersionProvider{appVersion: appVersion}
}

func (s *AppVersionProvider) GetLastVersion(ctx context.Context, req *appversion.GetLastAppVersionRequest) (*appversion.AppVersion, error) {
	res, err := s.appVersion.GetLastVersion(ctx)
	if err != nil {
		return nil, err
	}
	return res.ConvertToDto(), nil
}
