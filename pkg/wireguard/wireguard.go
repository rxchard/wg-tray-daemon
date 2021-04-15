package wireguard

import (
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"os/exec"
	"strings"
)

type WgStat struct {
	wc *wgctrl.Client

	Enabled       bool
	Endpoint      string
	BytesSent     uint64
	BytesReceived uint64
}

func (s *WgStat) Update() error {
	if s.wc == nil {
		var err error

		s.wc, err = wgctrl.New()
		if err != nil {
			return err
		}
	}

	dev, err := s.wc.Device("wg-main0")
	if err != nil {
		return fmt.Errorf("device not found: %s", err)
	}

	if len(dev.Peers) == 0 {
		return fmt.Errorf("no peer found")
	}

	if len(dev.Peers) > 1 {
		return fmt.Errorf("only one peer is supported")
	}

	srv := dev.Peers[0]

	wgc := exec.Command("wg")
	buffer, err := wgc.Output()
	if err != nil {
		return err
	}

	s.Enabled = strings.Contains(string(buffer), srv.PublicKey.String())
	s.Endpoint = srv.Endpoint.IP.String()
	s.BytesSent = uint64(srv.TransmitBytes)
	s.BytesReceived = uint64(srv.ReceiveBytes)

	return nil
}
