package domain

type Activity struct {
	uuid     ActivityUUID
	userUUID UserUUID
	route    Route
}
