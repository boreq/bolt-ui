package application

import (
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/music"
	"github.com/boreq/velo/application/queries"
)

type Application struct {
	Auth    auth.Auth
	Music   Music
	Queries Queries
}

type Music struct {
	Thumbnail *music.ThumbnailHandler
	Track     *music.TrackHandler
	Browse    *music.BrowseHandler
}

type Queries struct {
	Stats *queries.StatsHandler
}
