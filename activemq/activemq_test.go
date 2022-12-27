package main

import (
	"log"
	"os"
	"testing"

	"github.com/go-stomp/stomp/v3"

	"github.com/ory/dockertest/v3"
)

var activeMqResource *dockertest.Resource

func TestMain(m *testing.M) {

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	activeMqResource, err = pool.Run("rmohr/activemq", "5.15.9", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		mqPort := activeMqResource.GetPort("61613/tcp")
		conn, err := stomp.Dial("tcp", "localhost:"+mqPort)
		if err != nil {
			return err
		}
		return conn.Disconnect()
	}); err != nil {
		log.Fatalf("Could not connect to mq: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(activeMqResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestActiveMQ(t *testing.T) {
	mqPort := activeMqResource.GetPort("61613/tcp")
	conn, err := stomp.Dial("tcp", "localhost:"+mqPort)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe("/queue/test-2", stomp.AckClient)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	err = conn.Send("/queue/test-2", "application/xml",
		[]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	msg, err := sub.Read()
	if err != nil {
		t.Fatal(err)
	}

	if string(msg.Body) != "hello" {
		t.Error("unexpected: " + string(msg.Body))
	}
}
