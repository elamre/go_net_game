package webrtc

import (
	"fmt"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net/webrtc/internal/webrtc_client"
	"log"
	"time"
)

type WebrtcClient struct {
	client    *webrtc_client.Client
	port      int
	ip        string
	packetIdx uint64
}

func NewWebrtcClient(ip string, port int) *WebrtcClient {
	c := WebrtcClient{ip: ip, port: port, packetIdx: 1}
	c.client = webrtc_client.New(webrtc_client.Options{
		IPAddress:     fmt.Sprintf("%s:%d", ip, port),
		ICEServerURLs: []string{"stun:127.0.0.1:3478"},
	})
	return &c
}

func (w *WebrtcClient) Connect() any {
	w.client.Start()
	startTime := time.Now()
	for time.Since(startTime) < time.Second*5 || w.client.GetLastError() != nil {
		if w.client.HasConnectedOnce() {
			return nil
		}
		if err := w.client.GetLastError(); err != nil {
			return err
		}
	}
	return fmt.Errorf("timed out connecting to %s:%d", w.ip, w.port)
}

func (w *WebrtcClient) Disconnect() any {
	return w.client.GetLastError()
}

func (w *WebrtcClient) WritePacket(packet packet_interface.Packet) any {
	if !w.client.IsConnected() {
		return fmt.Errorf("not connected")
	}
	rawPacket := NewRawPacket(packet)
	rawPacket.PacketId = w.packetIdx
	w.packetIdx++
	rawPacket.PacketTime = time.Now().UnixMilli()
	if err := w.client.Send(rawPacket.GetBytes()); err != nil {
		return err
	}

	return w.client.GetLastError()
}
func (w *WebrtcClient) SetPacketReceivedCallback(PacketReceived func(packet packet_interface.Packet)) {
	log.Println("Set received callback")
}

func (w *WebrtcClient) GotPacket() bool {
	log.Println("Checking packet")
	return true
}

func (w *WebrtcClient) ReadPacket() (packet_interface.Packet, error) {
	if data, success := w.client.Read(); success {
		rawPacket := RawPacketFrom(data)
		return rawPacket.ContainingPacket, nil
	}
	return nil, nil
}

func (w *WebrtcClient) IsConnected() bool {
	return w.client.IsConnected()
}

func (w *WebrtcClient) Close() any {
	return nil
}
