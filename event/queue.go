package event

import (
	"log"
	"os"
	"path"
	"tinylytics/config"
	"tinylytics/constants"

	"github.com/joncrlsn/dque"
)

type EventQueue struct {
	queue *dque.DQue
}

// ItemBuilder creates a new item and returns a pointer to it.
// This is used when we load a segment of the queue from disk.
func ItemBuilder() interface{} {
	return &ClientInfo{}
}

func (q *EventQueue) GetSize() int {
	return q.queue.Size()
}

func (q *EventQueue) Push(item *ClientInfo) {
	if err := q.queue.Enqueue(item); err != nil {
		log.Fatal("Error enqueueing item ", err)
	}
}

func (q *EventQueue) Pop() *ClientInfo {
	var iface interface{}
	var err interface{}

	if iface, err = q.queue.DequeueBlock(); err != nil {
		log.Fatal("Error dequeuing item ", err)
	}

	// Assert type of the response to an Item pointer so we can work with it
	item, ok := iface.(*ClientInfo)
	if !ok {
		log.Fatal("Dequeued object is not an Item pointer")
	}

	return item
}

func (q *EventQueue) Peek() *ClientInfo {
	var iface interface{}
	var err interface{}

	if iface, err = q.queue.PeekBlock(); err != nil {
		log.Fatal("Error dequeuing item ", err)
	}

	// Assert type of the response to an Item pointer so we can work with it
	item, ok := iface.(*ClientInfo)
	if !ok {
		log.Fatal("Dequeued object is not an Item pointer")
	}

	return item
}

// ExampleDQue shows how the queue works
func (q *EventQueue) Connect() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	qDir := path.Join(wd, config.Config.DataFolder)

	segmentSize := 50

	// Create a new queue with segment size of 50
	s, err := dque.NewOrOpen(constants.EVENT_QUEUE_NAME, qDir, segmentSize, ItemBuilder)
	if err != nil {
		log.Fatal("Error creating new dque ", err)
	}
	q.queue = s
}

// Listen processes events from the queue synchronously (one at a time).
// Events are processed in order: peek, handle, pop, repeat.
// This ensures no parallel processing of queue events.
func (q *EventQueue) Listen(handler func(item *ClientInfo)) {
	log.Println("Queue listener started - waiting for events...")
	go func() {
		for {
			b := q.Peek()
			log.Printf("[QUEUE] Picked up event: domain=%s page=%s IP=%s UserAgent=%s", b.Domain, b.Page, b.IP, b.UserAgent)
			// Process event synchronously - handler must complete before next event
			handler(b)
			q.Pop()
		}
	}()
}
