package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"khan/internal/models"
	"khan/internal/repository"
	"khan/internal/service"
	"khan/internal/ws"
)

// WSHandler upgrades HTTP to WebSocket and wires events to services
type WSHandler struct {
	hub      *ws.Hub
	auth     *service.AuthService
	msgs     *service.MessageService
	rooms    *service.RoomService
	users    *service.UserService
	roomRep  *repository.RoomRepo
	userRep  *repository.UserRepo
	msgRep   *repository.MessageRepo
}

func NewWSHandler(hub *ws.Hub, auth *service.AuthService, msgs *service.MessageService,
	rooms *service.RoomService, users *service.UserService,
	roomRep *repository.RoomRepo, userRep *repository.UserRepo, msgRep *repository.MessageRepo) *WSHandler {
	return &WSHandler{hub: hub, auth: auth, msgs: msgs, rooms: rooms, users: users, roomRep: roomRep, userRep: userRep, msgRep: msgRep}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true }, // LAN app
}

// ServeHTTP handles the /ws endpoint
func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		token = extractToken(r)
	}
	if token == "" {
		writeErr(w, http.StatusUnauthorized, "وارد نشده‌اید")
		return
	}
	u, err := h.auth.ValidateToken(token)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	display := u.DisplayName
	if display == "" {
		display = u.Username
	}
	client := ws.NewClient(h.hub, u.ID, u.Username, display, conn)
	h.hub.Register(client)

	// notify others of presence (skip hidden super admin)
	if !u.Hidden {
		h.hub.BroadcastToAll(ws.NewEvent(ws.EvPresence, 0, ws.PresencePayload{
			UserID: u.ID, Online: true, Name: display,
		}))
	}

	// clean up on disconnect
	defer func() {
		h.hub.Unregister(client)
		if !u.Hidden {
			h.hub.BroadcastToAll(ws.NewEvent(ws.EvPresence, 0, ws.PresencePayload{
				UserID: u.ID, Online: false,
			}))
		}
	}()

	go client.WritePump()
	client.ReadPump(h.handleEvent(u))
}

// handleEvent dispatches client events to services
func (h *WSHandler) handleEvent(u *models.User) func(ws.ClientEvent) {
	return func(ev ws.ClientEvent) {
		switch ev.Type {
		case ws.CevSendMessage:
			view, err := h.msgs.Send(u, ev.RoomID, ev.Text, ev.FileID)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			members, _ := h.roomRep.ListMembers(ev.RoomID)
			ids := memberIDs(members)
			payload := ws.MessagePayload{
				ID: view.ID, SenderID: view.SenderID, SenderName: view.SenderName,
				Text: view.Text, FileID: view.FileID, CreatedAt: view.CreatedAt,
			}
			h.hub.BroadcastToRoom(ev.RoomID, ids, ws.NewEvent(ws.EvMessage, ev.RoomID, payload))

		case ws.CevEditMessage:
			view, err := h.msgs.EditMessage(u, ev.MessageID, ev.Text)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			members, _ := h.roomRep.ListMembers(view.RoomID)
			h.hub.BroadcastToRoom(view.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvMessageEdit, view.RoomID, map[string]interface{}{
					"id": view.ID, "text": view.Text, "edited_at": view.EditedAt,
				}))

		case ws.CevDeleteMsg:
			msg, err := h.getMsgRoom(ev.MessageID)
			if err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			if err := h.msgs.DeleteMessage(u, ev.MessageID); err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			members, _ := h.roomRep.ListMembers(msg.RoomID)
			h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvMessageDelete, msg.RoomID, ev.MessageID))

		case ws.CevAddReaction:
			if err := h.msgs.AddReaction(u, ev.MessageID, ev.Emoji); err != nil {
				h.sendErr(u.ID, err.Error())
				return
			}
			msg, _ := h.getMsgRoom(ev.MessageID)
			if msg != nil {
				members, _ := h.roomRep.ListMembers(msg.RoomID)
				h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvReaction, msg.RoomID, map[string]interface{}{
						"message_id": ev.MessageID, "emoji": ev.Emoji, "user_id": u.ID,
					}))
			}

		case ws.CevDelReaction:
			msg, _ := h.getMsgRoom(ev.MessageID)
			if msg != nil {
				members, _ := h.roomRep.ListMembers(msg.RoomID)
				h.hub.BroadcastToRoom(msg.RoomID, memberIDs(members),
					ws.NewEvent(ws.EvReaction, msg.RoomID, map[string]interface{}{
						"message_id": ev.MessageID, "emoji": ev.Emoji, "user_id": u.ID, "remove": true,
					}))
			}

		case ws.CevTyping:
			members, _ := h.roomRep.ListMembers(ev.RoomID)
			h.hub.BroadcastToRoom(ev.RoomID, memberIDs(members),
				ws.NewEvent(ws.EvTyping, ev.RoomID, ws.TypingPayload{UserID: u.ID, Name: u.DisplayName}))

		case ws.CevMarkRead:
			// presence only — no persistence needed for v1
		}
	}
}

// getMsgRoom finds the room of a message
func (h *WSHandler) getMsgRoom(msgID int64) (*models.Message, error) {
	msg, err := h.msgRep.GetByID(msgID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, http.ErrNoLocation
	}
	return msg, nil
}

// sendErr sends an error event to a user
func (h *WSHandler) sendErr(userID int64, msg string) {
	h.hub.SendToUser(userID, ws.NewEvent("error", 0, map[string]string{"error": msg}))
}

func memberIDs(users []models.User) []int64 {
	ids := make([]int64, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	return ids
}

var _ = json.Marshal
