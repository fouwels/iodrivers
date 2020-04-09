package mcp3208

import (
	"fmt"
	"math"
	"time"

	"periph.io/x/periph/conn/spi"
)

const _maxChannels = 8

//Mcp3208 ..
type Mcp3208 struct {
	spi   spi.Conn
	label string
}

//NewMcp3208 ..
func NewMcp3208(spi spi.Conn, label string) (*Mcp3208, error) {

	mc := Mcp3208{
		spi:   spi,
		label: label,
	}

	return &mc, nil
}

//Label ..
func (m *Mcp3208) Label() string {
	return m.label
}

//GetValues retrieves n channels from channel chstart. Returns 12 bit output scaled to 16 bits.
func (m *Mcp3208) GetValues(chstart uint, n uint) ([]uint16, time.Time, error) {

	values := []uint16{}

	for i := chstart; i < chstart+n; i++ {

		tx := make([]byte, 3)
		tx[0] = 0x06 + (byte(i) >> 2)
		tx[1] = (byte(i) & 0x03) << 6
		tx[2] = 0x00

		rx := make([]byte, 3)

		err := m.spi.Tx(tx, rx)

		if err != nil {
			return []uint16{}, time.Time{}, fmt.Errorf("Failed to transact SPI: %v", err)
		}
		if len(rx) != 3 {
			return []uint16{}, time.Time{}, fmt.Errorf("SPI response wrong length, got %v, expected %v: %v", len(rx), 3, err)
		}
		result12bit := int((rx[1]&0xf))<<8 + int(rx[2])
		result16bit := m.rescale12to16(result12bit)

		values = append(values, result16bit)
	}

	timestamp := time.Now()
	return values, timestamp, nil
}

func (m *Mcp3208) rescale12to16(in int) uint16 {

	con := (float64(math.Pow(2, 16)-1) / float64(math.Pow(2, 12)-1))

	return uint16(con * float64(in))
}

func (m *Mcp3208) channelToBitmask(channel uint) (D2 bool, D1 bool, D0 bool) {

	D2 = (channel / 4) == 1
	D1 = ((channel % 2) / 2) != 0
	D0 = ((channel % 2) != 0) && channel != 0

	return
}
