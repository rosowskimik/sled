package trigger

import (
	"fmt"
	"os"
	"path/filepath"
)

type Timer struct {
	delayOn  uint
	delayOff uint
}

func NewTimer(delayOn, delayOff uint) *Timer {
	return &Timer{
		delayOn,
		delayOff,
	}
}

func (t *Timer) Setup(root string) error {
	onPath := filepath.Join(root, "delay_on")
	onFile, err := os.OpenFile(onPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer onFile.Close()

	offPath := filepath.Join(root, "delay_off")
	offFile, err := os.OpenFile(offPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer offFile.Close()

	if _, err := fmt.Fprintf(onFile, "%d", t.delayOn); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(offFile, "%d", t.delayOff); err != nil {
		return err
	}

	return nil

}

func (*Timer) Name() string {
	return "timer"
}

func (*Timer) Cleanup() {}
