package common_packets

import (
	"github.com/elamre/go_net_game/net/packet_interface"
)

func Register() {
	packet_interface.RegisterPacket(ChatPacket{})
	packet_interface.RegisterPacket(ConnectionPacket{})
}
