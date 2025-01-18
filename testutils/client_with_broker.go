package testutils

import (
	"fmt"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"log"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
)

const maxCh = 100
const start = 1800

// Semaphore представляет семафор синхронизации.
type Semaphore struct {
	pool  chan struct{}
	n     int
	ports map[int]struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{pool: make(chan struct{}, n), n: n, ports: make(map[int]struct{})}
}

func (s *Semaphore) Acquire() int {
	select {
	case s.pool <- struct{}{}:
		var newPort = start
		for {
			newPort = randRange(start, start+maxCh)
			_, ok := s.ports[newPort]
			if !ok {
				break
			}
		}

		s.ports[newPort] = struct{}{}
		return newPort
	}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.pool <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release(port int) {
	select {
	case <-s.pool:
		delete(s.ports, port)
	default:
		delete(s.ports, port)
	}
}

var sema = NewSemaphore(maxCh)

func GetClientWithBroker() (wb.ClientInterface, *mochi.Server, func()) {
	port := sema.Acquire()

	server := mochi.New(&mochi.Options{
		InlineClient: true,
	})

	_ = server.AddHook(new(auth.AllowHook), nil)
	tcp := listeners.NewTCP(listeners.Config{
		ID:      fmt.Sprintf("tcp-%d", port),
		Address: fmt.Sprintf(":%d", port),
	})
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	options := wb.Options{
		Broker:   fmt.Sprintf("localhost:%d", port),
		ClientId: "test-client",
		QoS:      1,
	}

	client, err := wb.NewClient(options)
	if err != nil {
		log.Fatal(err)
	}

	disconnect := func() {
		client.Disconnect(100)
		err := server.Close()
		if err != nil {
			log.Fatal(err)
		}

		sema.Release(port)
	}

	return client, server, disconnect
}

func AddOnHandler(server *mochi.Server) {
	callbackFn := func(cl *mochi.Client, sub packets.Subscription, pk packets.Packet) {
		server.Log.Info("inline client received message from subscription", "client", cl.ID, "subscriptionId", sub.Identifier, "topic", pk.TopicName, "payload", string(pk.Payload))
		topic := strings.Replace(pk.TopicName, "/on", "", 1)
		_ = server.Publish(topic, pk.Payload, true, 1)
	}
	server.Log.Info("inline client subscribing")
	subTopic := fmt.Sprintf(conventions.CONV_CONTROL_ON_VALUE_FMT, "+", "+")
	_ = server.Subscribe(subTopic, 1, callbackFn)
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
