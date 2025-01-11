package net

import (
	"fmt"
)

type NetPlayer struct {
	HasRegistered bool
	Ready         bool
	Name          string
	Id            int32
	RoomId        string
}

func (n NetPlayer) String() string {
	retString := fmt.Sprintf("NetPlayer[%s[%d]]", n.Name, n.Id)
	if !n.HasRegistered {
		retString += " [not registered]"
	}
	if n.Ready {
		retString += " [ready]"
	}
	if n.RoomId != "" {
		retString += fmt.Sprintf(" [room: %s]", n.RoomId)
	}
	return retString
}
