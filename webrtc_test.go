package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestWebrtc(t *testing.T) {
	const (
		DATA_CHANNEL_LABLE = "chat"
	)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to Upgrade connection: %s\n", err.Error())
		}
		defer conn.Close()

		peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
		if err != nil {
			t.Fatalf("Failed to create PeerConnection: %s\n", err.Error())
		}
		defer peerConnection.Close()

		dataChannel, err := peerConnection.CreateDataChannel(DATA_CHANNEL_LABLE, nil)
		if err != nil {
			t.Fatalf("Failed to create DataChannel: %s\n", err.Error())
		}

		dataChannel.OnOpen(func() {
			t.Logf("DataChannel is open.\n")
			go func() {
				for {
					time.Sleep(5 * time.Second)
					err := dataChannel.SendText(time.Now().String())
					if err != nil {
						t.Logf("Failed to send message: %s\n", err.Error())
						break
					} else {
						t.Logf(time.Now().String())
					}
				}
			}()
		})

		dataChannel.OnClose(func() {
			t.Logf("DataChannel is closed.")
		})

		dataChannel.OnError(func(err error) {
			t.Logf("DataChannel error: %s\n", err)
		})

		dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			t.Logf("Received message: %s\n", string(msg.Data))
		})

		peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate != nil {
				candidateJSON, _ := json.Marshal(candidate.ToJSON())
				conn.WriteMessage(websocket.TextMessage, candidateJSON)
			}
		})

		peerConnection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
			t.Logf("WebRTC connection state: %s", pcs.String())
		})

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				t.Logf("Read error: %s\n", err)
				break
			}

			var msg map[string]any
			if err := json.Unmarshal(message, &msg); err != nil {
				t.Logf("JSON error: %s\n", err)
				continue
			}

			if msg["type"] == "offer" {
				offer := webrtc.SessionDescription{
					Type: webrtc.SDPTypeOffer,
					SDP:  msg["sdp"].(string),
				}
				if err := peerConnection.SetRemoteDescription(offer); err != nil {
					t.Logf("SetRemoteDescription error: %s\n", err.Error())
					continue
				} else {
					t.Logf("SetRemoteDescription succeed")
				}

				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					t.Logf("CreateAnswer error: %s\n", err.Error())
					continue
				} else {
					t.Logf("CreateAnswer succeed")
				}

				if err := peerConnection.SetLocalDescription(answer); err != nil {
					t.Logf("SetLocalDescription error: %s\n", err.Error())
					continue
				} else {
					t.Logf("SetLocalDescription succeed")
				}

				answerJSON, _ := json.Marshal(answer)
				conn.WriteMessage(websocket.TextMessage, answerJSON)

			} else if msg["candidate"] != nil {
				candidate := webrtc.ICECandidateInit{
					Candidate: msg["candidate"].(string),
				}
				if err := peerConnection.AddICECandidate(candidate); err != nil {
					t.Logf("AddICECandidate error: %s\n", err.Error())
				} else {
					t.Logf("AddICECandidate succeed")
				}
			}
		}
	})

	t.Logf("Server started at :8080")
	t.Fatal(http.ListenAndServe(":8080", nil))
}
