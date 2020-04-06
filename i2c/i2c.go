package i2c

import (
	"fmt"
	"os"
	"syscall"
)

//I2C ..
type I2C struct {
	addr uint8
	bus  int
	file *os.File
}

//NewI2C ..
func NewI2C(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	_, _, errNo := syscall.Syscall6(syscall.SYS_IOCTL, f.Fd(), I2C_SLAVE, uintptr(addr), 0, 0, 0)
	if errNo != 0 {
		return nil, fmt.Errorf("Syscall to set address failed: %v", errNo)
	}

	v := &I2C{file: f, bus: bus, addr: addr}
	return v, nil
}

//GetBus ..
func (v *I2C) GetBus() int {
	return v.bus
}

//GetAddr ..
func (v *I2C) GetAddr() uint8 {
	return v.addr
}

//WriteBytes ..
func (v *I2C) WriteBytes(buf []byte) (int, error) {
	//log.Printf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return v.file.Write(buf)
}

//ReadBytes ..
func (v *I2C) ReadBytes(buf []byte) (int, error) {
	n, err := v.file.Read(buf)
	if err != nil {
		return n, err
	}
	//log.Printf("Read %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return n, nil
}

//Close and release file
func (v *I2C) Close() error {
	return v.file.Close()
}
