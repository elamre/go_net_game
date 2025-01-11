package common

import (
	"github.com/elamre/go_net_game/net"
	"log"
	"sync"
)

type NetRoomSettings struct {
	HasPassword bool
}

type NetRoom struct {
	playerSync sync.Mutex
	RoomName   string
	Owner      int32
	Players    []*net.NetPlayer
	password   string
}

func (n *NetRoom) IsReady() bool {
	for _, p := range n.Players {
		if !p.Ready {
			log.Printf("%s not ready", p.Name)
			return false
		}
	}
	return true
}
