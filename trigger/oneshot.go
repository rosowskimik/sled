package trigger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	invertTrue  = "1"
	invertFalse = "0"
)

type Oneshot struct {
	delayOn  time.Duration
	delayOff time.Duration
	invert   bool
	c        chan<- interface{}
}

func NewOneshot(delayOn, delayOff time.Duration, invert bool) *Oneshot {
	return &Oneshot{
		delayOn:  delayOn,
		delayOff: delayOff,
		invert:   invert,
		c:        nil,
	}
}

func (o *Oneshot) Shoot() error {
	if o.c == nil {
		return errors.New("Attempt to shoot on uninitialized Oneshot trigger")
	}
	o.c <- nil
	return nil
}

func (o *Oneshot) Setup(root string) error {
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

	invertPath := filepath.Join(root, "invert")
	invertFile, err := os.OpenFile(invertPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer offFile.Close()

	if _, err := fmt.Fprintf(onFile, "%d", ledDelay(o.delayOn)); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(offFile, "%d", ledDelay(o.delayOff)); err != nil {
		return err
	}

	var inv string
	if o.invert {
		inv = invertTrue
	} else {
		inv = invertFalse
	}
	if _, err := invertFile.WriteString(inv); err != nil {
		return err
	}

	shotPath := filepath.Join(root, "shot")
	shotFile, err := os.OpenFile(shotPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	c := make(chan interface{})
	go func(f *os.File, c <-chan interface{}) {
		for range c {
			f.WriteString("1")
		}
		f.Close()
	}(shotFile, c)
	o.c = c

	return nil
}

func (*Oneshot) Name() string {
	return "oneshot"
}

func (o *Oneshot) Cleanup() {
	if o.c != nil {
		close(o.c)
		o.c = nil
	}
}
