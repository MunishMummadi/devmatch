package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MunishMummadi/devmatch/server/internal/models"            // Use new module path
	"github.com/MunishMummadi/devmatch/server/internal/services/database" // Use new module path
	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ChatHandler handles WebSocket and potentially REST endpoints for chat.
type ChatHandler struct {
	dbService *database.DBService
	hub       *Hub
}

// --- WebSocket Handling ---

// Define WebSocket timing/size constants (adjust as needed)
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 4
)

// Configure the WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a single WebSocket connection.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *WebSocketMessage
	register   chan *Client
	unregister chan *Client
	dbService  *database.DBService
	mu         sync.Mutex
}

// WebSocketMessage defines the structure for messages sent/received via WebSocket
type WebSocketMessage struct {
	Type           string `json:"type"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id,omitempty"`
	Content        string `json:"content"`
}

// NewHub creates and runs a new Hub.
func NewHub(dbService *database.DBService) *Hub {
	hub := &Hub{
		broadcast:  make(chan *WebSocketMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		dbService:  dbService,
	}
	go hub.run()
	return hub
}

// run starts the hub's processing loop.
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s. Total clients: %d", client.userID, len(h.clients))
			// Optionally send connection confirmation or user list update
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send) // Close the client's send channel
				log.Printf("Client unregistered: %s. Total clients: %d", client.userID, len(h.clients))
				// Optionally send user list update
			}
			h.mu.Unlock()
		case msg := <-h.broadcast: // msg is *WebSocketMessage from readPump
			log.Printf("Hub received broadcast message type: %s from %s for conv %s", msg.Type, msg.SenderID, msg.ConversationID)

			// 1. Create the database message model & generate ID
			newMessage := models.Message{
				ID:             uuid.NewString(), // Generate unique ID here
				ConversationID: msg.ConversationID,
				SenderID:       msg.SenderID, // Use senderID from authenticated client
				Content:        msg.Content,
				// Timestamp will be set by SaveMessage if zero
			}

			// 2. Save the message to the database
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Add timeout
			err := h.dbService.SaveMessage(ctx, &newMessage)                        // Pass pointer
			cancel()                                                                // Release context resources
			if err != nil {
				log.Printf("Error saving message from %s to DB for conv %s: %v", newMessage.SenderID, newMessage.ConversationID, err)
				// Should we notify the sender? For now, just log and skip broadcast.
				continue
			}
			log.Printf("Message %s saved to DB for conv %s (Timestamp: %v)", newMessage.ID, newMessage.ConversationID, newMessage.Timestamp)

			// --- Start Broadcasting Logic ---

			// 3. Get participants for this conversation
			participantsCtx, participantsCancel := context.WithTimeout(context.Background(), 3*time.Second)
			participantIDs, err := h.dbService.GetConversationParticipants(participantsCtx, newMessage.ConversationID)
			participantsCancel()
			if err != nil {
				log.Printf("Error fetching participants for conv %s: %v. Cannot broadcast.", newMessage.ConversationID, err)
				continue
			}
			if len(participantIDs) == 0 {
				log.Printf("Warning: No participants found for conv %s. Message not broadcast.", newMessage.ConversationID)
				continue
			}
			participantsSet := make(map[string]struct{}, len(participantIDs))
			for _, id := range participantIDs {
				participantsSet[id] = struct{}{}
			}

			// 4. Fetch sender's name (needed for FrontendMessage 'From' field)
			senderName := "Unknown User" // Default
			senderProfileCtx, senderProfileCancel := context.WithTimeout(context.Background(), 3*time.Second)
			senderProfile, err := h.dbService.GetUserProfileByID(senderProfileCtx, newMessage.SenderID)
			senderProfileCancel()
			if err == nil && senderProfile != nil && senderProfile.Username != nil {
				senderName = *senderProfile.Username
			} else if err != nil {
				log.Printf("Could not fetch sender profile for ID %s: %v", newMessage.SenderID, err)
				// Continue with default senderName
			}

			// 5. Iterate through *connected* clients and send if they are participants
			h.mu.Lock() // Lock before iterating clients map
			clientCount := len(h.clients)
			broadcastCount := 0
			log.Printf("Checking %d connected clients for participation in conv %s", clientCount, newMessage.ConversationID)

			for client := range h.clients {
				// Check if this client is a participant in the conversation
				if _, isParticipant := participantsSet[client.userID]; !isParticipant {
					continue // Skip client if not a participant
				}

				// Determine message type and sender name for this specific client
				messageType := models.GetMessageType(newMessage.SenderID, client.userID)
				fromName := senderName
				if messageType == "outgoing" {
					fromName = "You"
				}

				// Format the message specifically for this recipient
				frontendMsg := models.FrontendMessage{
					From:      fromName,
					Text:      newMessage.Content,
					Timestamp: models.FormatRelativeTime(newMessage.Timestamp),
					Type:      messageType,
				}

				// Marshal the tailored FrontendMessage to JSON
				messageBytes, err := json.Marshal(frontendMsg)
				if err != nil {
					log.Printf("Error marshalling frontend message for client %s, conv %s: %v", client.userID, newMessage.ConversationID, err)
					continue // Skip this client if marshalling fails
				}

				// Send the JSON message to the participant client
				select {
				case client.send <- messageBytes:
					broadcastCount++
					log.Printf("Sent message %s to participant client %s in conv %s", newMessage.ID, client.userID, newMessage.ConversationID)
				default:
					// If client's send buffer is full, assume client is lagging/dead.
					log.Printf("Client %s send buffer full for conv %s. Closing and unregistering.", client.userID, newMessage.ConversationID)
					close(client.send)
					delete(h.clients, client) // Delete while iterating requires care, but okay with range over map
				}
			}
			h.mu.Unlock() // Unlock after iterating clients map
			log.Printf("Broadcast attempt finished for message %s in conv %s. Sent to %d participants.", newMessage.ID, newMessage.ConversationID, broadcastCount)
		}
	}
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(db *database.DBService, hub *Hub) *ChatHandler {
	return &ChatHandler{
		dbService: db,
		hub:       hub, // Store the hub
	}
}

