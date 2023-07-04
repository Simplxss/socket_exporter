package client

import (
	"encoding/binary"
	"math"
	"net"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/YuFanXing/socket_exporter/config"
)

type Client struct {
	endpoint config.Endpoint
	conn     net.Conn
	data     []byte

	Desc []*prometheus.Desc
}

func (c *Client) Connect() {
	switch c.endpoint.Type {
	case "tcp":
		tcpAddr, err := net.ResolveTCPAddr("tcp", c.endpoint.Address)
		if err != nil {
			println("ResolveTCPAddr failed:", err.Error())
			os.Exit(1)
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			println("Dial failed:", err.Error())
			os.Exit(1)
		}
		c.conn = conn
		go func() {
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					println("Read failed:", err.Error())
					continue
				}
				if n == c.endpoint.Length {
					c.data = buf
				}
			}
		}()
	case "udp":
		udpAddr, err := net.ResolveUDPAddr("udp", c.endpoint.Address)
		if err != nil {
			println("ResolveUDPAddr failed:", err.Error())
			os.Exit(1)
		}
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			println("Dial failed:", err.Error())
			os.Exit(1)
		}
		c.conn = conn
		go func() {
			for {
				buf := make([]byte, 2048)
				n, err := conn.Read(buf)
				if err != nil {
					println("Read failed:", err.Error())
					continue
				}
				if n == c.endpoint.Length {
					c.data = buf
				}
			}
		}()
	default:
		println("Unknown protocol:", c.endpoint.Type)
		os.Exit(1)
	}
}

func (c *Client) PrasePacket([]byte) {

}

func (c *Client) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.Desc {
		ch <- desc
	}
}

func (c *Client) Collect(ch chan<- prometheus.Metric) {
	if c.data != nil {
		for i, node := range c.endpoint.Protocol {
			switch node.Datatype {
			case "bool":
				var value float64
				if c.data[int(math.Floor(node.Offset))]&(1<<int(node.Offset*10)%10) == 0 {
					value = 0
				} else {
					value = 1
				}
				ch <- prometheus.MustNewConstMetric(c.Desc[i], prometheus.GaugeValue, value)
			case "real":
				ch <- prometheus.MustNewConstMetric(c.Desc[i], prometheus.GaugeValue, float64(math.Float32frombits(binary.BigEndian.Uint32(c.data[int(math.Floor(node.Offset)):int(math.Floor(node.Offset))+4]))))
			}
		}
	}
}

func NewClient(endpoint config.Endpoint, reg *prometheus.Registry) *Client {
	var client Client
	client.endpoint = endpoint
	client.Connect()

	for _, node := range endpoint.Protocol {
		Label := strings.Split(endpoint.Label, "=")
		client.Desc = append(client.Desc, prometheus.NewDesc(node.Name, node.Help, nil, prometheus.Labels{Label[0]: Label[1]}))
	}

	reg.MustRegister(&client)

	return &client
}
