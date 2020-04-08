package mcp3208

import (
	"fmt"
	"math"

	"periph.io/x/periph/conn/spi"
)

const _maxChannels = 8

//Mcp3208 ..
type Mcp3208 struct {
	spi spi.Conn
}

//NewMcp3208 ..
func NewMcp3208(spi spi.Conn) (*Mcp3208, error) {

	mc := Mcp3208{
		spi: spi,
	}

	return &mc, nil
}

//GetValues retrieves n channels from channel chstart. Returns 12 bit output scaled to 16 bits.
func (m *Mcp3208) GetValues(chstart uint, n uint) ([]uint16, error) {

	values := []uint16{}

	for i := chstart; i < chstart+n; i++ {

		tx := make([]byte, 3)
		tx[0] = 0x06 + (byte(i) >> 2)
		tx[1] = (byte(i) & 0x03) << 6
		tx[2] = 0x00

		rx := make([]byte, 3)

		err := m.spi.Tx(tx, rx)

		if err != nil {
			return values, fmt.Errorf("Failed to transact SPI: %v", err)
		}
		if len(rx) != 3 {
			return values, fmt.Errorf("SPI response wrong length, got %v, expected %v: %v", len(rx), 3, err)
		}

		result12bit := int((rx[1]&0xf))<<8 + int(rx[2])
		result16bit := m.rescale12to16(result12bit)

		values = append(values, result16bit)
	}
	return values, nil
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
