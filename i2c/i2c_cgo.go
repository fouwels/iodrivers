// +build linux,cgo

package i2c

// #include <linux/i2c-dev.h>
import "C"

const (
	I2C_SLAVE = C.I2C_SLAVE
)
