package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client 客户端连接
type Client struct {
	ID              uuid.UUID
	Conn            *websocket.Conn
	StoreID         uuid.UUID
	DeviceID        uuid.UUID
	DeviceFingerprint string
	Send            chan []byte
}

// Hub WebSocket 连接管理器
type Hub struct {
	Clients    map[uuid.UUID]*Client // StoreID -> Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mu         sync.RWMutex
}

// NewHub 创建 Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// Run 运行 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.StoreID] = client
			h.mu.Unlock()
			log.Printf("客户端已注册: StoreID=%s", client.StoreID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.StoreID]; ok {
				delete(h.Clients, client.StoreID)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("客户端已断开: StoreID=%s", client.StoreID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.StoreID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetClient 获取客户端
func (h *Hub) GetClient(storeID uuid.UUID) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.Clients[storeID]
	return client, ok
}

// IsOnline 检查门店是否在线
func (h *Hub) IsOnline(storeID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.Clients[storeID]
	return ok
}

// SendToStore 发送消息给指定门店
func (h *Hub) SendToStore(storeID uuid.UUID, message *Message) error {
	h.mu.RLock()
	client, ok := h.Clients[storeID]
	h.mu.RUnlock()

	if !ok {
		return ErrStoreOffline
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client.Send <- data
	return nil
}

// GetOnlineCount 获取在线数量
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}

var ErrStoreOffline = &StoreOfflineError{}

type StoreOfflineError struct{}

func (e *StoreOfflineError) Error() string {
	return "门店离线"
}

// ReadPump 读取消息
func (c *Client) ReadPump(onMessage func(*Message)) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		onMessage(&msg)
	}
}

// WritePump 发送消息
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}