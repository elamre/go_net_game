package lobby_system_packets

import "github.com/elamre/go_net_game/net/packet_interface"

func Register() {
	packet_interface.RegisterPacket(RoomPacket{})
	packet_interface.RegisterPacket(RoomUpdatePacket{})
}
