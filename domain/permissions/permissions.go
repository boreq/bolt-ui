package permissions

import (
	"fmt"

	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

func CanListActivity(activity *domain.Activity, user *auth.ReadUser) (bool, error) {
	switch activity.Visibility() {
	case domain.PublicActivityVisibility:
		return true, nil
	case domain.UnlistedActivityVisibility:
		return user != nil && activity.UserUUID() == user.UUID, nil
	case domain.PrivateActivityVisibility:
		return user != nil && activity.UserUUID() == user.UUID, nil
	default:
		return false, fmt.Errorf("unsupported visibility '%s'", activity.Visibility())
	}
}

func CanViewActivity(activity *domain.Activity, user *auth.ReadUser) (bool, error) {
	switch activity.Visibility() {
	case domain.PublicActivityVisibility:
		return true, nil
	case domain.UnlistedActivityVisibility:
		return true, nil
	case domain.PrivateActivityVisibility:
		return user != nil && activity.UserUUID() == user.UUID, nil
	default:
		return false, fmt.Errorf("unsupported visibility '%s'", activity.Visibility())
	}
}

func CanEditActivity(activity *domain.Activity, user *auth.ReadUser) (bool, error) {
	return user != nil && activity.UserUUID() == user.UUID, nil
}

func CanDeleteActivity(activity *domain.Activity, user *auth.ReadUser) (bool, error) {
	return user != nil && activity.UserUUID() == user.UUID, nil
}

func CanViewPrivacyZone(privacyZone *domain.PrivacyZone, user *auth.ReadUser) (bool, error) {
	return user != nil && privacyZone.UserUUID() == user.UUID, nil
}

func CanListPrivacyZones(userUUID auth.UserUUID, user *auth.ReadUser) bool {
	return user != nil && userUUID == user.UUID
}
