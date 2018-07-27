package uploader

import (
	"time"
	"kuber-controller/buffering"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"net/http"
	"bytes"
)

const READ_SIZE uint32 = 50

type Payload struct {
	Key          string
	EventType    string
	Namespace    string
	ResourceType string
	Data         *interface{}
}

func UploadData(ringBuffer *buffering.RingBuffer) {
	for true {
		ringBuffer.PrintDetails()

		for true {
			data, size := ringBuffer.ReadN(READ_SIZE)

			if size == 0 {
				log.Debug("There is no data to upload.")
				break
			}

			resp, err := SendData(data)
			if err != nil {
				log.Warn("Error while sending data ", err)
				break
			}

			if resp.StatusCode == 200 {
				log.Info("Data is posted successfully")
				ringBuffer.RemoveN(size)
				ringBuffer.PrintDetails()
			} else {
				log.Error("Data posting is failed with error ", resp.StatusCode)
				break
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func SendData(payload []*interface{}) (*http.Response, error) {
	url := "http://localhost:8002/le-mans/v1/streams/purser-stream"
	jsonStr, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer lGCVEpnJavdaq49qiQEH5pxwGj5wo9ZQ")

	client := &http.Client{}
	return client.Do(req)
}
