package main

import (
	"fmt"
	"log"
	"time"
)

/*
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <linux/types.h>
#include <linux/spi/spidev.h>

#define SPI_SPEED 4000000

uint8_t mode=0;
uint8_t bits=8;
uint32_t speed=SPI_SPEED;
uint16_t delay=5;

int spi_open(const char *device) {
  int fd = open(device, O_RDWR);
  int ret;

  if (fd < 0) {
    printf("can't open device");
    return -1;
  }
  ret = ioctl(fd, SPI_IOC_WR_MODE, &mode);
  if (ret == -1) {
    printf("can't set spi mode");
    return -1;
  }

  ret = ioctl(fd, SPI_IOC_RD_MODE, &mode);
  if (ret == -1) {
    printf("can't get spi mode");
    return -1;
  }

  ret = ioctl(fd, SPI_IOC_WR_BITS_PER_WORD, &bits);
  if (ret == -1) {
    printf("can't set bits per word");
    return -1;
  }

  ret = ioctl(fd, SPI_IOC_RD_BITS_PER_WORD, &bits);
  if (ret == -1) {
    printf("can't get bits per word");
    return -1;
  }

  ret = ioctl(fd, SPI_IOC_WR_MAX_SPEED_HZ, &speed);
  if (ret == -1) {
    printf("can't set max speed hz");
    return -1;
  }

  ret = ioctl(fd, SPI_IOC_RD_MAX_SPEED_HZ, &speed);
  if (ret == -1) {
    printf("can't get max speed hz");
    return -1;
  }

  return fd;
}

int spi_xfer(int fd, char* tx, char* rx, int length) {
  struct spi_ioc_transfer tr = {
    .tx_buf = (unsigned long)tx,
    .rx_buf = (unsigned long)rx,
    .len = length,
    .delay_usecs = delay,
    .speed_hz = speed,
    .bits_per_word = bits,
  };

  int ret = ioctl(fd, SPI_IOC_MESSAGE(1), &tr);
  if (ret < 1)
    return -1;

  return 0;
}
*/
import "C"
import "unsafe"

import "errors"

// SPIDevice device
type SPIDevice struct {
	fd C.int
}

// NewSPIDevice opens the device
func NewSPIDevice(devPath string) (*SPIDevice, error) {
	name := C.CString(devPath)
	defer C.free(unsafe.Pointer(name))
	i := C.spi_open(name)
	if i < 0 {
		return nil, errors.New("could not open")
	}
	return &SPIDevice{i}, nil
}

// Xfer cross transfer
func (d *SPIDevice) Xfer(tx []byte) ([]byte, error) {
	length := len(tx)
	rx := make([]byte, length)
	ret := C.spi_xfer(
		d.fd,
		(*C.char)(unsafe.Pointer(&tx[0])),
		(*C.char)(unsafe.Pointer(&rx[0])),
		C.int(length))
	if ret < 0 {
		return nil, errors.New("could not xfer")
	}
	return rx, nil
}

// Close closes the fd
func (d *SPIDevice) Close() {
	C.close(d.fd)
}

type Max7219Reg byte

const (
	MAX7219_REG_NOOP   Max7219Reg = 0
	MAX7219_REG_DIGIT0            = iota
	MAX7219_REG_DIGIT1
	MAX7219_REG_DIGIT2
	MAX7219_REG_DIGIT3
	MAX7219_REG_DIGIT4
	MAX7219_REG_DIGIT5
	MAX7219_REG_DIGIT6
	MAX7219_REG_DIGIT7
	MAX7219_REG_DECODEMODE
	MAX7219_REG_INTENSITY
	MAX7219_REG_SCANLIMIT
	MAX7219_REG_SHUTDOWN
	MAX7219_REG_DISPLAYTEST = 0x0F
	MAX7219_REG_LASTDIGIT   = MAX7219_REG_DIGIT7
)

const MAX7219_DIGIT_COUNT = MAX7219_REG_LASTDIGIT -
	MAX7219_REG_DIGIT0 + 1

const cascaded = 4

type Device struct {
	buffer []byte
	spi    *SPIDevice
}

func NewDevice() *Device {
	buf := make([]byte, MAX7219_DIGIT_COUNT*cascaded)
	this := &Device{buffer: buf}
	return this
}

func (this *Device) Open(spibus int, spidevice int, brightness byte) error {
	devstr := fmt.Sprintf("/dev/spidev%d.%d", spibus, spidevice)
	spi, err := NewSPIDevice(devstr)
	if err != nil {
		return err
	}
	this.spi = spi
	// Initialize Max7219 led driver.
	this.Command(MAX7219_REG_SCANLIMIT, 7)   // show all 8 digits
	this.Command(MAX7219_REG_DECODEMODE, 0)  // use matrix (not digits)
	this.Command(MAX7219_REG_DISPLAYTEST, 0) // no display test
	this.Command(MAX7219_REG_SHUTDOWN, 1)    // not shutdown mode
	this.Command(MAX7219_REG_INTENSITY, brightness)
	this.ClearAll()
	this.Flush()
	return nil
}

func (this *Device) Close() {
	this.spi.Close()
}

func (this *Device) Command(reg Max7219Reg, value byte) error {
	buf := []byte{byte(reg), value}
	for i := 0; i < cascaded; i++ {
		_, err := this.spi.Xfer(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Device) SetBufferLine(
	cascadeId int, position int, value byte) {
	this.buffer[cascadeId*MAX7219_DIGIT_COUNT+position] = value
}

func (this *Device) Flush() {
	buf := make([]byte, cascaded*2)
	for position := 0; position < MAX7219_DIGIT_COUNT; position++ {
		for i := 0; i < cascaded; i++ {
			buf[i*2] = byte(MAX7219_REG_DIGIT0 + position)
			buf[i*2+1] = this.buffer[i*MAX7219_DIGIT_COUNT+position]
		}
		this.spi.Xfer(buf)
	}
}

func (this *Device) ClearAll() {
	for cId := 0; cId < cascaded; cId++ {
		for i := 0; i < MAX7219_DIGIT_COUNT; i++ {
			this.buffer[cId*MAX7219_DIGIT_COUNT+i] = 0
		}
	}
}

func main() {
	device := NewDevice()
	err := device.Open(0, 0, 7)
	if err != nil {
		log.Fatal(err)
	}
	defer device.Close()
	var b byte = 0x55
	for t := 1; t < 100; t++ {
		for l := 0; l < 8; l++ {
			for i := 0; i < 4; i++ {
				device.SetBufferLine(i, l, b)
			}
			b ^= 0xFF
		}
		device.Flush()
		b ^= 0xFF
		time.Sleep(200 * time.Millisecond)
	}
}
