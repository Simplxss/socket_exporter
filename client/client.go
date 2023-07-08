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

const namespace = "socket"
const subsystem = "endpoint"

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
				if c.data[int(math.Floor(node.Offset))]&(1<<(int(node.Offset*10)%10)) == 0 {
					value = 0
				} else {
					if node.TrueValue == 0 {
						value = 1
					} else {
						value = float64(node.TrueValue)
					}
				}
				ch <- prometheus.MustNewConstMetric(c.Desc[i], prometheus.GaugeValue, value)
			case "real":
				ch <- prometheus.MustNewConstMetric(c.Desc[i], prometheus.GaugeValue, float64(math.Float32frombits(binary.BigEndian.Uint32(c.data[int(node.Offset):int(node.Offset)+4]))))
			}
		}
	}
}

func NewClient(endpoint config.Endpoint, reg *prometheus.Registry) *Client {
	var client Client
	client.endpoint = endpoint
	client.Connect()

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := client.conn.Read(buf)
			if err != nil {
				println("Read failed:", err.Error())
				if err.Error() == "EOF" {
					println("Reconnect:" + client.endpoint.Address)
					client.Connect()
				}
				continue
			}
			if n == client.endpoint.Length {
				client.data = buf
			}
		}
	}()

	for _, node := range endpoint.Protocol {
		Label := make(prometheus.Labels)
		if endpoint.Label != "" {
			Label1 := strings.Split(endpoint.Label, "=")
			Label[Label1[0]] = Label1[1]
		}
		if node.Label != "" {
			Label2 := strings.Split(node.Label, "=")
			Label[Label2[0]] = Label2[1]
		}

		client.Desc = append(client.Desc, prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, node.Name), node.Help, nil, Label))
	}

	reg.MustRegister(&client)

	return &client
}