// --- WebSocket Handler ---

// HandleWebSocket upgrades the connection and starts client processing.
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	// 1. Get UserID from Clerk claims BEFORE upgrading
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		log.Println("WebSocket upgrade failed: User ID not found in context")
		// Cannot set headers after connection hijack attempt, just return
		c.AbortWithStatus(http.StatusUnauthorized) // Use Abort to prevent further writes
		return
	}
	log.Printf("Attempting WebSocket upgrade for user: %s", claims.Subject)

	// 2. Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed for user %s: %v", claims.Subject, err)
		// Upgrader handles sending the error response, so just return
		return
	}
	log.Printf("WebSocket upgrade successful for user: %s", claims.Subject)

	// 3. Create a new client for this connection
	client := &Client{
		hub:    h.hub,                  // Reference to the hub
		conn:   conn,                   // The WebSocket connection
		send:   make(chan []byte, 256), // Buffered channel for outgoing messages
		userID: claims.Subject,         // Authenticated user ID
	}

	// 4. Register the client with the hub
	client.hub.register <- client
	log.Printf("Client %s registered with hub", claims.Subject)

	// 5. Start goroutines for this client's read and write operations.
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump() // Handles sending messages from client.send to the WebSocket
	go client.readPump()  // Handles reading messages from WebSocket and sending to hub.broadcast

	// The handler function returns now, but the goroutines keep the connection alive.
	log.Printf("WebSocket handler finished setup for user %s", claims.Subject)
}

