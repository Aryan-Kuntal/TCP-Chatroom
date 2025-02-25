package chat

import "net"

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(c *client, msg string) {
	for addr, m := range r.members {
		if addr == c.conn.RemoteAddr() {
			continue
		}
		m.msg(msg)
	}
}
