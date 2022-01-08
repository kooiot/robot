package hardware

import (
	"os/exec"
	"strconv"
)

type NamedLed struct {
	name string
}

func (s *NamedLed) Set(value int) error {
	_, err := exec.Command("echo", strconv.Itoa(value), ">", "/sys/class/led/"+s.name+"/brightness").Output()
	if err != nil {
		return err
	}
	return nil
}

func (s *NamedLed) Get(value int) (int, error) {
	out, err := exec.Command("cat", "/sys/class/led/"+s.name+"/brightness").Output()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(out))
}

func NewNamedLed(name string) *NamedLed {
	return &NamedLed{
		name: name,
	}
}
