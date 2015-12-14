package consul

import (
	"github.com/hashicorp/serf/serf"
	"net"
	"strconv"
	"strings"
)

type NodeNameAddress struct {
	nodeName string
	port     uint16
}

const nodeNameNetwork = "nodename-net"

func (addr NodeNameAddress) Network() string { return nodeNameNetwork }
func (addr NodeNameAddress) String() string  { return addr.nodeName + ":" + strconv.Itoa(int(addr.port)) }

type NodeNameIPResolver struct {
	serf **serf.Serf
}

func (t *NodeNameIPResolver) ResolveAddr(addr net.Addr) net.Addr {
	if addr.Network() == nodeNameNetwork {
		nodeNameAddress := addr.(NodeNameAddress)
		for _, member := range (*t.serf).Members() {
			if member.Name == nodeNameAddress.nodeName {
				addr = &net.TCPAddr{IP: member.Addr, Port: int(nodeNameAddress.port)}
			}
		}
	}
	return addr
}

func (t *NodeNameIPResolver) ResolveString(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 && t.serf != nil {
		for _, member := range (*t.serf).Members() {
			if member.Name == parts[0] {
				addr = net.JoinHostPort(member.Addr.String(), parts[1])
			}
		}
	}
	return addr
}
