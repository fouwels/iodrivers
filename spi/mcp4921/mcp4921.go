package mcp4921

import (
	"encoding/binary"
	"fmt"

	"periph.io/x/periph/conn/spi"
)

//EnumBuffered ..
type EnumBuffered int

//EnumOutputGain ..
type EnumOutputGain int

//EnumShutdownMode ..
type EnumShutdownMode int

const (
	//EnumBufferedTrue Buffer the output registers
	EnumBufferedTrue EnumBuffered = 1
	//EnumBufferedFalse Do not buffer the output register
	EnumBufferedFalse EnumBuffered = 0

	//EnumOutputGain1x Set output to (1 * VREF) * input
	EnumOutputGain1x EnumOutputGain = 1
	//EnumOutputGain2x Set output to (2 * VREF) * input
	EnumOutputGain2x EnumOutputGain = 0

	//EnumShutdownModeActive Enable Vout when shut down
	EnumShutdownModeActive EnumShutdownMode = 1
	//EnumShutdownModeHighImpedence Set Vout to high impedence when shut down (500k typical)
	EnumShutdownModeHighImpedence EnumShutdownMode = 0
)

//Mcp4921 ..
type Mcp4921 struct {
	spi          spi.Conn
	label        string
	buffered     EnumBuffered
	outputGain   EnumOutputGain
	shutdownMode EnumShutdownMode
}

//NewMcp4921 ..
func NewMcp4921(spi spi.Conn, label string, buffered EnumBuffered, outputGain EnumOutputGain, shutdownMode EnumShutdownMode) (*Mcp4921, error) {

	mc := Mcp4921{
		spi:          spi,
		label:        label,
		buffered:     buffered,
		outputGain:   outputGain,
		shutdownMode: shutdownMode,
	}

	return &mc, nil
}

func (e *Mcp4921) Write(value uint16) error {

	if value >= 4096 {
		return fmt.Errorf("Input value out of range for a 12 bit ADC: %v", value)
	}

	buffer := uint16(0)
	buffer = value

	if e.shutdownMode == EnumShutdownModeActive {
		buffer = buffer | (1 << 12)
	}

	if e.outputGain == EnumOutputGain1x {
		buffer = buffer | (1 << 13)
	}

	if e.buffered == EnumBufferedTrue {
		buffer = buffer | (1 << 14)
	}

	tx := make([]byte, 2)
	binary.BigEndian.PutUint16(tx, buffer)

	rx := make([]byte, 2)
	err := e.spi.Tx(tx, rx)

	if err != nil {
		return err
	}

	return nil
}
