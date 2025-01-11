package ping_system

import (
	"github.com/elamre/go_net_game/net"
	"github.com/elamre/go_net_game/net/packet_interface"
	common_system2 "github.com/elamre/go_net_game/net_systems/common_system"
	"github.com/elamre/go_net_game/net_systems/ping_system/ping_packets"
	"time"
)

const (
	checkFrequency  = 10 * time.Second
	timeoutDuration = 30 * time.Second
)

type PingServerSystem struct {
	pastCheck      time.Time
	lastTimeStamps map[*net.ServerPlayer]time.Time
	timedOutList   []*net.ServerPlayer
}

func NewPingServerSystem() *PingServerSystem {
	return &PingServerSystem{lastTimeStamps: map[*net.ServerPlayer]time.Time{}, timedOutList: make([]*net.ServerPlayer, 0)}
}

func (p *PingServerSystem) pingCallback(c *net.ServerPlayer, d common_system2.ServerRegulator, pack packet_interface.Packet) {
	c.WritePacket(pack)
	p.lastTimeStamps[c] = pack.(ping_packets.PingPacket).CreationTime
}

func (p *PingServerSystem) RegisterCallbacks(r common_system2.ServerRegulator) {
	r.RegisterPacketCallback(p.pingCallback, ping_packets.PingPacket{})
}

func (p *PingServerSystem) Update(r common_system2.ServerRegulator) {
	curTime := time.Now()
	if curTime.Sub(p.pastCheck) > checkFrequency {
		for c, t := range p.lastTimeStamps {
			if curTime.Sub(t) > timeoutDuration {
				p.timedOutList = append(p.timedOutList, c)
			}
		}
		if len(p.timedOutList) != 0 {
			for _, s := range p.timedOutList {
				r.GetSubsystem("serverUsers").(*common_system2.UserManagement).OnDisconnectionReason(s.Client, "timed-out")
				delete(p.lastTimeStamps, s)
			}
			p.timedOutList = make([]*net.ServerPlayer, 0)
		}
		p.pastCheck = curTime
	}
}
