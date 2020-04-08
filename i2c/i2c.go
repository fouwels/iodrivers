package i2c

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

//I2C ..
type I2C struct {
	bus  int
	file *os.File
	sync.Mutex
}

//NewI2C ..
func NewI2C(bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// _, _, errNo := syscall.Syscall6(syscall.SYS_IOCTL, f.Fd(), I2C_SLAVE, uintptr(addr), 0, 0, 0)
	// if errNo != 0 {
	// 	return nil, fmt.Errorf("Syscall to set address failed: %v", errNo)
	// }

	v := &I2C{file: f, bus: bus}

	return v, nil
}

//SetAddr ..
func (v *I2C) SetAddr(addr uint8) error {
	v.Lock()
	defer v.Unlock()

	_, _, errNo := syscall.Syscall6(syscall.SYS_IOCTL, v.file.Fd(), I2C_SLAVE, uintptr(addr), 0, 0, 0)
	if errNo != 0 {
		return fmt.Errorf("Syscall to set address failed: %v", errNo)
	}

	return nil
}

//GetBus ..
func (v *I2C) GetBus() int {
	return v.bus
}

//WriteBytes ..
func (v *I2C) WriteBytes(buf []byte) (int, error) {
	v.Lock()
	defer v.Unlock()

	//log.Printf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return v.file.Write(buf)
}

//ReadBytes ..
func (v *I2C) ReadBytes(buf []byte) (int, error) {
	v.Lock()
	defer v.Unlock()

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
