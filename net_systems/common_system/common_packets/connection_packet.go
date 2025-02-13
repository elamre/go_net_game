package common_packets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/elamre/go_helpers/pkg/misc"
)

const (
	_                             = iota
	ConnectionRegisterAction      = iota
	ConnectionAcceptedAction      = iota
	ConnectionRefusedAction       = iota
	ConnectionNewRegisteredAction = iota
	ConnectionDisconnectedAction  = iota
)

type ConnectionAction uint32

func (c ConnectionAction) String() string {
	switch c {
	case ConnectionRegisterAction:
		return "ConnectionRegisterAction"
	case ConnectionAcceptedAction:
		return "ConnectionAcceptedAction"
	case ConnectionRefusedAction:
		return "ConnectionRefusedAction"
	case ConnectionNewRegisteredAction:
		return "ConnectionNewRegisteredAction"
	case ConnectionDisconnectedAction:
		return "ConnectionDisconnectedAction"
	}
	return "unknown"
}

type ConnectionPacket struct {
	UserId  int32
	Action  ConnectionAction
	Message string
}

func (c ConnectionPacket) String() string {
	return fmt.Sprintf("Id: %d Message: %s Action: %s", c.UserId, c.Message, c.Action.String())
}

func NewRegisterPacket(name string) ConnectionPacket {
	return ConnectionPacket{Action: ConnectionRegisterAction, Message: name}
}

func (c ConnectionPacket) ConnectionSuccessful() bool {
	return c.Action == ConnectionAcceptedAction
}

func (c ConnectionPacket) ToWriter(w *bytes.Buffer) {
	misc.CheckError(binary.Write(w, binary.LittleEndian, c.UserId))
	misc.CheckError(binary.Write(w, binary.LittleEndian, c.Action))
	misc.CheckErrorRetVal(w.WriteString(c.Message))
}
func (c ConnectionPacket) FromReader(r *bytes.Reader) any {
	misc.CheckError(binary.Read(r, binary.LittleEndian, &c.UserId))
	misc.CheckError(binary.Read(r, binary.LittleEndian, &c.Action))
	sstring := make([]byte, r.Len())
	misc.CheckErrorRetVal(r.Read(sstring))
	c.Message = string(sstring)
	return c
}

func (c ConnectionPacket) AckRequired() bool {
	return true
}
