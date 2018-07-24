package uploader

import (
	"time"
	"kuber-controller/buffering"
	log "github.com/Sirupsen/logrus"
)

const READ_SIZE uint32 = 50

type Payload struct {
	Key          string
	EventType    string
	Namespace    string
	ResourceType string
	Data		*interface{}
}

func UploadData(ringBuffer *buffering.RingBuffer)  {
	for true {
		ringBuffer.PrintDetails()
		log.Println("Uploading data")
		_, size := ringBuffer.ReadN(READ_SIZE)
		log.Printf("Read %d from buffer\n", size)
		ringBuffer.RemoveN(size)
		ringBuffer.PrintDetails()
		time.Sleep(10 * time.Second)
	}
}