package main

import (
	"log"

	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
)

func main() {
	pulse, e := pulseaudio.New()
	if e != nil {
		log.Panicln("connect", e)
	}

	client := &Client{pulse}
	pulse.Register(client)

	pulse.Listen()
}

type Client struct {
	client *pulseaudio.Client
}

func (c *Client) NewSink(path dbus.ObjectPath) {
	log.Println("new sink", path)
}

func (c *Client) SinkRemoved(path dbus.ObjectPath) {
	log.Println("sink removed", path)
}

func (c *Client) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	log.Println("device volume", path, values)
}

func (c *Client) DeviceMuteUpdated(path dbus.ObjectPath, muted bool) {
	log.Println("mute updated", path, muted)
}

func (c *Client) DeviceActivePortUpdated(path dbus.ObjectPath, sommat dbus.ObjectPath) {
	log.Println("active port updated", path, sommat)
}
