package mcp4921_test

import (
	"log"
	"os"
	"testing"
	"time"

	"periph.io/x/periph/conn/physic"

	"github.com/kaelanfouwels/iodrivers/spi/mcp4921"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const spiDevice = "/dev/spidev0.0"

var _mcp *mcp4921.Mcp4921

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

	mcp, err := mcp4921.NewMcp4921(conn, "DAC1", mcp4921.EnumBufferedTrue, mcp4921.EnumOutputGain1x, mcp4921.EnumShutdownModeActive)
	if err != nil {
		log.Fatalf("Failed to create DAC: %v", err)
	}

	_mcp = mcp

	result := m.Run()

	err = s.Close()
	if err != nil {
		log.Fatalf("Failed to close SPI: %v", err)
	}
	os.Exit(result)
}

func TestWrite(t *testing.T) {

	err := _mcp.Write((4096 / 100) * 14) //2^12/2
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
}

func TestVolume(t *testing.T) {

	duration := 1000 * time.Millisecond
	percentage := 10.00

	timer := time.NewTimer(duration)

	err := _mcp.Write(uint16((4096.00 / 100.00) * percentage))
	if err != nil {
		log.Fatalf("Failed to open: %v", err)
	}

	<-timer.C

	err = _mcp.Write(0)
	if err != nil {
		log.Fatalf("Failed to close: %v", err)
	}

}

func TestWriteOOR(t *testing.T) {
	err := _mcp.Write(4096)
	if err == nil {
		log.Fatalf("Test failed, OOR value was not rejected")
	}
}
