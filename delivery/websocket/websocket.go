package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"planning-poker/application/planningsvc"
	"planning-poker/domain/planning"
	"planning-poker/infra"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler struct {
	planningSvc *planningsvc.PlanningService
	logger      *zap.Logger
	sessions    map[string]map[*websocket.Conn]bool
	mu          sync.Mutex
}

func NewWebsocketHandler(planningSvc *planningsvc.PlanningService) *WebsocketHandler {
	handler := &WebsocketHandler{
		planningSvc: planningSvc,
		logger:      infra.GetLogger(),
		sessions:    make(map[string]map[*websocket.Conn]bool),
	}
	go handler.Stats() // Start the stats logging in a separate goroutine
	return handler
}

// Stats logs the current number of active sessions
func (h *WebsocketHandler) Stats() {
	for {
		h.mu.Lock()
		activeSessions := len(h.sessions)
		h.mu.Unlock()
		h.logger.Info("Active WebSocket sessions", zap.Int("count", activeSessions))
		// Sleep for a while before logging again
		time.Sleep(10 * time.Second)
	}
}

func (h *WebsocketHandler) register(planningId string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.sessions[planningId]; !ok {
		h.sessions[planningId] = make(map[*websocket.Conn]bool)
	}
	h.sessions[planningId][conn] = true
}

func (h *WebsocketHandler) unregister(planningId string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.sessions[planningId]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.sessions, planningId)
		}
	}
}

func (h *WebsocketHandler) broadcast(planningId string, eventType string, payload interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, ok := h.sessions[planningId]
	if !ok {
		return
	}

	event := struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}{
		Type:    eventType,
		Payload: payload,
	}

	msg, err := json.Marshal(event)
	if err != nil {
		h.logger.Error("failed to marshal broadcast event", zap.Error(err))
		return
	}

	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			h.logger.Error("failed to write message during broadcast", zap.Error(err))
		}
	}
}

func (h *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("failed to upgrade connection", zap.Error(err))
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			h.logger.Error("failed to close connection", zap.Error(err))
		}
	}(conn)

	var planningId string
	var playerId string

	defer func() {
		if planningId != "" {
			h.unregister(planningId, conn)
			if p, err := h.planningSvc.Leave(planningId, playerId); err == nil {
				h.broadcast(planningId, "player_left", p)
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("unexpected close error", zap.Error(err))
			}
			break
		}

		var event struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}

		if err := json.Unmarshal(msg, &event); err != nil {
			h.logger.Error("failed to unmarshal event", zap.Error(err))
			continue
		}

		var newPlanningId string
		var newPlayerId string
		var ok bool

		switch event.Type {
		case "create":
			newPlanningId, newPlayerId, ok = h.handleCreate(event.Payload)
			if ok {
				planningId = newPlanningId
				playerId = newPlayerId
				h.register(planningId, conn)
			}
		case "join":
			newPlanningId, newPlayerId, ok = h.handleJoin(event.Payload)
			if ok {
				if planningId != "" && planningId != newPlanningId {
					h.unregister(planningId, conn)
				}
				planningId = newPlanningId
				playerId = newPlayerId
				h.register(planningId, conn)
			}
		case "vote":
			h.handleVote(event.Payload)
		case "reveal":
			h.handleReveal(event.Payload)
		case "reset":
			h.handleReset(event.Payload)
		case "close":
			h.handleClose(event.Payload)
		default:
			h.logger.Warn("unknown event type", zap.String("type", event.Type))
		}

		if planningId != "" {
			p, err := h.planningSvc.GetById(planningId, "")
			if err != nil {
				h.logger.Error("failed to get planning for broadcast", zap.Error(err))
				continue
			}
			h.broadcast(planningId, event.Type, p)
		}
	}
}

func (h *WebsocketHandler) handleCreate(payload json.RawMessage) (string, string, bool) {
	var p planning.Planning
	if err := json.Unmarshal(payload, &p); err != nil {
		h.logger.Error("failed to unmarshal create payload", zap.Error(err))
		return "", "", false
	}

	if err := h.planningSvc.Create(&p); err != nil {
		h.logger.Error("failed to create planning", zap.Error(err))
		return "", "", false
	}

	return p.Id, p.Owner.Id, true
}

func (h *WebsocketHandler) handleJoin(payload json.RawMessage) (string, string, bool) {
	var req struct {
		PlanningId string          `json:"planningId"`
		Player     planning.Player `json:"player"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		h.logger.Error("failed to unmarshal join payload", zap.Error(err))
		return "", "", false
	}

	_, err := h.planningSvc.Join(req.PlanningId, &req.Player)
	if err != nil {
		h.logger.Error("failed to join planning", zap.Error(err))
		return "", "", false
	}

	return req.PlanningId, req.Player.Id, true
}

func (h *WebsocketHandler) handleVote(payload json.RawMessage) {
	var req struct {
		PlanningId string `json:"planningId"`
		PlayerId   string `json:"playerId"`
		Value      int    `json:"value"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		h.logger.Error("failed to unmarshal vote payload", zap.Error(err))
	}

	if err := h.planningSvc.Vote(req.PlanningId, req.PlayerId, req.Value); err != nil {
		h.logger.Error("failed to vote", zap.Error(err))
	}
}

func (h *WebsocketHandler) handleReveal(payload json.RawMessage) {
	var req struct {
		PlanningId string `json:"planningId"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		h.logger.Error("failed to unmarshal reveal payload", zap.Error(err))
	}

	if _, err := h.planningSvc.RevealVotes(req.PlanningId); err != nil {
		h.logger.Error("failed to reveal votes", zap.Error(err))
	}
}

func (h *WebsocketHandler) handleReset(payload json.RawMessage) {
	var req struct {
		PlanningId string `json:"planningId"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		h.logger.Error("failed to unmarshal reset payload", zap.Error(err))
	}

	if err := h.planningSvc.ResetVotes(req.PlanningId); err != nil {
		h.logger.Error("failed to reset votes", zap.Error(err))
	}
}

func (h *WebsocketHandler) handleClose(payload json.RawMessage) {
	var req struct {
		PlanningId string `json:"planningId"`
	}

	if err := json.Unmarshal(payload, &req); err != nil {
		h.logger.Error("failed to unmarshal close payload", zap.Error(err))
	}

	if err := h.planningSvc.Close(req.PlanningId); err != nil {
		h.logger.Error("failed to close planning", zap.Error(err))
	}
}
