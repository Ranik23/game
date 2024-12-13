package handlers

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)


func sendMessage(conn *websocket.Conn, action, message string) {
    errMsg := map[string]string{
        "Action":  action,
        "Message": message,
    }
    msgBytes, err := json.Marshal(errMsg)
    if err != nil {
        log.Println("Failed to marshal error message:", err)
        return
    }

    if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
        log.Println("Failed to send error message:", err)
    }
}