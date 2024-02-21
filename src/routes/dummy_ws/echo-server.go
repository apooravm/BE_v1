package dummy_ws

import (
	"strconv"
	"net/http"
	"fmt"

	"github.com/apooravm/multi-serve/src/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func EchoDummyWS(c echo.Context) error {
	ConnUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := ConnUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("WS_ERR:", err.Error())
		return err
	}

	defer conn.Close()

	client_id := strconv.Itoa(Id_Gen.GenerateNewID())
	var newClient Client
	newClient.Id = client_id

	for {
		var clientData []byte
		_, clientData, err := conn.ReadMessage()
		if err != nil {
			clientDisconnect(conn, client_id)
			return nil
		}

		switch string(clientData) {
		case "subscribe":
			_, hasSubbed := SubscribedUsersMap.GetClient(client_id)
			if !hasSubbed {
				SubscribedUsersMap.AddClient(client_id, &newClient)
			}

		default:
			broadcastToSubscribers(string(clientData))
		}
	}
}

func broadcastToSubscribers(data string) {
	for _, client := range SubscribedUsersMap.Clients {
		if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
			fmt.Println("BROAD_ALL_ERR:", err.Error())
			// utils.LogData(err.Error(), utils.CHAT_DEBUG)
			utils.LogData(err.Error(), utils.DUMMY_WS_LOG_PATH)
			return
		}
	}
}

// Handle Client socket disconnection
// Graceful handling prevents error logs
func clientDisconnect(conn *websocket.Conn, clientID string) {
	if _, isSubbed := SubscribedUsersMap.GetClient(clientID); isSubbed {
		SubscribedUsersMap.DeleteClient(clientID)
	}
	if _, isFound := UserMap.GetClient(clientID); isFound {
		UserMap.DeleteClient(clientID)
	}

	conn.Close()
	// Handle disconnection or error here
	// // Delete client from the map
	// ChatClientsMap.DeleteClient(client_id)
	// message := client.Username + " left!"
	// LogData(message, CHAT_LOG)
	// BroadcastServerMessageAll(message)
}
