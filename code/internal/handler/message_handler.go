package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"khan/internal/service"
)

// MessageHandler serves /api/messages/*
type MessageHandler struct {
	msgs *service.MessageService
}

func NewMessageHandler(msgs *service.MessageService) *MessageHandler {
	return &MessageHandler{msgs: msgs}
}

// List returns messages for a room (paginated by ?before=X&limit=N)
func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	roomID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	before, _ := strconv.ParseInt(r.URL.Query().Get("before"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	msgs, err := h.msgs.ListMessages(u, roomID, before, limit)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, msgs)
}

func (h *MessageHandler) Edit(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	view, err := h.msgs.EditMessage(u, msgID, req.Text)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	if err := h.msgs.DeleteMessage(u, msgID); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type reactionReq struct {
	Emoji string `json:"emoji"`
}

func (h *MessageHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req reactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.msgs.AddReaction(u, msgID, req.Emoji); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *MessageHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	u := userFrom(r)
	msgID, err := pathInt64(r, "id")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "شناسه نامعتبر است")
		return
	}
	var req reactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "درخواست نامعتبر است")
		return
	}
	if err := h.msgs.RemoveReaction(u, msgID, req.Emoji); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