// readPump pumps messages from the WebSocket connection to the hub.
// The application runs readPump in a per-connection goroutine. It ensures
// that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	// Cleanup: Unregister client and close connection when readPump exits.
	defer func() {
		log.Printf("readPump closing for client %s", c.userID)
		c.hub.unregister <- c
		c.conn.Close()
		log.Printf("Client %s unregistered and connection closed", c.userID)
	}()

	// Set read limits and deadlines for security and resource management.
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Error setting read deadline for client %s: %v", c.userID, err)
		return // Exit if we can't set the deadline
	}
	c.conn.SetPongHandler(func(string) error {
		log.Printf("Pong received from client %s", c.userID)
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("Error setting read deadline after pong for client %s: %v", c.userID, err)
			// We might want to close the connection here if the deadline can't be extended
		}
		return nil
	})

	// Main read loop
	for {
		messageType, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			// Log specific close errors vs general errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket unexpected close error for client %s: %v", c.userID, err)
			} else {
				// Normal closure or expected errors (like connection closed by peer)
				log.Printf("WebSocket read error/closed for client %s: %v", c.userID, err)
			}
			break // Exit loop on any read error or closure
		}

		// We only process text messages containing JSON
		if messageType != websocket.TextMessage {
			log.Printf("Received non-text message type %d from client %s. Ignoring.", messageType, c.userID)
			continue
		}

		log.Printf("Received raw message from client %s: %s", c.userID, string(messageBytes))

		// Unmarshal the JSON message into our WebSocketMessage struct
		var wsMsg WebSocketMessage
		if err := json.Unmarshal(messageBytes, &wsMsg); err != nil {
			log.Printf("Error unmarshalling message from client %s: %v. Message: %s", c.userID, err, string(messageBytes))
			// Optionally send an error message back to the client? For now, just ignore malformed messages.
			continue
		}

		// *** Crucial: Add the authenticated SenderID ***
		wsMsg.SenderID = c.userID
		// Note: Message ID should be generated within the Hub or just before saving
		// as the WebSocketMessage struct doesn't carry the final DB message model.

		// We only handle 'newMessage' type from clients for now
		if wsMsg.Type != "newMessage" {
			log.Printf("Received unhandled message type '%s' from client %s", wsMsg.Type, c.userID)
			continue
		}
		if wsMsg.Content == "" || wsMsg.ConversationID == "" {
			log.Printf("Received incomplete 'newMessage' from client %s (missing content or conversation_id)", c.userID)
			continue
		}

		// Pass the validated message (as struct) to the hub's broadcast channel
		// The hub will handle saving it and distributing to relevant clients.
		// Use a select with a timeout to prevent blocking indefinitely if the hub is stuck
		select {
		case c.hub.broadcast <- &wsMsg: // Pass pointer to the struct
			log.Printf("Message from client %s sent to hub broadcast channel", c.userID)
		case <-time.After(2 * time.Second): // Timeout if hub is blocked for too long
			log.Printf("Hub broadcast channel full/blocked receiving from client %s. Message dropped.", c.userID)
			// Optionally notify the client that the message wasn't processed immediately
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) // Ticker for sending pings
	defer func() {
		ticker.Stop()  // Stop the ticker
		c.conn.Close() // Close the WebSocket connection
		log.Printf("writePump closing for client %s", c.userID)
		// Unregistration should be handled by readPump's defer when the connection breaks
	}()

	for {
		select {
		case messageBytes, ok := <-c.send:
			// Set write deadline before writing
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error setting write deadline for client %s: %v", c.userID, err)
				return // Exit pump
			}
			if !ok {
				// The hub closed the channel (likely because the client was unregistered).
				log.Printf("Hub closed send channel for client %s. Sending close message.", c.userID)
				// Attempt to send a WebSocket close message
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return // Exit pump
			}

			// Get a writer for the next message
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting next writer for client %s: %v", c.userID, err)
				return // Exit pump
			}

			// Write the message bytes
			if _, err := w.Write(messageBytes); err != nil {
				log.Printf("Error writing message bytes for client %s: %v", c.userID, err)
				// Attempt to close the writer even on error
				if closeErr := w.Close(); closeErr != nil {
					log.Printf("Error closing writer after write error for client %s: %v", c.userID, closeErr)
				}
				return // Exit pump
			}
			log.Printf("Message sent to client %s: %s", c.userID, string(messageBytes))

			// Optimize: Add queued chat messages to the current websocket message.
			// This reduces network overhead by batching messages that arrive close together.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	// Write a newline separator (optional, depends on client handling)
			// 	if _, err := w.Write([]byte{'\n'}); err != nil {
			// 		log.Printf("Error writing newline separator for client %s: %v", c.userID, err)
			// 		return
			// 	}
			// 	queuedMsg := <-c.send
			// 	if _, err := w.Write(queuedMsg); err != nil {
			// 		log.Printf("Error writing queued message for client %s: %v", c.userID, err)
			// 		return
			// 	}
			// 	log.Printf("Queued message sent to client %s: %s", c.userID, string(queuedMsg))
			// }

			// Close the writer to flush the message to the connection.
			if err := w.Close(); err != nil {
				log.Printf("Error closing writer for client %s: %v", c.userID, err)
				return // Exit pump
			}

		case <-ticker.C:
			// Send a ping message periodically to keep the connection alive and check responsiveness.
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error setting write deadline for ping for client %s: %v", c.userID, err)
				return // Exit pump
			}
			log.Printf("Sending ping to client %s", c.userID)
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to client %s: %v", c.userID, err)
				return // Assume connection is broken, exit pump
			}
		}
	}
}

