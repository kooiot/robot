package hardware

import (
	"os/exec"
	"strconv"
)

type NamedGPIO struct {
	name string
}

func (s *NamedGPIO) Set(value int) error {
	_, err := exec.Command("echo", strconv.Itoa(value), ">", "/sys/class/gpio/"+s.name+"/value").Output()
	if err != nil {
		return err
	}
	return nil
}

func (s *NamedGPIO) Get(value int) (int, error) {
	out, err := exec.Command("cat", "/sys/class/gpio/"+s.name+"/value").Output()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(out))
}

func NewNamedGPIO(name string) *NamedGPIO {
	return &NamedGPIO{
		name: name,
	}
}
