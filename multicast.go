package multicast

import (
	"net"
)

type (
	Group struct {
		IPV4Address *net.UDPAddr
		IPV4Conn    *net.UDPConn
	}
	Packet struct {
		Source      *net.UDPAddr
		Destination *net.UDPAddr
		Size        int
		Body        []byte
	}
)

const (
	ipv4address     = "239.0.0.0:1013"
	maxDatagramSize = 8192
)

func New() (*Group, error) {
	a4, err := net.ResolveUDPAddr("udp4", ipv4address)
	if err != nil {
		return nil, err
	}

	c4, err := net.DialUDP("udp4", nil, a4)
	if err != nil {
		return nil, err
	}

	g := &Group{
		IPV4Address: a4,
		IPV4Conn:    c4,
	}

	return g, nil
}

func (g *Group) Listen() (<-chan *Packet, error) {
	conn, err := net.ListenMulticastUDP("udp4", nil, g.IPV4Address)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(maxDatagramSize)

	out := make(chan *Packet)

	go func() {
		defer close(out)
		for {
			b := make([]byte, maxDatagramSize)
			n, src, err := conn.ReadFromUDP(b)
			if err != nil {
				return
			}
			if src.String() == g.IPV4Conn.LocalAddr().String() {
				continue
			}
			out <- &Packet{
				Source:      src,
				Destination: g.IPV4Address,
				Size:        n,
				Body:        b,
			}
		}
	}()

	return out, nil
}

func (g *Group) Write(b []byte) (int, error) {
	return g.IPV4Conn.Write(b)
}
