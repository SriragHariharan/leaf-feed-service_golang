package events

type UserEvent struct {
	UserID         string  `json:"userID"`
	Username       string  `json:"username"`
	ProfilePicture *string `json:"profilePicture"`
}
