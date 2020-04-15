package sfm3000

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

const i2cBus = "1"
const i2cAddress = 0x40 //64

var _sfm *SFM3000

func TestMain(m *testing.M) {

	_, err := host.Init()
	if err != nil {
		log.Fatalf("Failed to init: %v", err)
	}

	bus, err := i2creg.Open(i2cBus)
	if err != nil {
		log.Fatalf("Failed to create SPI: %v", err)
	}

	dev := i2c.Dev{
		Bus:  bus,
		Addr: i2cAddress,
	}

	sfm, err := NewSFM3000(&dev, i2cAddress, true, "TESTSENSOR")
	if err != nil {
		log.Fatalf("Failed to create SFM3000: %v", err)
	}

	_sfm = sfm

	result := m.Run()

	err = bus.Close()
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
		value, crc, _, err := _sfm.getRaw()
		if err != nil {
			t.Fatalf("[%v] Failed to get value: %v", i, err)
		}
		t.Logf("[%v] Value: %v CRC: %v", i, value, crc)
	}
}

func TestGetValue(t *testing.T) {

	for i := 0; i < 6; i++ {
		value, crc, _, err := _sfm.GetValue()
		if err != nil {
			t.Fatalf("[%v] Failed to get value: %v", i, err)
		}
		t.Logf("[%v] Value: %v CRC: %v", i, value, crc)
	}
}

func TestCaptureDatalog(t *testing.T) {

	const fileName string = "capture_datalog.csv"
	const sampleRate int = 1000

	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		t.Fatalf("Failed to crete+open file: %v", err)
	}

	startTime := time.Now()

	cs := csv.NewWriter(f)
	defer cs.Flush()

	lines := [][]string{}

	t.Logf("Starting datalog for 10 seconds at %v", sampleRate)
	for i := 0; i < 10*sampleRate; i++ {

		value, _, timestamp, err := _sfm.GetValue()
		if err != nil {
			t.Fatalf("[%v] Failed to get value: %v", i, err)
		}

		//Save a offset to force go to use monotonic time...
		line := []string{timestamp.Sub(startTime).String(), fmt.Sprintf("%v", value)}

		lines = append(lines, line)
		time.Sleep((1 * time.Second) / time.Duration(sampleRate))
	}

	t.Logf("Captured %v records", len(lines))

	cs.WriteAll(lines)
	cs.Flush()
	if err := cs.Error(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}

	t.Logf("Finished datalog")
}
