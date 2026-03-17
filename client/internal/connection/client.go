package connection

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"store-ops-client/internal/config"
	"store-ops-client/internal/resource"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypePing           MessageType = "ping"
	MessageTypePong           MessageType = "pong"
	MessageTypeRegister       MessageType = "register"
	MessageTypeRegisterAck    MessageType = "register_ack"
	MessageTypeResourceReport MessageType = "resource_report"
	MessageTypeCommand        MessageType = "command"
	MessageTypeCommandResult  MessageType = "command_result"
)

type Message struct {
	Type    MessageType              `json:"type"`
	Payload map[string]interface{}   `json:"payload"`
}

type Client struct {
	conn        *websocket.Conn
	url         string
	token       string
	deviceID    string
	fingerprint string
	send        chan []byte
	recv        chan Message
	done        chan struct{}
	mu          sync.Mutex
	connected   bool
}

var client *Client

func Init() error {
	cfg := config.Get()
	client = &Client{
		url:         cfg.ServerURL,
		token:       cfg.Token,
		deviceID:    cfg.DeviceID,
		fingerprint: cfg.DeviceFingerprint,
		send:        make(chan []byte, 256),
		recv:        make(chan Message, 256),
		done:        make(chan struct{}),
	}
	return nil
}

func GetClient() *Client {
	return client
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	// 构建WebSocket URL
	url := c.url + "?token=" + c.token + "&fingerprint=" + c.fingerprint

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	c.connected = true

	// 启动读写协程
	go c.readPump()
	go c.writePump()
	go c.heartbeat()
	go c.resourceReporter()

	// 发送注册消息
	c.register()

	return nil
}

func (c *Client) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return
	}

	close(c.done)
	c.conn.Close()
	c.connected = false
}

func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

func (c *Client) register() {
	cfg := config.Get()
	msg := Message{
		Type: MessageTypeRegister,
		Payload: map[string]interface{}{
			"device_id":          c.deviceID,
			"device_fingerprint": c.fingerprint,
			"store_name":         cfg.StoreName,
			"rustdesk_id":        cfg.RustDeskID,
		},
	}
	c.Send(msg)
}

func (c *Client) Send(msg Message) {
	data, _ := json.Marshal(msg)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connected {
		c.send <- data
	}
}

func (c *Client) readPump() {
	defer c.Disconnect()

	for {
		select {
		case <-c.done:
			return
		default:
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("读取消息错误: %v", err)
				return
			}

			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				continue
			}

			c.handleMessage(msg)
		}
	}
}

func (c *Client) writePump() {
	for {
		select {
		case <-c.done:
			return
		case data := <-c.send:
			c.mu.Lock()
			if c.connected {
				c.conn.WriteMessage(websocket.TextMessage, data)
			}
			c.mu.Unlock()
		}
	}
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.Send(Message{Type: MessageTypePing})
		}
	}
}

func (c *Client) resourceReporter() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			info, err := resource.Collect()
			if err != nil {
				continue
			}
			c.Send(Message{
				Type: MessageTypeResourceReport,
				Payload: map[string]interface{}{
					"cpu_percent":        info.CPUPercent,
					"memory_percent":     info.MemoryPercent,
					"memory_available_gb": info.MemoryAvailableGB,
					"disks":              info.Disks,
				},
			})
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case MessageTypePong:
		// 心跳响应
	case MessageTypeRegisterAck:
		log.Printf("注册响应: %+v", msg.Payload)
	case MessageTypeCommand:
		c.handleCommand(msg)
	}
}

func (c *Client) handleCommand(msg Message) {
	// TODO: 实现命令执行
	log.Printf("收到命令: %+v", msg.Payload)
}