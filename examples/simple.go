package main

import(
	gt "github.com/opesun/gotrigga"
	"fmt"
	"time"
)

var counter = 0

func read(rec *gt.Room) {
	for {
		dat, err := rec.Read()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		if string(dat) == "hello" {
			counter++
		} else {
			fmt.Println("rev wrong data:", string(dat))
		}
	}
}

func main() {
	recnum := 8
	for i:=0;i<recnum;i++{
		c, err := gt.Connect("127.0.0.1:8912")
		if err != nil {
			fmt.Println(err)
			continue
		}
		wh := c.Room("whatever")
		wh.Subscribe()
		go read(wh)
	}
	sender, err := gt.Connect("127.0.0.1:8912")
	if err != nil {
		panic(err)
	}
	sendnum := 5000
	t := time.Now()
	fmt.Println("sending messages.")
	for i:=0;i<sendnum;i++{
		sender.Room("whatever").Publish("hello")
	}
	last_counter := counter
	iterations := 0
	for {
		time.Sleep(50*time.Millisecond)
		fmt.Println("counter: ", counter)
		if counter == last_counter {
			break
		}
		last_counter = counter
		iterations++
	}
	fmt.Println(float64(sendnum)/time.Since(t).Seconds(), "messages per second with", recnum, "subscribers.")
}