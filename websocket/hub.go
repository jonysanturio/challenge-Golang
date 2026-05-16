package websocket

import (
        "fmt"
        "log"
        "sync"
        "net/http"
        
        "github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
   // Registered clients.
   clients map[*Client]bool

   // Inbound messages from the clients.
   broadcast chan []byte

   // Register requests from the clients.
   register chan *Client

   // Unregister requests from clients.
   unregister chan *Client

   // Mutex for thread safety
   mutex sync.RWMutex
 }

 // NewHub creates a new hub instance.
 func NewHub() *Hub {
   return &Hub{
           broadcast:  make(chan []byte),
           register:   make(chan *Client),
           unregister: make(chan *Client),
           clients:    make(map[*Client]bool),
   }
 }

 // Run starts the hub's event loop.
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

          case message := <-h.broadcast:
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

// Client is a middleman between the websocket connection and the hub.

type Client struct {
   hub *Hub

   // The websocket connection.
   conn *websocket.Conn

   // Buffered channel of outbound messages.
   send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
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
          // Handle incoming messages if needed
          // For now, we just broadcast to all clients
          c.hub.broadcast <- message
  }
}
// writePump pumps messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
  defer func() {
          c.conn.Close()
  }()
  for {
          select {
          case message, ok := <-c.send:
                  if !ok {
                          // The hub closed the channel.
                          c.conn.WriteMessage(websocket.CloseMessage)
                        
                           return
                   }

                   w, err := c.conn.NextWriter(websocket.TextMessage)
                   if err != nil {
                           return
                   }
                   w.Write(message)

                   // Add queued chat messages to the current webso
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

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Printf("WebSocket upgrade error: %v", err)
    return
  }
client := &Client{hub: hub, conn: conn, send: make(chan []byte,)}
client.hub.register <- client

// Allow collection of memory referenced by the caller by doing all work in
// new goroutines.
go client.writePump()
go client.readPump()
}

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    // Allow all connections for simplicity; adjust for prod

 return true
  },
}

// Notify functions for broadcasting events
func NotifyProductCreated(product interface{}) {
  // Implementation would serialize product and broadcast
  // This is a placeholder for actual implementation
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