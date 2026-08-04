package repository

import (
	"sort"
	"time"

	"khan/internal/database"
	"khan/internal/models"
)

// MessageRepo handles message persistence
type MessageRepo struct {
	store *database.Store
}

func NewMessageRepo(store *database.Store) *MessageRepo { return &MessageRepo{store: store} }

func recToMessage(r database.MessageRecord) *models.Message {
	m := &models.Message{
		ID:       r.ID,
		RoomID:   r.RoomID,
		SenderID: r.SenderID,
		Body:     []byte(r.Body),
		FileID:   r.FileID,
	}
	if t, err := time.Parse(time.RFC3339, r.CreatedAt); err == nil {
		m.CreatedAt = t
	}
	if r.EditedAt != "" {
		if t, err := time.Parse(time.RFC3339, r.EditedAt); err == nil {
			m.EditedAt = &t
		}
	}
	if r.DeletedAt != "" {
		if t, err := time.Parse(time.RFC3339, r.DeletedAt); err == nil {
			m.DeletedAt = &t
		}
	}
	return m
}

// Create inserts a message
func (r *MessageRepo) Create(m *models.Message) (int64, error) {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()

	m.ID = r.store.NextID("messages")
	rec := database.MessageRecord{
		ID:        m.ID,
		RoomID:    m.RoomID,
		SenderID:  m.SenderID,
		Body:      string(m.Body),
		FileID:    m.FileID,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
	}
	r.store.Data().Messages = append(r.store.Data().Messages, rec)
	return m.ID, r.store.SaveLocked()
}

// GetByID finds a message
func (r *MessageRepo) GetByID(id int64) (*models.Message, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	for _, rec := range r.store.Data().Messages {
		if rec.ID == id {
			return recToMessage(rec), nil
		}
	}
	return nil, nil
}

// ListByRoom returns messages for a room, paginated (before cursor), newest last
func (r *MessageRepo) ListByRoom(roomID int64, before int64, limit int) ([]models.Message, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()

	var msgs = make([]models.Message, 0)
	for _, rec := range r.store.Data().Messages {
		if rec.RoomID != roomID || rec.DeletedAt != "" {
			continue
		}
		if before > 0 && rec.ID >= before {
			continue
		}
		msgs = append(msgs, *recToMessage(rec))
	}
	sort.Slice(msgs, func(i, j int) bool { return msgs[i].ID < msgs[j].ID })
	if len(msgs) > limit {
		msgs = msgs[len(msgs)-limit:]
	}
	return msgs, nil
}

// UpdateBody updates message text (edit)
func (r *MessageRepo) UpdateBody(id int64, body []byte) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i, rec := range r.store.Data().Messages {
		if rec.ID == id {
			r.store.Data().Messages[i].Body = string(body)
			r.store.Data().Messages[i].EditedAt = time.Now().Format(time.RFC3339)
			return r.store.SaveLocked()
		}
	}
	return nil
}

// SoftDelete marks a message as deleted
func (r *MessageRepo) SoftDelete(id int64) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for i, rec := range r.store.Data().Messages {
		if rec.ID == id {
			r.store.Data().Messages[i].DeletedAt = time.Now().Format(time.RFC3339)
			return r.store.SaveLocked()
		}
	}
	return nil
}

// AddReaction inserts a reaction (idempotent per user+emoji)
func (r *MessageRepo) AddReaction(messageID, userID int64, emoji string) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	for _, re := range r.store.Data().Reactions {
		if re.MessageID == messageID && re.UserID == userID && re.Emoji == emoji {
			return nil
		}
	}
	rec := database.ReactionRecord{
		ID: r.store.NextID("reactions"), MessageID: messageID, UserID: userID, Emoji: emoji,
	}
	r.store.Data().Reactions = append(r.store.Data().Reactions, rec)
	return r.store.SaveLocked()
}

// RemoveReaction deletes a reaction
func (r *MessageRepo) RemoveReaction(messageID, userID int64, emoji string) error {
	r.store.Mu().Lock()
	defer r.store.Mu().Unlock()
	var kept []database.ReactionRecord
	for _, re := range r.store.Data().Reactions {
		if !(re.MessageID == messageID && re.UserID == userID && re.Emoji == emoji) {
			kept = append(kept, re)
		}
	}
	r.store.Data().Reactions = kept
	return r.store.SaveLocked()
}

// ListReactions returns reactions for a message
func (r *MessageRepo) ListReactions(messageID int64) ([]models.Reaction, error) {
	r.store.Mu().RLock()
	defer r.store.Mu().RUnlock()
	var reacs = make([]models.Reaction, 0)
	for _, rec := range r.store.Data().Reactions {
		if rec.MessageID == messageID {
			reacs = append(reacs, models.Reaction{
				ID: rec.ID, MessageID: rec.MessageID, UserID: rec.UserID, Emoji: rec.Emoji,
			})
		}
	}
	return reacs, nil
}
