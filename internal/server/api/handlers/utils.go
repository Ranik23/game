package handlers

import (
	"log"
	"github.com/gorilla/websocket"
)


func sendMessage(conn *websocket.Conn, action, message string) {
    msg := map[string]string{
        "Action":  action,
        "Message": message,
    }
    if err := conn.WriteJSON(msg); err != nil {
        log.Println("Failed to send error message:", err)
    }
}