package lobby_system

import (
	"fmt"
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net_systems/common_system"
	"github.com/elamre/go_net_game/net_systems/common_system/common_packets"
	"github.com/elamre/go_net_game/net_systems/lobby_system/common"
	. "github.com/elamre/go_net_game/net_systems/lobby_system/lobby_system_packets"
	"log"
)

type LobbyCallback func(success bool, reason string, player *net.NetPlayer)

type LobbyClientSystem struct {
	OnRegisteredAction LobbyCallback
	OnRoomCreate       LobbyCallback
	OnRoomJoin         LobbyCallback
	OnRoomStart        LobbyCallback
	CurPlayer          *net.NetPlayer
	client             net.Client
	curRoom            *common.NetRoom
}

func NewLobbyClientSystem(client net.Client) *LobbyClientSystem {
	c := &LobbyClientSystem{client: client, CurPlayer: &net.NetPlayer{}}
	return c
}

func (c *LobbyClientSystem) SendRoomPacket(pack RoomPacket) any {
	pack.UserId = c.CurPlayer.Id
	return c.client.WritePacket(pack)

}

func (c *LobbyClientSystem) SendPacket(pack packet_interface.Packet) any {
	return c.client.WritePacket(pack)
}

func (c *LobbyClientSystem) RegisterPlayer(name string) {
	if !c.CurPlayer.HasRegistered {
		c.client.WritePacket(common_packets.ConnectionPacket{
			Action:  common_packets.ConnectionRegisterAction,
			Message: name,
		})
		c.CurPlayer.HasRegistered = true
	}
}

func (p *LobbyClientSystem) connectionCallback(c net.Client, d common_system.ClientRegulator, pack packet_interface.Packet) {
	t := pack.(common_packets.ConnectionPacket)
	log.Printf("Client received: %s", t.String())
	if t.Action == common_packets.ConnectionRefusedAction {
		p.CurPlayer.HasRegistered = false
		if p.OnRegisteredAction != nil {
			p.OnRegisteredAction(false, t.Message, p.CurPlayer)
		}
		log.Printf("Unable to register: %s", t.Message)
	} else if t.Action == common_packets.ConnectionAcceptedAction {
		p.CurPlayer.Id = t.UserId
		if p.OnRegisteredAction != nil {
			p.OnRegisteredAction(true, t.Message, p.CurPlayer)
		}
	} else {
		log.Printf("Player with name: %s registered (%d)", t.Message, t.UserId)
	}
}
func (lp *LobbyClientSystem) roomupdateCallback(c net.Client, d common_system.ClientRegulator, pack packet_interface.Packet) {
	/*	t := pack.(packets2.RoomUpdatePacket)
		for _, room := range t.Rooms {
			log.Printf("Room: %s", room.RoomName)
			players := ""
			for _, p := range room.Players {
				players += p.Name + ", "
			}
			log.Println(players)
		}*/
}

func (lp *LobbyClientSystem) roomPacketCallback(c net.Client, d common_system.ClientRegulator, pack packet_interface.Packet) {
	t := pack.(RoomPacket)
	if t.Action == RoomJoinSuccessAction || t.Action == RoomCreateSuccessAction {
		lp.curRoom = &common.NetRoom{
			RoomName: t.Name,
			Owner:    t.UserId,
			Players:  nil,
		}
	}
	log.Printf("Received: %+v", t)
}

func (p *LobbyClientSystem) RegisterCallbacks(r common_system.ClientRegulator) {
	r.RegisterPacketCallback(p.connectionCallback, common_packets.ConnectionPacket{})
	r.RegisterPacketCallback(p.roomupdateCallback, RoomUpdatePacket{})
	r.RegisterPacketCallback(p.roomPacketCallback, RoomPacket{})
}

func (c *LobbyClientSystem) Update() {
	if !c.client.IsConnected() {
		return
	}
	// We can draw here something
}

func (c *LobbyClientSystem) CreateRoom(roomName string, password string) any {
	if !c.client.IsConnected() {
		return fmt.Errorf("not connected")
	}
	c.SendRoomPacket(RoomPacket{
		UserId:   c.CurPlayer.Id,
		Action:   RoomCreateAction,
		Password: password,
		Name:     roomName,
	})
	return nil
}

func (c *LobbyClientSystem) StartRoom() any {
	if !c.client.IsConnected() {
		return fmt.Errorf("not connected")
	}
	c.SendRoomPacket(RoomPacket{
		UserId: c.CurPlayer.Id,
		Action: RoomStartAction,
		Name:   c.curRoom.RoomName,
	})
	return nil
}

func (c *LobbyClientSystem) SetReady() any {
	if !c.client.IsConnected() {
		return fmt.Errorf("not connected")
	}
	c.SendRoomPacket(RoomPacket{
		UserId: c.CurPlayer.Id,
		Action: RoomReadyAction,
		Name:   c.curRoom.RoomName,
	})
	return nil
}
