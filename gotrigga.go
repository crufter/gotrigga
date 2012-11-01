package gotrigga


import(
	"net"
	"encoding/json"
	"github.com/opesun/trigga/binhelper"
	"sync"
	"time"
)

type Connection struct {
	conn	net.Conn
	chans	map[string]map[int64]chan[]byte
	mut		sync.RWMutex
}

func (c *Connection) sendOnChans(msg map[string]interface{}) {
	c.mut.Lock()
	defer c.mut.Unlock()
	m, ok := c.chans[msg["r"].(string)]
	if !ok {
		return
	}
	for _, v := range m {
		v <- []byte(msg["m"].(string))
	}
}

func (c *Connection) read() {
	for {
		msg, err := binhelper.ReadMsg(c.conn)
		if err != nil {
			panic(err)
		}
		var v interface{}
		err = json.Unmarshal(msg, &v)
		if err != nil {
			panic(err)
		}
		c.sendOnChans(v.(map[string]interface{}))
	}
}

func (c *Connection) Close() {
	c.conn.Close()
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

func (r *Room) regChan(id int64, c chan[]byte) {
	r.c.mut.Lock()
	defer r.c.mut.Unlock()
	m, ok := r.c.chans[r.name]
	if !ok {
		m = map[int64]chan[]byte{}
	}
	m[id] = c
	r.c.chans[r.name] = m
}

func (r *Room) unregChan(id int64, c chan[]byte) {
	r.c.mut.Lock()
	defer r.c.mut.Unlock()
	delete(r.c.chans[r.name], id)
	if len(r.c.chans[r.name]) == 0 {
		delete(r.c.chans, r.name)
	}
}

func (r *Room) Read() ([]byte, error) {
	id := time.Now().UnixNano()
	ch := make(chan[]byte)
	r.regChan(id, ch)	// I bet it's costly...
	defer r.unregChan(id, ch)
	msg := <- ch
	return msg, nil
}

func Connect(addr string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	ret := &Connection {
		conn,
		map[string]map[int64]chan[]byte{},
		sync.RWMutex{},
	}
	go ret.read()
	return ret, nil
}