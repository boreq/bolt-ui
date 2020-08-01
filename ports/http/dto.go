package http

type UserProfile struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type PostActivityResponse struct {
	ActivityUUID string `json:"activityUUID"`
}
