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
	"github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/api/movie"
	"github.com/liluoliluoli/gnboot/api/season"
	"github.com/liluoliluoli/gnboot/api/series"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/api/user"
	"github.com/liluoliluoli/gnboot/internal/adaptor"
	"github.com/liluoliluoli/gnboot/internal/conf"
	localMiddleware "github.com/liluoliluoli/gnboot/internal/server/middleware"
	"golang.org/x/text/language"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Bootstrap,
	movieProvider *adaptor.MovieProvider,
	episodeProvider *adaptor.EpisodeProvider,
	seasonProvider *adaptor.SeasonProvider,
	seriesProvider *adaptor.SeriesProvider,
	genreProvider *adaptor.GenreProvider,
	studioProvider *adaptor.StudioProvider,
	keywordProvider *adaptor.KeywordProvider,
	actorProvider *adaptor.ActorProvider,
	userProvider *adaptor.UserProvider,
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
	movie.RegisterMovieRemoteServiceServer(srv, movieProvider)
	episode.RegisterEpisodeRemoteServiceServer(srv, episodeProvider)
	season.RegisterSeasonRemoteServiceServer(srv, seasonProvider)
	series.RegisterSeriesRemoteServiceServer(srv, seriesProvider)
	genre.RegisterGenreRemoteServiceServer(srv, genreProvider)
	studio.RegisterStudioRemoteServiceServer(srv, studioProvider)
	keyword.RegisterKeywordRemoteServiceServer(srv, keywordProvider)
	actor.RegisterActorRemoteServiceServer(srv, actorProvider)
	user.RegisterUserRemoteServiceServer(srv, userProvider)
	//TODO 追加业务注册
	return srv
}
