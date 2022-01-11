package hardware

import (
	"os/exec"
	"strconv"

	"github.com/kooiot/robot/pkg/util/log"
)

type NamedGPIO struct {
	name string
}

func (s *NamedGPIO) Set(value int) error {
	cmd := "echo \"" + strconv.Itoa(value) + "\" > /sys/class/gpio/" + s.name + "/value"
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err
	}
	// log.Info(string(out))
	return nil
}

func (s *NamedGPIO) Get(value int) (int, error) {
	cmd := "cat /sys/class/gpio/" + s.name + "/value"
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return -1, err
	}
	log.Info(string(out))
	return strconv.Atoi(string(out))
}

func NewNamedGPIO(name string) *NamedGPIO {
	return &NamedGPIO{
		name: name,
	}
}
