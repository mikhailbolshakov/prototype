package listener

import (
	"encoding/json"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/queue/stan"
	"log"
	"math/rand"
	"testing"
	"time"
)

type TestMsg struct {
	V string
}

func Test(t *testing.T) {

	publisher := &stan.Stan{}
	err := publisher.Open(fmt.Sprintf("test_publisher_%d", rand.Intn(99999)))
	if err != nil {
		t.Fatal(err)
	}

	subscriber := &stan.Stan{}
	err = subscriber.Open(fmt.Sprintf("test_subscriber_%d", rand.Intn(99999)))
	if err != nil {
		t.Fatal(err)
	}

	h1 := func(v []byte) error {
		d := &TestMsg{}
		_ = json.Unmarshal(v, d)
		log.Printf("recieved h1: %v\n", d)
		return nil
	}

	h2 := func(v []byte) error {
		d := &TestMsg{}
		_ = json.Unmarshal(v, d)
		log.Printf("recieved h2: %v\n", d)
		return nil
	}
	//
	//c1 := make(chan []byte)
	//c2 := make(chan []byte)

	//f := func(cc chan []byte, h func(v []byte)) {
	//	select {
	//		case m := <-cc: h(m)
	//	}
	//}
	//
	//go f(c1, h1)
	//go f(c2, h2)
	//
	//_ = subscriber.Subscribe("test.topic", c1)
	//_ = subscriber.Subscribe("test.topic", c2)

	listener := NewQueueListener(subscriber)
	listener.Add("test.topic", h1, h2)
	listener.ListenAsync()
	time.Sleep(time.Second * 1)

	m := &TestMsg{"value1"}
	b, _ := json.Marshal(m)
	if err := publisher.Publish("test.topic", b); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 1000)

}
