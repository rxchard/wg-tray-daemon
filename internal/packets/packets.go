package packets

import (
	"fmt"
	"net"
)

type PacketHandler struct {
	Handler func(c *net.Conn) error
}

type PacketHandlerMgr struct {
	Handlers map[string]*PacketHandler
}

func (m PacketHandlerMgr) Add(name string, handler func(c *net.Conn) error) {
	m.Handlers[name] = &PacketHandler{Handler: handler}
}

func (m PacketHandlerMgr) Handle(name string, c *net.Conn) error {
	var handler = m.Handlers[name]

	if handler == nil {
		return fmt.Errorf("handler not found: %s", name)
	}

	// log.Println(name)

	if err := handler.Handler(c); err != nil {
		return fmt.Errorf("handler fail: %s %s", name, err)
	}

	return nil
}
