package main

import (
	"kuber-controller/controller"
	"kuber-controller/config"
	log "github.com/Sirupsen/logrus"
	"os"
	"kuber-controller/buffering"
	"sync"
	"kuber-controller/client"
	"kuber-controller/uploader"
)

var conf *config.Config

var groupcrdclient *client.GroupCrdClient
var subscriberclient *client.SubscriberCrdClient

func init() {
	setlogFile()
	conf = &config.Config{}
	//conf.Resource = config.Resource{Pod: true, Node:true, Services:true, ReplicaSet:true, Deployment:true, Job:true}
	conf.Resource = config.Resource{Pod: true}
	conf.RingBuffer = &buffering.RingBuffer{Size: buffering.BUFFER_SIZE, Mutex: &sync.Mutex{}}
	// initialize client for api extension server
	conf.Groupcrdclient, conf.Subscriberclient = controller.GetApiExtensionClient()
}

func setlogFile() {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	//log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	go uploader.UploadData(conf)
	controller.Start(conf)
	//controller.CreateCRDDefinition()
	// controller.TestCrdFlow()
}
