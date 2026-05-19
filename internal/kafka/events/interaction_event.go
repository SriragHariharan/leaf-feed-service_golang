package events

type InteractionEvent struct {
	EventType    string `json:"eventType"`
	ActorUserId  string `json:"actorUserId"`
	TargetUserId string `json:"targetUserId"`
	PostId       string `json:"postId"`
	Timestamp    string `json:"timestamp"`
	CommentId    string `json:"commentId,omitempty"`
}
