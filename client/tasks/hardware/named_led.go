package hardware

import (
	"os/exec"
	"strconv"

	"github.com/kooiot/robot/pkg/util/log"
)

type NamedLed struct {
	name string
}

func (s *NamedLed) Set(value int) error {
	cmd := "echo \"" + strconv.Itoa(value) + "\" > /sys/class/leds/" + s.name + "/brightness"
	// log.Info(cmd)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Error("Echo led error: %s", err.Error())
		return err
	}
	return nil
}

func (s *NamedLed) Get(value int) (int, error) {
	cmd := "cat /sys/class/leds/" + s.name + "/brightness"
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
