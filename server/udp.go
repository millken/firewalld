package server

import (
	"bytes"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"

	"github.com/millken/firewalld/worker"
)

const (
	maxQueueSize  = 1000000
	UDPPacketSize = 1500
)

type message struct {
	addr   net.Addr
	msg    []byte
	length int
}

type messageQueue chan message

func (mq messageQueue) enqueue(m message) {
	mq <- m
}

func (mq messageQueue) dequeue() {
	for m := range mq {
		handleMessage(m.addr, m.msg[0:m.length])
		bufferPool.Put(m.msg)
	}
}

var bufferPool sync.Pool
var mq messageQueue

func receive(c net.PacketConn) {
	defer c.Close()

	for {
		msg := bufferPool.Get().([]byte)
		nbytes, addr, err := c.ReadFrom(msg[0:])
		if err != nil {
			log.Printf("Error %s", err)
			continue
		}
		mq.enqueue(message{addr, msg, nbytes})
	}
}

func handleMessage(addr net.Addr, msg []byte) {
	cmdargs := bytes.Split(msg, []byte(" "))
	var args []string
	for _, arg := range cmdargs {
		a := strings.ToLower(strings.TrimSpace(string(arg)))
		if a == "" {
			continue
		}
		args = append(args, a)
	}
	if len(args) == 3 {
		worker.AddJob(args[0], args[1], args[2], "")
	}
	if len(args) == 4 {
		worker.AddJob(args[0], args[1], args[2], args[3])
	}
}

func (self *Server) serverUDP(addr string, stopC <-chan struct{}) {
	c, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Printf("listen udp : %s", err)
		return
	}
	bufferPool = sync.Pool{
		New: func() interface{} { return make([]byte, UDPPacketSize) },
	}
	mq = make(messageQueue, maxQueueSize)
	maxWorkers := runtime.NumCPU()
	for i := 0; i < maxWorkers; i++ {
		go mq.dequeue()
		go receive(c)
	}
}