// --- REST API Handlers for Chat (Example) ---

// GetConversations fetches the combined chat data (conversation summaries and initial messages)
// GET /chat/conversations (Clerk User ID taken from middleware)
func (h *ChatHandler) GetConversations(c *gin.Context) {
	// Get user ID from Clerk claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	userID := claims.Subject

	initialMessagesLimit := 20

	chatData, err := h.dbService.GetChatDataForUser(c.Request.Context(), userID, initialMessagesLimit)
	if err != nil {
		log.Printf("Error fetching chat data for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat data"})
		return
	}

	c.JSON(http.StatusOK, chatData)
}

// GetMessages fetches messages for a specific conversation.
// GET /chat/messages/:conversationId
func (h *ChatHandler) GetMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")

	limit := 50
	offset := 0

	messages, err := h.dbService.GetMessagesForConversation(c.Request.Context(), conversationID, limit, offset)
	if err != nil {
		log.Printf("Error fetching messages for conversation %s: %v", conversationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	if messages == nil {
		messages = []models.Message{}
	}
	c.JSON(http.StatusOK, messages)
}

// Define a struct to bind the SendMessage request body
type SendMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Content        string `json:"content" binding:"required"`
}

// SendMessage handles sending a new message via REST and broadcasts it.
// POST /chat/message
func (h *ChatHandler) SendMessage(c *gin.Context) {
	// Get sender ID from Clerk claims
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session claims found"})
		return
	}
	senderID := claims.Subject

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding SendMessage request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 1. Prepare the message model
	newMessage := models.Message{
		// ID will be generated by SaveMessage or DB trigger
		ConversationID: req.ConversationID,
		SenderID:       senderID,
		Content:        req.Content,
		// Timestamp will be set by SaveMessage
	}

	// 2. Save the message to the database
	err := h.dbService.SaveMessage(c.Request.Context(), &newMessage)
	if err != nil {
		log.Printf("Error saving message from user %s to conversation %s: %v", senderID, req.ConversationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}
	log.Printf("REST API saved message %s to DB for conv %s", newMessage.ID, newMessage.ConversationID)

	// 3. Broadcast the message via Hub (after successful save)
	// Create the WebSocketMessage structure for the hub
	broadcastMsg := &WebSocketMessage{
		Type:           "newMessage", // Indicate this is a new message
		ConversationID: newMessage.ConversationID,
		SenderID:       newMessage.SenderID,
		Content:        newMessage.Content,
		// The hub's run method will handle fetching the timestamp and formatting
	}

	// Send to the hub's broadcast channel (non-blocking send)
	// Use a select with a default to prevent blocking if the hub is busy
	select {
	case h.hub.broadcast <- broadcastMsg:
		log.Printf("Sent message %s from REST API to Hub broadcast channel for conv %s", newMessage.ID, newMessage.ConversationID)
	default:
		// This should ideally not happen if the hub is running correctly.
		// Log a warning if the broadcast channel is full.
		log.Printf("Warning: Hub broadcast channel full. Message %s from REST API for conv %s might not be broadcast immediately.", newMessage.ID, newMessage.ConversationID)
	}

	// 4. Return the created message object in the HTTP response
	c.JSON(http.StatusCreated, newMessage)
}
