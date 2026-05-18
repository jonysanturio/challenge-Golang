package websocket

import (
        "fmt"
        "log"
        "sync"
        "net/http"
        
        "github.com/gorilla/websocket"
)

type Hub struct {
   clients map[*Client]bool
   Broadcast chan []byte
   register chan *Client
   unregister chan *Client
   mutex sync.RWMutex
 }

 func NewHub() *Hub {
   return &Hub{
           Broadcast:  make(chan []byte),
           register:   make(chan *Client),
           unregister: make(chan *Client),
           clients:    make(map[*Client]bool),
   }
 }

 func (h *Hub) Run() {
   for {
           select {
           case client := <-h.register:
                   h.mutex.Lock()
                   h.clients[client] = true
                   h.mutex.Unlock()
                   log.Printf("Client registered. Total clients: %d", len(h.clients))

           case client := <-h.unregister:
                   h.mutex.Lock()
                   if _, ok := h.clients[client]; ok {
                           delete(h.clients, client)
                           close(client.send)
                           log.Printf("Client unregistered. Total clients: %d", len(h.clients))
                  }
                  h.mutex.Unlock()

          case message := <-h.Broadcast:
                  h.mutex.RLock()
                  for client := range h.clients {
                          select {
                          case client.send <- message:
                          default:
                                  close(client.send)
                                  delete(h.clients, client)
                          }
                  }
                  h.mutex.RUnlock()
          }
  }
}


type Client struct {
   hub *Hub
   conn *websocket.Conn
   send chan []byte
}


func (c *Client) readPump() {
  defer func() {
          c.hub.unregister <- c
          c.conn.Close()
  }()
  for {
          _, message, err := c.conn.ReadMessage()
          if err != nil {
                  if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                          log.Printf("WebSocket error: %v", err)
                  }
                  break
          }
          c.hub.Broadcast <- message
  }
}

func (c *Client) writePump() {
  defer func() {
          c.conn.Close()
  }()
  for {
  select {
          case message, ok := <-c.send:
                  if !ok {
                          c.conn.WriteMessage(websocket.PingMessage, []byte{})
                           return
                  }

                   w, err := c.conn.NextWriter(websocket.TextMessage)
                   if err != nil {
                           return
                   }
                   w.Write(message)

                   n := len(c.send)
                   for i := 0; i < n; i++ {
                           w.Write(<-c.send)
                   }

                   if err := w.Close(); err != nil {
                           return
                   }
           }
   }
}

// ServeWs maneja las peticiones websocket
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Printf("WebSocket upgrade error: %v", err)
    return
  }
client := &Client{hub: hub, conn: conn, send: make(chan []byte,)}
client.hub.register <- client

go client.writePump()
go client.readPump()
}
var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {

 return true
  },
}

// Notificacion para eventos 
func NotifyProductCreated(product interface{}) {
  fmt.Printf("Product created: %v\n", product)
}

func NotifyProductUpdated(product interface{}) {
  fmt.Printf("Product updated: %v\n", product)
}

func NotifyProductDeleted(productID uint) {
  fmt.Printf("Product deleted: %d\n", productID)
}

func NotifyCategoryCreated(category interface{}) {
  fmt.Printf("Category created: %v\n", category)
}

func NotifyCategoryUpdated(category interface{}) {
  fmt.Printf("Category updated: %v\n", category)
}

func NotifyCategoryDeleted(categoryID uint) {
  fmt.Printf("Category deleted: %d\n", categoryID)
}