package main

import (
	"fmt"
	"net"

	"time"

	"github.com/alecthomas/kingpin"
	"github.com/nickw444/miio-go"
	"github.com/nickw444/miio-go/device"
	"github.com/nickw444/miio-go/protocol"

	"github.com/nickw444/miio-go/common"
	"github.com/nickw444/miio-go/subscription"
	"github.com/sirupsen/logrus"
)

func onNewDevice(dev common.Device) {
	switch dev.(type) {
	case *device.Yeelight:
		fmt.Printf("Found Yeelight :)\n")
	case *device.PowerPlug:
		fmt.Println("Found PowerPlug")
		d := dev.(*device.PowerPlug)
		sub, err := d.NewSubscription()
		if err != nil {
			panic(err)
		}
		go watchSubscription(sub)
		go tick(d)

	default:
		fmt.Printf("Unknown device type %T\n", dev)
	}
}

func watchSubscription(sub subscription.Subscription) {
	for event := range sub.Events() {
		fmt.Printf("New Sub Event: %T\n", event)
	}
}

func tick(d *device.PowerPlug) {
	currState := false
	for {
		select {
		case <-time.After(time.Second * 5):
			var s common.PowerState
			if currState {
				s = common.PowerStateOn
			} else {
				s = common.PowerStateOff
			}
			currState = !currState
			d.SetPower(s)
		}
	}
}

func main() {
	var (
		local = kingpin.Flag("local", "Send broadcast to 127.0.0.1 instead of 255.255.255.255 (For use with locally hosted simulator)").Bool()
	)

	kingpin.Parse()

	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	common.SetLogger(l)

	addr := net.IPv4bcast
	if *local {
		addr = net.IPv4(127, 0, 0, 1)
	}

	proto, err := protocol.NewProtocol(protocol.ProtocolConfig{
		BroadcastIP: addr,
	})
	if err != nil {
		panic(err)
	}

	client, err := miio.NewClientWithProtocol(proto)
	if err != nil {
		panic(err)
	}

	sub, err := client.NewSubscription()
	if err != nil {
		panic(err)
	}

	for event := range sub.Events() {
		switch event.(type) {
		case common.EventNewDevice:
			dev := event.(common.EventNewDevice).Device
			onNewDevice(dev)
			fmt.Printf("New device event %T\n", dev)
		case common.EventExpiredDevice:
			fmt.Println("Expired device event")
		default:
			fmt.Println("Uknown event")
		}
	}
}
