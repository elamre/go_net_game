package webrtc

import (
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net/webrtc/internal/webrtc_server"
	"time"
)

type WebrtcHostClient struct {
	connection *webrtc_server.Connection
	packetIdx  uint64
}

func NewWebrtcHostClient(connection *webrtc_server.Connection) net.ServerClient {
	return &WebrtcHostClient{connection: connection}
}

func (h *WebrtcHostClient) WritePacket(packet packet_interface.Packet) any {
	rawPack := NewRawPacket(packet)
	rawPack.PacketId = h.packetIdx
	rawPack.PacketTime = time.Now().UnixMilli()
	h.packetIdx++
	return h.connection.Send(rawPack.GetBytes())
}
func (h *WebrtcHostClient) GotPacket() bool {
	return true
}
func (h *WebrtcHostClient) ReadPacket() (packet packet_interface.Packet, err error) {
	bytes, data := h.connection.Read()
	if data {
		rawPack := RawPacketFrom(bytes)
		return rawPack.ContainingPacket, nil
	}
	return nil, nil
}
func (h *WebrtcHostClient) Close() any {
	return h.connection.CloseButDontFree
}
