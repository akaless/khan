package service

import (
	"errors"
	"time"

	"khan/internal/models"
	"khan/internal/repository"
)

// MessageService handles message send/edit/delete/reactions with crypto at rest
type MessageService struct {
	msgs   *repository.MessageRepo
	rooms  *repository.RoomRepo
	users  *repository.UserRepo
	crypto *CryptoService
}

func NewMessageService(msgs *repository.MessageRepo, rooms *repository.RoomRepo, users *repository.UserRepo, crypto *CryptoService) *MessageService {
	return &MessageService{msgs: msgs, rooms: rooms, users: users, crypto: crypto}
}

// Send creates and stores an encrypted message
func (s *MessageService) Send(actor *models.User, roomID int64, text string, fileID *int64) (*models.MessageView, error) {
	isMember, err := s.rooms.IsMember(roomID, actor.ID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("شما عضو این اتاق نیستید")
	}
	if text == "" && fileID == nil {
		return nil, errors.New("پیام خالی است")
	}
	if len(text) > 8000 {
		return nil, errors.New("پیام خیلی طولانی است")
	}

	enc, err := s.crypto.Encrypt([]byte(text))
	if err != nil {
		return nil, err
	}

	msg := &models.Message{
		RoomID:    roomID,
		SenderID:  actor.ID,
		Body:      enc,
		FileID:    fileID,
		CreatedAt: time.Now(),
	}
	id, err := s.msgs.Create(msg)
	if err != nil {
		return nil, err
	}
	msg.ID = id

	view := &models.MessageView{
		ID:         msg.ID,
		RoomID:     msg.RoomID,
		SenderID:   msg.SenderID,
		SenderName: actor.DisplayName,
		Text:       text,
		FileID:     msg.FileID,
		CreatedAt:  msg.CreatedAt,
	}
	if view.SenderName == "" {
		view.SenderName = actor.Username
	}
	return view, nil
}

// ListMessages returns decrypted messages for a room (paginated)
func (s *MessageService) ListMessages(actor *models.User, roomID int64, before int64, limit int) ([]models.MessageView, error) {
	isMember, err := s.rooms.IsMember(roomID, actor.ID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("شما عضو این اتاق نیستید")
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	msgs, err := s.msgs.ListByRoom(roomID, before, limit)
	if err != nil {
		return nil, err
	}

	views := make([]models.MessageView, 0)
	for _, m := range msgs {
		text, decErr := s.crypto.Decrypt(m.Body)
		if decErr != nil {
			text = []byte("🔒")
		}
		view := models.MessageView{
			ID:        m.ID,
			RoomID:    m.RoomID,
			SenderID:  m.SenderID,
			Text:      string(text),
			FileID:    m.FileID,
			CreatedAt: m.CreatedAt,
			EditedAt:  m.EditedAt,
		}
		// sender name (deleted user → generic)
		if m.SenderID == -1 {
			view.SenderName = "کاربر حذف‌شده"
		} else if sender, err := s.users.GetByID(m.SenderID); err == nil && sender != nil {
			view.SenderName = sender.DisplayName
			if view.SenderName == "" {
				view.SenderName = sender.Username
			}
		} else {
			view.SenderName = "کاربر حذف‌شده"
		}
		// reactions
		if reacs, err := s.msgs.ListReactions(m.ID); err == nil {
			view.Reactions = reacs
		}
		views = append(views, view)
	}
	return views, nil
}

// EditMessage updates message text (owner only)
func (s *MessageService) EditMessage(actor *models.User, messageID int64, newText string) (*models.MessageView, error) {
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return nil, errors.New("پیام یافت نشد")
	}
	if msg.SenderID != actor.ID {
		return nil, errors.New("فقط فرستنده می‌تواند پیام را ویرایش کند")
	}
	if newText == "" {
		return nil, errors.New("متن پیام خالی است")
	}
	enc, err := s.crypto.Encrypt([]byte(newText))
	if err != nil {
		return nil, err
	}
	if err := s.msgs.UpdateBody(messageID, enc); err != nil {
		return nil, err
	}
	updated, _ := s.msgs.GetByID(messageID)
	view := &models.MessageView{
		ID:        messageID,
		RoomID:    msg.RoomID,
		SenderID:  actor.ID,
		SenderName: actor.DisplayName,
		Text:      newText,
		CreatedAt: msg.CreatedAt,
		EditedAt:  updated.EditedAt,
	}
	return view, nil
}

// DeleteMessage removes a message (owner or supervisor+ in room)
func (s *MessageService) DeleteMessage(actor *models.User, messageID int64) error {
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return errors.New("پیام یافت نشد")
	}

	if msg.SenderID != actor.ID {
		// supervisor+ can delete others' messages in rooms they manage
		if !actor.CanManageRooms() {
			return errors.New("فقط فرستنده می‌تواند پیام را حذف کند")
		}
		isMember, err := s.rooms.IsMember(msg.RoomID, actor.ID)
		if err != nil || !isMember {
			return errors.New("شما عضو این اتاق نیستید")
		}
	}
	return s.msgs.SoftDelete(messageID)
}

// AddReaction adds a reaction (member only)
func (s *MessageService) AddReaction(actor *models.User, messageID int64, emoji string) error {
	if emoji == "" || len([]rune(emoji)) > 8 {
		return errors.New("ایموجی نامعتبر است")
	}
	msg, err := s.msgs.GetByID(messageID)
	if err != nil || msg == nil {
		return errors.New("پیام یافت نشد")
	}
	isMember, err := s.rooms.IsMember(msg.RoomID, actor.ID)
	if err != nil || !isMember {
		return errors.New("شما عضو این اتاق نیستید")
	}
	return s.msgs.AddReaction(messageID, actor.ID, emoji)
}

// RemoveReaction removes a reaction
func (s *MessageService) RemoveReaction(actor *models.User, messageID int64, emoji string) error {
	return s.msgs.RemoveReaction(messageID, actor.ID, emoji)
}
