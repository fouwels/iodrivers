package sfm3000

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
)

const i2cBus = "1"
const i2cAddress = 0x40 //64

var _sfm *SFM3000

func TestMain(m *testing.M) {

	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	i2c, err := i2c.NewI2C(i2cAddress, 1)
	if err != nil {
		log.Fatalf("Failed to create I2C: %v", err)
	}

	sfm, err := NewSFM3000(i2c, true)
	if err != nil {
		log.Fatalf("Failed to create SFM3000: %v", err)
	}

	_sfm = sfm

	result := m.Run()

	err = i2c.Close()
	if err != nil {
		log.Fatalf("Failed to close I2C: %v", err)
	}
	os.Exit(result)
}

func TestSoftReset(t *testing.T) {
	err := _sfm.SoftReset()
	if err != nil {
		t.Fatalf("Failed to soft reset: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
}

func TestGetSerial(t *testing.T) {
	serial, err := _sfm.GetSerial()
	if err != nil {
		t.Fatalf("Failed to get serial no: %v", err)
	}

	t.Logf("Serial number: %v", serial)
}

func TestGetRaw(t *testing.T) {

	for i := 0; i < 6; i++ {
		value, crc, err := _sfm.getRaw()
		if err != nil {
			t.Fatalf("[%v] Failed to get value: %v", i, err)
		}
		t.Logf("[%v] Value: %v CRC: %v", i, value, crc)
	}
}

func TestGetValue(t *testing.T) {

	for i := 0; i < 6; i++ {
		value, crc, err := _sfm.GetValue()
		if err != nil {
			t.Fatalf("[%v] Failed to get value: %v", i, err)
		}
		t.Logf("[%v] Value: %v CRC: %v", i, value, crc)
	}
}
