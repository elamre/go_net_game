package common_system

import (
	"fmt"
	"github.com/elamre/go_helpers/pkg/slice_helpers"
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	"github.com/elamre/go_net_game/net_systems/common_system/common_packets"
	"log"
	"strings"
	"sync"
)

// These 2 functions are for ease

type UserManagement struct {
	playersAdjustMutex   sync.Mutex
	OnServerPlayerChange func(c *net.ServerPlayer, reason common_packets.ConnectionAction)
	players              []*net.ServerPlayer
	ClientToPlayer       map[net.ServerClient]*net.ServerPlayer
	NameToClient         map[string]*net.ServerPlayer
	IdToClient           map[int32]*net.ServerPlayer
	clientIdx            int32
}

func NewUserManagement(server net.Server) *UserManagement {
	s := &UserManagement{
		players:        make([]*net.ServerPlayer, 0),
		NameToClient:   make(map[string]*net.ServerPlayer),
		IdToClient:     make(map[int32]*net.ServerPlayer),
		ClientToPlayer: make(map[net.ServerClient]*net.ServerPlayer),
		clientIdx:      3,
	}
	server.AddConnectionCallback(s.OnConnection)
	server.AddDisconnectionCallback(s.OnDisconnection)
	return s
}

func (u *UserManagement) BroadcastFilter(pack packet_interface.Packet, filter func(player *net.ServerPlayer) bool) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	for _, p := range u.players {
		if filter(p) {
			p.Client.WritePacket(pack)
		}
	}
}

func (u *UserManagement) OnConnection(player *net.ServerPlayer) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()
	u.players = append(u.players, player)
	u.ClientToPlayer[player.Client] = player
	log.Printf("Connected: %p-[%+v], list: %+v, serverPlayer %p-[%+v]", player.Client, player.Client, u.ClientToPlayer, player, player)
}

func (u *UserManagement) OnDisconnectionReason(client net.ServerClient, reason string) {
	netP := u.ClientToPlayer[client].NetPlayer
	for _, p := range u.players {
		p.WritePacket(common_packets.ConnectionPacket{
			UserId:  netP.Id,
			Action:  common_packets.ConnectionDisconnectedAction,
			Message: fmt.Sprintf("%s %s", netP.Name, reason),
		})
	}
	u.OnDisconnection(u.ClientToPlayer[client])
}

func (u *UserManagement) OnDisconnection(player *net.ServerPlayer) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()

	serverClient := u.ClientToPlayer[player.Client]
	tempList := slice_helpers.RemoveFromList[*net.ServerPlayer](u.ClientToPlayer[player.Client], u.players)
	if tempList == nil {
		log.Println("client not found")
		return
	}
	u.players = tempList
	if serverClient.NetPlayer.HasRegistered {
		if u.OnServerPlayerChange != nil {
			u.OnServerPlayerChange(serverClient, common_packets.ConnectionDisconnectedAction)
		}
		serverClient.NetPlayer.HasRegistered = false
		delete(u.NameToClient, serverClient.NetPlayer.Name)
		delete(u.IdToClient, serverClient.NetPlayer.Id)
	} else {
		log.Println("we were never registered")
	}
	delete(u.ClientToPlayer, player.Client)
}

func (u *UserManagement) HandleConnectionPacket(c *net.ServerPlayer, d ServerRegulator, pack packet_interface.Packet) {
	u.playersAdjustMutex.Lock()
	defer u.playersAdjustMutex.Unlock()

	conPacket := pack.(common_packets.ConnectionPacket)
	if _, ok := u.ClientToPlayer[c.Client]; !ok {
		log.Printf("We don't exist yet2!")
	}

	log.Printf("connection packet: %+v ", conPacket)

	if u.ClientToPlayer[c.Client].NetPlayer.HasRegistered {
		c.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "Already registered",
		})
		log.Printf("Already registered!")
		return
	}
	if len(conPacket.Message) == 0 {
		c.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "No name given",
		})
		log.Printf("No name given!")
		return
	}
	name := strings.TrimSpace(strings.ToLower(conPacket.Message))
	if _, ok := u.NameToClient[name]; ok {
		c.WritePacket(common_packets.ConnectionPacket{
			UserId:  0,
			Action:  common_packets.ConnectionRefusedAction,
			Message: "Name already registered",
		})
		log.Printf("Name already exists!")
		return
	}
	u.ClientToPlayer[c.Client].NetPlayer.HasRegistered = true
	u.ClientToPlayer[c.Client].NetPlayer.Name = name
	u.ClientToPlayer[c.Client].NetPlayer.Id = u.clientIdx

	u.IdToClient[u.clientIdx] = u.ClientToPlayer[c.Client]
	u.NameToClient[name] = u.ClientToPlayer[c.Client]

	conPacket.Action = common_packets.ConnectionAcceptedAction
	conPacket.Message = name
	conPacket.UserId = u.clientIdx
	log.Printf("sending register packet: %+v", conPacket)
	c.WritePacket(conPacket)
	// TODO maybe broadcast this? we got a new user!
	u.clientIdx++
	if u.OnServerPlayerChange != nil {
		u.OnServerPlayerChange(u.ClientToPlayer[c.Client], conPacket.Action)
	}
}

func (s *UserManagement) RegisterCallbacks(r ServerRegulator) {
	r.RegisterPacketCallback(s.HandleConnectionPacket, common_packets.ConnectionPacket{})
}
func (s *UserManagement) Update(r ServerRegulator) {

}
