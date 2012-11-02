// This package tests both the server and the client.
// Assumes a runnin server on port 127.0.0.1:8912
package gotrigga_test

import(
	gt "github.com/opesun/gotrigga"
	"testing"
	"time"
)

type Waiter struct {
	interval	time.Duration
	limit		time.Duration
}

func (w *Waiter) Interval(i time.Duration) *Waiter {
	w.interval = i
	return w
}

func (w *Waiter) Limit(i time.Duration) *Waiter {
	w.limit = i
	return w
}

// Returns true if f returned true under given time.
func (w *Waiter) Wait(f func() bool) bool {
	t := time.Now()
	for {
		if f() == true {
			return true
		}
		if time.Since(t) > w.limit {
			return false
		}
		time.Sleep(w.interval)
	}
	panic("Should never happen.")
}

func TestSameconnSendRec(t *testing.T) {
	roomACounter := 0
	da := "hello"
	read := func(rec *gt.Room, c *int) {
		for {
			dat, err := rec.Read()
			if err != nil {
				t.Fatal(err)
				return
			}
			if string(dat) == da {
				*c++
			} else {
				t.Fatal("Wrong data:", string(dat))
			}
		}
	}
	conn, err := gt.Connect("127.0.0.1:8912")
	if err != nil {
		t.Fatal(err)
	}
	roomA := conn.Room("roomA")
	roomA.Subscribe()
	go read(roomA, &roomACounter)
	conn.Room("roomA").Publish(da)
	v := &Waiter{}
	ok := v.Interval(50*time.Millisecond).Limit(2*time.Second).Wait(func() bool {
		return roomACounter == 1
	})
	if !ok {
		t.Fatal(roomACounter)
	}
}

func TestMultiroom(t *testing.T) {
	roomACounter := 0
	roomBCounter := 0
	da := "hello"
	read := func(rec *gt.Room, c *int) {
		for {
			dat, err := rec.Read()
			if err != nil {
				t.Fatal(err)
				return
			}
			if string(dat) == da {
				*c++
			} else {
				t.Fatal("Wrong data:", string(dat))
			}
		}
	}
	conn, err := gt.Connect("127.0.0.1:8912")
	if err != nil {
		t.Fatal(err)
	}
	conn1, err := gt.Connect("127.0.0.1:8912")
	if err != nil {
		t.Fatal(err)
	}
	roomA := conn.Room("roomA")
	err = roomA.Subscribe()
	if err != nil {
		panic(err)
	}
	roomB := conn.Room("roomB")
	err = roomB.Subscribe()
	if err != nil {
		panic(err)
	}
	go read(roomA, &roomACounter)
	go read(roomB, &roomBCounter)
	sendn := 1000
	for i:=0;i<sendn;i++{
		conn1.Room("roomA").Publish(da)
	}
	for i:=0;i<sendn;i++{
		conn1.Room("roomB").Publish(da)
	}
	v := &Waiter{}
	ok := v.Interval(50*time.Millisecond).Limit(2*time.Second).Wait(func() bool {
		return roomACounter == sendn && roomBCounter == sendn
	})
	if !ok {
		t.Fatal(roomACounter, roomBCounter)
	}
}

