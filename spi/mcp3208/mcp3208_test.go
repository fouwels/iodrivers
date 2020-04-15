package mcp3208

import (
	"log"
	"math"
	"os"
	"testing"

	"periph.io/x/periph/conn/physic"

	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const spiDevice = "/dev/spidev0.1"

var _mcp *Mcp3208

func TestMain(m *testing.M) {

	_, err := host.Init()
	if err != nil {
		log.Fatalf("Failed to init: %v", err)
	}

	s, err := spireg.Open(spiDevice)
	if err != nil {
		log.Fatalf("Failed to create SPI: %v", err)
	}

	conn, err := s.Connect(physic.Frequency(1*physic.MegaHertz), spi.Mode0, 8)
	if err != nil {
		log.Fatalf("Failed to connect SPI: %v", err)
	}

	mcp, err := NewMcp3208(conn, "ADC1")
	if err != nil {
		log.Fatalf("Failed to create MCP: %v", err)
	}

	_mcp = mcp

	result := m.Run()

	err = s.Close()
	if err != nil {
		log.Fatalf("Failed to close SPI: %v", err)
	}
	os.Exit(result)
}

func TestRead(t *testing.T) {

	for i := 0; i < 5; i++ {
		vals, tstamp, err := _mcp.GetValues(0, 8)
		if err != nil {
			log.Fatalf("Failed: %v", err)
		}

		for i, v := range vals {
			log.Printf("%v [%v] %v", tstamp, i, float64(v)*(1/math.Pow(2, 16)))
		}
	}
}

func TestRescale(t *testing.T) {
	out := _mcp.rescale12to16(0)
	if out != 0 {
		log.Fatalf("Failed, got %v, expected %v", out, 0)
	}
	out = _mcp.rescale12to16(4095)
	if out != 65535 {
		log.Fatalf("Failed,got %v, expected %v", out, 65535)
	}

	out = _mcp.rescale12to16(2000)
	if out != 32007 {
		log.Fatalf("Failed,got %v, expected %v", out, 32007)
	}
}

func BenchmarkRescale(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = _mcp.rescale12to16(0)
		_ = _mcp.rescale12to16(2000)
		_ = _mcp.rescale12to16(3000)
		_ = _mcp.rescale12to16(4000)
		_ = _mcp.rescale12to16(65535)
	}
}

func TestChannelToBitMask(t *testing.T) {
	a, b, c := _mcp.channelToBitmask(0)
	if !(a == false && b == false && c == false) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(1)
	if !(a == false && b == false && c == true) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(2)
	if !(a == false && b == true && c == false) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(3)
	if !(a == false && b == true && c == true) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(4)
	if !(a == true && b == false && c == false) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(5)
	if !(a == true && b == false && c == true) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(6)
	if !(a == true && b == true && c == false) {
		log.Fatalf("Failed")
	}

	a, b, c = _mcp.channelToBitmask(7)
	if !(a == true && b == true && c == true) {
		log.Fatalf("Failed")
	}
}
