package models

import "time"

// Role constants (obfuscated storage — short codes)
const (
	RoleUser       = "user"
	RoleSupervisor = "sup"
	RoleAdmin      = "adm"
	RoleSuperAdmin = "sadm"
)

// Room types
const (
	RoomDM      = "dm"
	RoomGroup   = "group"
	RoomPublic  = "public"
	RoomPrivate = "private"
)

// User represents a chat user
type User struct {
	ID            int64     `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	DisplayName   string    `json:"display_name" db:"display_name"`
	AvatarPath    string    `json:"avatar_path,omitempty" db:"avatar_path"`
	Role          string    `json:"role" db:"role"`
	Hidden        bool      `json:"-" db:"hidden"`
	Active        bool      `json:"active" db:"active"`
	MustChangePwd bool      `json:"must_change_pwd" db:"must_change_pwd"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	LastSeen      *time.Time `json:"last_seen,omitempty" db:"last_seen"`
}

// Sanitized returns a user safe for API responses (no sensitive fields)
func (u *User) Sanitized() *User {
	c := *u
	c.PasswordHash = ""
	c.Hidden = false
	return &c
}

// IsHidden reports whether this user should be invisible to others
func (u *User) IsHidden() bool {
	return u.Hidden
}

// CanManageRooms checks supervisor+ level
func (u *User) CanManageRooms() bool {
	return u.Role == RoleSupervisor || u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// CanManageUsers checks admin+ level
func (u *User) CanManageUsers() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// CanManageAdmins checks super admin only
func (u *User) CanManageAdmins() bool {
	return u.Role == RoleSuperAdmin
}

// Room represents a chat room
type Room struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Type      string    `json:"type" db:"type"`
	CreatorID int64     `json:"creator_id" db:"creator_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RoomMember represents membership of a user in a room
type RoomMember struct {
	ID       int64     `json:"id" db:"id"`
	RoomID   int64     `json:"room_id" db:"room_id"`
	UserID   int64     `json:"user_id" db:"user_id"`
	Role     string    `json:"role" db:"role"` // "owner" | "member"
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

// Message represents a chat message (body is encrypted at rest)
type Message struct {
	ID        int64      `json:"id" db:"id"`
	RoomID    int64      `json:"room_id" db:"room_id"`
	SenderID  int64      `json:"sender_id" db:"sender_id"`
	Body      []byte     `json:"body" db:"body"`
	IV        []byte     `json:"-" db:"iv"`
	FileID    *int64     `json:"file_id,omitempty" db:"file_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	EditedAt  *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// MessageView is the API-facing message with sender info
type MessageView struct {
	ID         int64      `json:"id"`
	RoomID     int64      `json:"room_id"`
	SenderID   int64      `json:"sender_id"`
	SenderName string     `json:"sender_name"`
	Text       string     `json:"text"`
	FileID     *int64     `json:"file_id,omitempty"`
	FileName   string     `json:"file_name,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	EditedAt   *time.Time `json:"edited_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Reactions  []Reaction `json:"reactions,omitempty"`
}

// Reaction represents an emoji reaction on a message
type Reaction struct {
	ID        int64     `json:"id" db:"id"`
	MessageID int64     `json:"message_id" db:"message_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Emoji     string    `json:"emoji" db:"emoji"`
}

// File represents an uploaded file
type File struct {
	ID        int64     `json:"id" db:"id"`
	OwnerID   int64     `json:"owner_id" db:"owner_id"`
	RoomID    int64     `json:"room_id" db:"room_id"`
	FileName  string    `json:"file_name" db:"file_name"`
	StoredAs  string    `json:"-" db:"stored_as"`
	Size      int64     `json:"size" db:"size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Session represents an active login session
type Session struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"-" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
