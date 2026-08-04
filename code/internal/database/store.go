package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store is a lightweight JSON-file datastore (replaces SQLite for 20-user LAN scale).
// - One file per table, atomic writes via temp+rename
// - In-memory index + mutex
// - Auto-increment IDs
type Store struct {
	dir  string
	mu   sync.RWMutex
	data *Data
}

// Data holds all tables in memory, persisted to disk
type Data struct {
	Users       []UserRecord       `json:"users"`
	Rooms       []RoomRecord       `json:"rooms"`
	RoomMembers []RoomMemberRecord `json:"room_members"`
	Messages    []MessageRecord    `json:"messages"`
	Reactions   []ReactionRecord   `json:"reactions"`
	Sessions    []SessionRecord    `json:"sessions"`
	Files       []FileRecord       `json:"files"`
	Seq         map[string]int64   `json:"seq"` // per-table id counters
}

// UserRecord mirrors the users table
type UserRecord struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	PasswordHash  string `json:"password_hash"`
	DisplayName   string `json:"display_name"`
	AvatarPath    string `json:"avatar_path,omitempty"`
	Role          string `json:"role"`
	Hidden        bool   `json:"hidden"`
	Active        bool   `json:"active"`
	MustChangePwd bool   `json:"must_change_pwd"`
	CreatedAt     string `json:"created_at"`
	LastSeen      string `json:"last_seen,omitempty"`
}

type RoomRecord struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatorID int64  `json:"creator_id"`
	CreatedAt string `json:"created_at"`
}

type RoomMemberRecord struct {
	ID       int64  `json:"id"`
	RoomID   int64  `json:"room_id"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

type MessageRecord struct {
	ID        int64  `json:"id"`
	RoomID    int64  `json:"room_id"`
	SenderID  int64  `json:"sender_id"` // -1 = deleted user
	Body      string `json:"body"`      // base64 AES-GCM
	FileID    *int64 `json:"file_id,omitempty"`
	CreatedAt string `json:"created_at"`
	EditedAt  string `json:"edited_at,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
}

type ReactionRecord struct {
	ID        int64  `json:"id"`
	MessageID int64  `json:"message_id"`
	UserID    int64  `json:"user_id"`
	Emoji     string `json:"emoji"`
}

type SessionRecord struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

type FileRecord struct {
	ID        int64  `json:"id"`
	OwnerID   int64  `json:"owner_id"`
	RoomID    int64  `json:"room_id"`
	FileName  string `json:"file_name"`
	StoredAs  string `json:"stored_as"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime_type"`
	CreatedAt string `json:"created_at"`
}

// Open loads (or creates) the store in dir
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, data: &Data{
		Seq:       map[string]int64{},
		Users:       []UserRecord{},
		Rooms:       []RoomRecord{},
		RoomMembers: []RoomMemberRecord{},
		Messages:    []MessageRecord{},
		Reactions:   []ReactionRecord{},
		Sessions:    []SessionRecord{},
		Files:       []FileRecord{},
	}}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) file(name string) string { return filepath.Join(s.dir, name) }

func (s *Store) load() error {
	data, err := os.ReadFile(s.file("khan.db.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return s.save() // create empty
		}
		return err
	}
	if err := json.Unmarshal(data, s.data); err != nil {
		return fmt.Errorf("corrupt store: %w", err)
	}
	if s.data.Seq == nil {
		s.data.Seq = map[string]int64{}
	}
	// Initialize nil slices to empty slices (JSON null -> nil in Go)
	if s.data.Users == nil { s.data.Users = []UserRecord{} }
	if s.data.Rooms == nil { s.data.Rooms = []RoomRecord{} }
	if s.data.RoomMembers == nil { s.data.RoomMembers = []RoomMemberRecord{} }
	if s.data.Messages == nil { s.data.Messages = []MessageRecord{} }
	if s.data.Reactions == nil { s.data.Reactions = []ReactionRecord{} }
	if s.data.Sessions == nil { s.data.Sessions = []SessionRecord{} }
	if s.data.Files == nil { s.data.Files = []FileRecord{} }
	return nil
}

// save atomically writes the store to disk
func (s *Store) save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.data, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	tmp := s.file("khan.db.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.file("khan.db.json"))
}

// SaveNow forces a flush to disk
func (s *Store) SaveNow() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *Store) saveLocked() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.file("khan.db.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.file("khan.db.json"))
}

// nextID returns the next auto-increment id for a table
func (s *Store) nextID(table string) int64 {
	s.data.Seq[table]++
	return s.data.Seq[table]
}

// Backup copies the data file
func (s *Store) Backup(destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("khan-backup-%s.json", time.Now().Format("20060102-150405"))
	dest := filepath.Join(destDir, name)
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(dest, data, 0o600); err != nil {
		return "", err
	}
	return dest, nil
}

// ---------- helpers ----------

func nowStr() string { return time.Now().Format(time.RFC3339) }

// Mu exposes the store mutex for repo transactions
func (s *Store) Mu() *sync.RWMutex { return &s.mu }

// Data exposes the in-memory dataset for repo transactions
func (s *Store) Data() *Data { return s.data }

// NextID returns the next auto-increment id for a table (caller holds lock)
func (s *Store) NextID(table string) int64 { return s.nextID(table) }

// SaveLocked persists while the caller holds the write lock
func (s *Store) SaveLocked() error { return s.saveLocked() }
