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
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"github.com/liluoliluoli/gnboot/api/actor"
	"github.com/liluoliluoli/gnboot/api/episode"
	"github.com/liluoliluoli/gnboot/api/genre"
	"github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/api/movie"
	"github.com/liluoliluoli/gnboot/api/season"
	"github.com/liluoliluoli/gnboot/api/series"
	"github.com/liluoliluoli/gnboot/api/studio"
	"github.com/liluoliluoli/gnboot/internal/adaptor"
	"github.com/liluoliluoli/gnboot/internal/conf"
	localMiddleware "github.com/liluoliluoli/gnboot/internal/server/middleware"
	"golang.org/x/text/language"
	nethttp "net/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Bootstrap,
	movieProvider *adaptor.MovieProvider,
	episodeProvider *adaptor.EpisodeProvider,
	seasonProvider *adaptor.SeasonProvider,
	seriesProvider *adaptor.SeriesProvider,
	genreProvider *adaptor.GenreProvider,
	studioProvider *adaptor.StudioProvider,
	keywordProvider *adaptor.KeywordProvider,
	actorProvider *adaptor.ActorProvider,
) *http.Server {
	middlewares := []middleware.Middleware{
		recovery.Recovery(),
		tenantMiddleware.Tenant(),
		ratelimit.Server(),
		localMiddleware.Header(),
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
	var opts = []http.ServerOption{http.Middleware(middlewares...)}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	movie.RegisterMovieRemoteServiceHTTPServer(srv, movieProvider)
	episode.RegisterEpisodeRemoteServiceHTTPServer(srv, episodeProvider)
	season.RegisterSeasonRemoteServiceHTTPServer(srv, seasonProvider)
	series.RegisterSeriesRemoteServiceHTTPServer(srv, seriesProvider)
	genre.RegisterGenreRemoteServiceHTTPServer(srv, genreProvider)
	studio.RegisterStudioRemoteServiceHTTPServer(srv, studioProvider)
	keyword.RegisterKeywordRemoteServiceHTTPServer(srv, keywordProvider)
	actor.RegisterActorRemoteServiceHTTPServer(srv, actorProvider)
	//TODO 追加业务注册
	srv.HandlePrefix("/debug/pprof", pprof.NewHandler())
	srv.HandlePrefix("/pub/healthcheck", HealthHandler())
	return srv
}

func HealthHandler() *nethttp.ServeMux {
	mux := nethttp.NewServeMux()
	mux.HandleFunc("/pub/healthcheck", adaptor.HealthCheck)
	return mux
}
