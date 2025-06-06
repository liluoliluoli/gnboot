package server

import (
	"github.com/go-cinch/common/i18n"
	i18nMiddleware "github.com/go-cinch/common/middleware/i18n"
	"github.com/go-cinch/common/middleware/logging"
	tenantMiddleware "github.com/go-cinch/common/middleware/tenant"
	traceMiddleware "github.com/go-cinch/common/middleware/trace"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/liluoliluoli/gnboot/api/appversion"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/api/user"
	"github.com/liluoliluoli/gnboot/api/video"
	"github.com/liluoliluoli/gnboot/internal/adaptor"
	"github.com/liluoliluoli/gnboot/internal/conf"
	localMiddleware "github.com/liluoliluoli/gnboot/internal/server/middleware"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Bootstrap,
	videoProvider *adaptor.VideoProvider,
	episodeProvider *adaptor.EpisodeProvider,
	userProvider *adaptor.UserProvider,
	appVersionProvider *adaptor.AppVersionProvider,
	client redis.UniversalClient,
) *grpc.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		tenantMiddleware.Tenant(),
		ratelimit.Server(),
	}
	if c.Tracer.Enable {
		middlewares = append(middlewares, tracing.Server(), traceMiddleware.Id())
	}

	middlewares = append(
		middlewares,
		logging.Server(),
		i18nMiddleware.Translator(i18n.WithLanguage(language.Make(c.Server.Language)), i18n.WithFs(locales)),
		metadata.Server(),
	)
	if c.Server.Idempotent {
		middlewares = append(middlewares, localMiddleware.Idempotent())
	}
	if c.Server.Validate {
		middlewares = append(middlewares, validate.Validator())
	}
	middlewares = append(middlewares, localMiddleware.Auth(client))
	var opts = []grpc.ServerOption{grpc.Middleware(middlewares...)}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()))
	}
	opts = append(opts, grpc.UnaryInterceptor(localMiddleware.GrpcUnaryDisableTimeoutPropagation()))
	srv := grpc.NewServer(opts...)
	video.RegisterVideoRemoteServiceServer(srv, videoProvider)
	episode.RegisterEpisodeRemoteServiceServer(srv, episodeProvider)
	user.RegisterUserRemoteServiceServer(srv, userProvider)
	appversion.RegisterAppVersionRemoteServiceServer(srv, appVersionProvider)
	return srv
}
