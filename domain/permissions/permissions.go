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
