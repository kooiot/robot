package hardware

import (
	"os/exec"
	"strconv"
)

type NamedLed struct {
	name string
}

func (s *NamedLed) Set(value int) error {
	cmd := "echo \"" + strconv.Itoa(value) + "\" > /sys/class/led/" + s.name + "/brightness"
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err
	}
	return nil
}

func (s *NamedLed) Get(value int) (int, error) {
	cmd := "cat /sys/class/led/" + s.name + "/brightness"
	out, err := exec.Command("sh", "-c", cmd).Output()
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
