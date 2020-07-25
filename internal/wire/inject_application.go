package wire

import (
	authAdapters "github.com/boreq/velo/adapters/auth"
	trackerAdapters "github.com/boreq/velo/adapters/tracker"
	"github.com/boreq/velo/application"
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/music"
	"github.com/boreq/velo/application/queries"
	"github.com/boreq/velo/application/tracker"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var appSet = wire.NewSet(
	wire.Struct(new(application.Application), "*"),

	wire.Struct(new(auth.Auth), "*"),
	auth.NewRegisterInitialHandler,
	auth.NewLoginHandler,
	auth.NewLogoutHandler,
	auth.NewCheckAccessTokenHandler,
	auth.NewListHandler,
	auth.NewCreateInvitationHandler,
	auth.NewRegisterHandler,
	auth.NewRemoveHandler,
	auth.NewSetPasswordHandler,

	wire.Struct(new(application.Music), "*"),
	music.NewTrackHandler,
	music.NewThumbnailHandler,
	music.NewBrowseHandler,

	wire.Struct(new(application.Queries), "*"),
	queries.NewStatsHandler,

	authAdapters.NewAuthTransactionProvider,
	wire.Bind(new(auth.TransactionProvider), new(*authAdapters.AuthTransactionProvider)),

	authAdapters.NewQueryTransactionProvider,
	wire.Bind(new(queries.TransactionProvider), new(*authAdapters.QueryTransactionProvider)),

	wire.Struct(new(auth.TransactableRepositories), "*"),
	wire.Struct(new(queries.TransactableRepositories), "*"),

	newQueryRepositoriesProvider,
	wire.Bind(new(authAdapters.QueryRepositoriesProvider), new(*queryRepositoriesProvider)),

	newAuthRepositoriesProvider,
	wire.Bind(new(authAdapters.AuthRepositoriesProvider), new(*authRepositoriesProvider)),

	wire.Bind(new(queries.UserRepository), new(*authAdapters.UserRepository)),
	wire.Bind(new(auth.UserRepository), new(*authAdapters.UserRepository)),
	authAdapters.NewUserRepository,

	wire.Bind(new(auth.InvitationRepository), new(*authAdapters.InvitationRepository)),
	authAdapters.NewInvitationRepository,

	wire.Bind(new(auth.PasswordHasher), new(*authAdapters.BcryptPasswordHasher)),
	authAdapters.NewBcryptPasswordHasher,

	wire.Bind(new(auth.AccessTokenGenerator), new(*authAdapters.CryptoAccessTokenGenerator)),
	authAdapters.NewCryptoAccessTokenGenerator,

	authAdapters.NewCryptoStringGenerator,
	wire.Bind(new(auth.CryptoStringGenerator), new(*authAdapters.CryptoStringGenerator)),
)

type authRepositoriesProvider struct {
}

func newAuthRepositoriesProvider() *authRepositoriesProvider {
	return &authRepositoriesProvider{}
}

func (p *authRepositoriesProvider) Provide(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	return BuildTransactableAuthRepositories(tx)
}

type queryRepositoriesProvider struct {
}

func newQueryRepositoriesProvider() *queryRepositoriesProvider {
	return &queryRepositoriesProvider{}
}

func (p *queryRepositoriesProvider) Provide(tx *bolt.Tx) (*queries.TransactableRepositories, error) {
	return BuildTransactableQueryRepositories(tx)
}

//lint:ignore U1000 because
var trackerSet = wire.NewSet(
	wire.Struct(new(tracker.Tracker), "*"),
	tracker.NewAddActivityHandler,

	trackerAdapters.NewTrackerTransactionProvider,
	wire.Bind(new(tracker.TransactionProvider), new(*trackerAdapters.TrackerTransactionProvider)),

	wire.Struct(new(tracker.TransactableRepositories), "*"),

	newTrackerRepositoriesProvider,
	wire.Bind(new(trackerAdapters.TrackerRepositoriesProvider), new(*trackerRepositoriesProvider)),
)

type trackerRepositoriesProvider struct {
}

func newTrackerRepositoriesProvider() *trackerRepositoriesProvider {
	return &trackerRepositoriesProvider{}
}

func (p *trackerRepositoriesProvider) Provide(tx *bolt.Tx) (*tracker.TransactableRepositories, error) {
	return BuildTransactableTrackerRepositories(tx)
}
