package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rxchard/wg-tray-daemon/internal/packets"
	"github.com/rxchard/wg-tray-daemon/pkg/wireguard"
	"log"
	"net"
	"os"
	"os/exec"
)

const socketAddr = "/var/run/wg-tray-daemon.sock"

var wgc *wireguard.WgStat
var pmgr *packets.PacketHandlerMgr

// wgServer is basically a connection handler.
// It deals with packets etc.
func wgServer(wgx context.Context, con net.Conn) {
	log.Println("client open")

	serverx, cancel := context.WithCancel(wgx)
	defer cancel()

	go func() {
		<-serverx.Done()
		_ = con.Close()
		log.Println("client end")
	}()

	buffer := make([]byte, 4096)

	for {
		n, err := con.Read(buffer)
		if err != nil {
			return
		}

		if err = pmgr.Handle(string(buffer[0:n]), &con); err != nil {
			log.Printf("handle error: %s\n", err)
		}
	}
}

// wgDeleteSocket deletes the unix socket.
// It "fails hard" if an error occurs and exits the application
func wgDeleteSocket() {
	if err := os.RemoveAll(socketAddr); err != nil {
		log.Fatal(err)
	}
}

// Starts listening on the constant socket address.
// Starts a new "thread" for each connection
// Technically we'd just have to accept ONE connection at a time, however it currently "supports" multiple connections.
func wgListen(parent context.Context) error {
	listx, cancel := context.WithCancel(parent)
	defer cancel()

	listener, err := net.Listen("unix", socketAddr)
	if err != nil {
		return err
	}

	go func() {
		<-listx.Done()
		listener.Close()
		wgDeleteSocket()
	}()

	err = os.Chmod(socketAddr, 0766)
	if err != nil {
		return err
	}

	for {
		var con net.Conn
		con, err = listener.Accept()

		// did the context exit
		if listx.Err() != nil {
			log.Println("user exit")
			return nil
		}

		if err != nil {
			return err
		}

		go wgServer(listx, con)
	}
}

// Initializes the packet handler
// It directly defines handler functions (TODO: replace this)
func pckInitHandlers() {
	wgc = &wireguard.WgStat{}
	if err := wgc.Update(); err != nil {
		wgc.Enabled = false
	}

	pmgr = &packets.PacketHandlerMgr{
		Handlers: map[string]*packets.PacketHandler{},
	}

	pmgr.Add("toggle", func(c *net.Conn) error {
		if wgc.Enabled {
			err := exec.Command("wg-quick", "down", "wg-main0").Run()
			if err != nil {
				return err
			}

			wgc.Enabled = false
			return nil
		}

		err := exec.Command("wg-quick", "up", "wg-main0").Run()
		if err != nil {
			return err
		}

		wgc.Enabled = true
		return nil
	})

	pmgr.Add("status", func(c *net.Conn) error {
		if err := wgc.Update(); err != nil {
			wgc.Enabled = false
		}

		buffer, err := json.Marshal(wgc)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(*c, "status:", base64.URLEncoding.EncodeToString(buffer))
		if err != nil {
			return err
		}

		return nil
	})
}

// Base Executor
func Execute(parent context.Context) error {
	wgDeleteSocket()
	pckInitHandlers()

	if err := wgListen(parent); err != nil {
		return err
	}

	return nil
}
