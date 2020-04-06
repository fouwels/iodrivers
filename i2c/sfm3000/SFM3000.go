package sfm3000

import (
	"encoding/binary"
	"fmt"

	"github.com/go-daq/crc8"
	"github.com/kaelanfouwels/iodrivers/i2c"
)

const crcPolynomial = 0x31
const flowOffset = 32000
const flowScaleFactorAirN2 = 140
const flowScaleFactorO2 = 142.8

//SFM3000 is the i2C driver for the SFM3000 Low Pressure Drop Digital Flow Meter
type SFM3000 struct {
	i2c      *i2c.I2C
	crcTable *crc8.Table
	readMode bool
	isAir    bool
}

//NewSFM3000 create a new SFM3000 driver
func NewSFM3000(i2c *i2c.I2C, isAir bool) (*SFM3000, error) {

	return &SFM3000{
		i2c:      i2c,
		readMode: false,
		crcTable: crc8.MakeTable(crcPolynomial),
		isAir:    isAir,
	}, nil
}

//SoftReset ..
func (e *SFM3000) SoftReset() error {

	e.readMode = false

	w := []byte{0x20, 0x00}
	_, err := e.i2c.WriteBytes(w)
	if err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	return nil
}

//GetSerial ..
func (e *SFM3000) GetSerial() ([4]byte, error) {

	e.readMode = false

	serial := [4]byte{}

	w := []byte{0x31, 0xAE}

	_, err := e.i2c.WriteBytes(w)
	if err != nil {
		return serial, fmt.Errorf("failed to write command: %w", err)
	}

	r := make([]byte, 4, 4)
	_, err = e.i2c.ReadBytes(r)
	if err != nil {
		return serial, fmt.Errorf("failed to read response: %w", err)
	}

	if len(r) != 4 {
		return serial, fmt.Errorf("response length unexpected (bytes), got %v, expected %v", len(r), 4)
	}

	copy(serial[:], r)

	return serial, nil
}

//GetValue Returns data uint16, crc uint8, error
func (e *SFM3000) GetValue() (float32, uint8, error) {

	value, crc, err := e.getRaw()
	if err != nil {
		return 0, 0, err
	}

	var scalefactor float32
	if e.isAir {
		scalefactor = flowScaleFactorAirN2
	} else {
		scalefactor = flowScaleFactorO2
	}

	flow := (float32(value) - flowOffset) / (scalefactor)
	return flow, crc, nil
}

//getRaw Returns data uint16, crc uint8, error
func (e *SFM3000) getRaw() (uint16, uint8, error) {

	if !(e.readMode) {

		w := []byte{0x10, 00}

		_, err := e.i2c.WriteBytes(w)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to write command: %w", err)
		}

		e.readMode = true
	}

	r := make([]byte, 3, 3)
	_, err := e.i2c.ReadBytes(r)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read command: %w", err)
	}

	dataCRC := byte(crc8.Checksum(r[:2], e.crcTable))
	sensorCRC := r[2]

	if dataCRC != sensorCRC {
		return 0, 0, fmt.Errorf("CRC Check failed, got %v, expected %v", sensorCRC, dataCRC)
	}

	data := binary.BigEndian.Uint16(r[:2])

	return data, dataCRC, nil
}
