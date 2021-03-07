package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

type AddPrivacyZone struct {
	UserUUID auth.UserUUID
	Position domain.Position
	Circle   domain.Circle
	Name     domain.PrivacyZoneName
}

type AddPrivacyZoneHandler struct {
	transactionProvider TransactionProvider
	routeFileParser     RouteFileParser
	uuidGenerator       UUIDGenerator
}

func NewAddPrivacyZoneHandler(transactionProvider TransactionProvider, routeFileParser RouteFileParser, uuidGenerator UUIDGenerator) *AddPrivacyZoneHandler {
	return &AddPrivacyZoneHandler{
		transactionProvider: transactionProvider,
		routeFileParser:     routeFileParser,
		uuidGenerator:       uuidGenerator,
	}
}

func (h *AddPrivacyZoneHandler) Execute(cmd AddPrivacyZone) (domain.PrivacyZoneUUID, error) {
	zone, err := h.createPrivacyZone(cmd)
	if err != nil {
		return domain.PrivacyZoneUUID{}, errors.Wrap(err, "could not create a privacy zone")
	}

	if err := h.transactionProvider.Write(func(repositories *TransactableRepositories) error {
		if err := repositories.PrivacyZone.Save(zone); err != nil {
			return errors.Wrap(err, "could not save the privacy zone")
		}

		if err := repositories.UserToPrivacyZone.Assign(zone.UserUUID(), zone.UUID()); err != nil {
			return errors.Wrap(err, "could not map a privacy zone to a user")
		}

		return nil
	}); err != nil {
		return domain.PrivacyZoneUUID{}, errors.Wrap(err, "transaction failed")
	}

	return zone.UUID(), nil
}

func (h *AddPrivacyZoneHandler) createPrivacyZone(cmd AddPrivacyZone) (*domain.PrivacyZone, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	privacyZoneUUID, err := domain.NewPrivacyZoneUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create an activity uuid")
	}

	return domain.NewPrivacyZone(
		privacyZoneUUID,
		cmd.UserUUID,
		cmd.Position,
		cmd.Circle,
		cmd.Name,
	)
}
