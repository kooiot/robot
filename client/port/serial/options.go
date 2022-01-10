package serial

import (
	"go.bug.st/serial"
)

// Options
type Options struct {
	Port     string
	Baudrate int

	// Size of the character (must be 5, 6, 7 or 8)
	DataBits int

	// Parity is the bit to use and defaults to NoParity (no parity bit).
	Parity serial.Parity

	// Number of stop bits to use. Default is OneStopBit (1 stop bit).
	StopBits serial.StopBits
}

// Option ...
type Option func(*Options)

func newOptions(opt ...Option) *Options {
	opts := Options{
		Port:     "",
		Baudrate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Port == "" {
		opts.Port = "/dev/ttyS1"
	}

	return &opts
}

func Port(port string) Option {
	return func(o *Options) {
		o.Port = port
	}
}

func Baudrate(baud int) Option {
	return func(o *Options) {
		o.Baudrate = baud
	}
}

func DataBits(data_bits int) Option {
	return func(o *Options) {
		o.DataBits = data_bits
	}
}

func Parity(parity serial.Parity) Option {
	return func(o *Options) {
		o.Parity = parity
	}
}

func StopBits(stop_bits serial.StopBits) Option {
	return func(o *Options) {
		o.StopBits = stop_bits
	}
}
