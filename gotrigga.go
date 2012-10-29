package gotrigga


import(
	"net"
	"encoding/json"
	"github.com/opesun/trigga/binhelper"
)

type Connection struct {
	conn	net.Conn
}

type Room struct {
	name	string
	c		*Connection
}

func (c *Connection) Room(roomName string) *Room {
	return &Room{
		roomName,
		c,
	}
}

func (r *Room) Subscribe() error {
	cmd := map[string]interface{}{
		"r":	r.name,
		"c":	"s",
	}
	return r.send(cmd)
}

func (r *Room) Unsubscribe() error {
	cmd := map[string]interface{}{
		"r":	r.name,
		"c":	"u",
	}
	return r.send(cmd)
}

func (r *Room) Publish(msg string) error {
	cmd := map[string]interface{}{
		"r":	r.name,
		"c":	"p",
		"m":	msg,
	}
	return r.send(cmd)
}

func (r *Room) send(cmd map[string]interface{}) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return binhelper.WriteMsg(r.c.conn, b)
}

func (r *Room) Read() ([]byte, error) {
	return binhelper.ReadMsg(r.c.conn)
}

func Connect(addr string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Connection {
		conn,
	}, nil
}