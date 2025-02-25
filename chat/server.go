package chat

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func NewServer() *server {
	log.Println("Server ready to accept conections")
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Run() {

	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_MESSAGE:
			s.msg(cmd.client, cmd.args)
		}
	}

}

func (s *server) NewClient(conn net.Conn) {
	log.Print("New client has joined", conn.RemoteAddr().String())

	c := &client{
		name:     "anon",
		conn:     conn,
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	c.name = args[1]
	c.msg(fmt.Sprintf("Alright,I will call you %s", args[1]))

}

func (s *server) msg(c *client, args []string) {

	if c.room == nil {
		c.err(fmt.Errorf("you need to join a room to send a message"))
		return
	}

	message := strings.Join(args[1:], " ")
	c.room.broadcast(c, fmt.Sprintf("%s:%s", c.name, message))
}

func (s *server) quit(c *client) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())
	s.quitRoom(c)
	c.msg("Hope to see you back. Bye!")
	c.conn.Close()
}

func (s *server) listRooms(c *client) {
	c.msg("Available rooms are:")
	for roomName := range s.rooms {
		c.msg(roomName)
	}
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	s.quitRoom(c)
	r.members[c.conn.RemoteAddr()] = c
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.name))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) quitRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}
