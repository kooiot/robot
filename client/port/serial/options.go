package serial

import (
	"time"
)

// Options
type Options struct {
	Port     string
	Baudrate int

	// Size is the number of data bits. If 0, DefaultSize is used.
	DataBits byte

	// Parity is the bit to use and defaults to ParityNone (no parity bit).
	Parity byte

	// Number of stop bits to use. Default is 1 (1 stop bit).
	StopBits byte

	// Total timeout
	ReadTimeout time.Duration
}

// Option ...
type Option func(*Options)

func newOptions(opt ...Option) *Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Port == "" {
		opts.Port = "/dev/ttyS1"
	}
	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = time.Millisecond * 500
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

func DataBits(data_bits byte) Option {
	return func(o *Options) {
		o.DataBits = data_bits
	}
}

func Parity(parity byte) Option {
	return func(o *Options) {
		o.Parity = parity
	}
}

func StopBits(stop_bits byte) Option {
	return func(o *Options) {
		o.StopBits = stop_bits
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = timeout
	}
}
