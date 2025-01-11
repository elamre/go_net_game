package game_system_packets

import (
	"github.com/elamre/go_net_game/net/packet_interface"
)

func Register() {
	packet_interface.RegisterPacket(EntityStatePacket{})
	packet_interface.RegisterPacket(PhysicsPacket{})
	packet_interface.RegisterPacket(InputStatePacket{})

}
