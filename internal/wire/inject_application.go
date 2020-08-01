package wire

import (
	authAdapters "github.com/boreq/velo/adapters/auth"
	"github.com/boreq/velo/application"
	"github.com/boreq/velo/application/auth"
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
	auth.NewGetUserHandler,

	authAdapters.NewAuthTransactionProvider,
	wire.Bind(new(auth.TransactionProvider), new(*authAdapters.AuthTransactionProvider)),

	wire.Struct(new(auth.TransactableRepositories), "*"),

	newAuthRepositoriesProvider,
	wire.Bind(new(authAdapters.AuthRepositoriesProvider), new(*authRepositoriesProvider)),

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

	trackerSet,
)

type authRepositoriesProvider struct {
}

func newAuthRepositoriesProvider() *authRepositoriesProvider {
	return &authRepositoriesProvider{}
}

func (p *authRepositoriesProvider) Provide(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	return BuildTransactableAuthRepositories(tx)
}

//lint:ignore U1000 because
var trackerSet = wire.NewSet(
	wire.Struct(new(tracker.Tracker), "*"),
	tracker.NewAddActivityHandler,
	tracker.NewGetActivityHandler,

	trackerRepositoriesSet,
	trackerTransactableRepositoriesSet,
)
