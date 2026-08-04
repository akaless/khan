package ws

import "time"

// EventType constants (server → client)
const (
	EvMessage       = "message"
	EvMessageEdit   = "message_edited"
	EvMessageDelete = "message_deleted"
	EvReaction      = "reaction"
	EvTyping        = "typing"
	EvPresence      = "presence"
	EvUserAdded     = "user_added"
	EvUserRemoved   = "user_removed"
	EvRoleChanged   = "role_changed"
	EvUserDeleted   = "user_deleted"
	EvForceLogout   = "force_logout"
)

// Event is a server → client message
type Event struct {
	Type    string      `json:"type"`
	RoomID  int64       `json:"room_id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// NewEvent builds an event
func NewEvent(typ string, roomID int64, payload interface{}) Event {
	return Event{Type: typ, RoomID: roomID, Payload: payload}
}

// ClientEventType constants (client → server)
const (
	CevSendMessage = "send_message"
	CevEditMessage = "edit_message"
	CevDeleteMsg   = "delete_message"
	CevAddReaction = "add_reaction"
	CevDelReaction = "remove_reaction"
	CevTyping      = "typing"
	CevMarkRead    = "mark_read"
	CevJoinRoom    = "join_room"
)

// ClientEvent is a client → server message
type ClientEvent struct {
	Type     string `json:"type"`
	RoomID   int64  `json:"room_id"`
	MessageID int64 `json:"message_id,omitempty"`
	Text     string `json:"text,omitempty"`
	Emoji    string `json:"emoji,omitempty"`
	FileID   *int64 `json:"file_id,omitempty"`
	LastID   int64  `json:"last_message_id,omitempty"`
	UserID   int64  `json:"-"` // injected server-side
}

// MessagePayload for EvMessage
type MessagePayload struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Text       string    `json:"text"`
	FileID     *int64    `json:"file_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// PresencePayload for EvPresence
type PresencePayload struct {
	UserID int64  `json:"user_id"`
	Online bool   `json:"online"`
	Name   string `json:"name,omitempty"`
}

// TypingPayload for EvTyping
type TypingPayload struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
}
