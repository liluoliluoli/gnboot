package adaptor

import (
	"github.com/google/wire"
)

// ProviderSet is adaptor providers.
var ProviderSet = wire.NewSet(NewMovieProvider, NewEpisodeProvider, NewSeasonProvider, NewSeriesProvider, NewGenreProvider, NewStudioProvider, NewKeywordProvider, NewActorProvider, NewUserProvider)
