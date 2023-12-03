package echovr

import (
	"log"
	"net/http"

	"github.com/unusualnorm/echovr_lib/messages"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print(err)
			return
		}

		sendMessage := func(message Message) error {
			log.Printf("send: %v", message)

			symbol := message.Symbol()
			b, err := message.Serialize()
			if err != nil {
				return err
			}

			packet := &Packet{Header: PACKET_HEADER, Symbol: symbol, Data: b}
			p, err := packet.Serialize()
			if err != nil {
				return err
			}

			return conn.WriteMessage(websocket.BinaryMessage, p)
		}

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			packet := &Packet{}
			err = packet.Deserialize(p)
			if err != nil {
				log.Println("deserialize:", err)
				continue
			}

			message := Message(nil)
			switch packet.Symbol {
			case messages.STcpConnectionUnrequireEventSymbol:
				message = &messages.STcpConnectionUnrequireEvent{}
			case messages.SNSConfigRequestv2Symbol:
				message = &messages.SNSConfigRequestv2{}
			case messages.SNSConfigSuccessv2Symbol:
				message = &messages.SNSConfigSuccessv2{}
			case messages.SNSConfigFailurev2Symbol:
				message = &messages.SNSConfigFailurev2{}
			}

			if message != nil {
				err = message.Deserialize(packet.Data)
				if err != nil {
					log.Println("deserialize:", err)
					continue
				}
				log.Printf("recv: %v", message)
			} else {
				log.Printf("recv: %v", packet)
			}

			sendMessage(&messages.STcpConnectionUnrequireEvent{Unused: 0x00})
		}
	}))
}
