package sled

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rosowskimik/sled/trigger"
)

const LED_BASE = "/sys/class/leds"

type LED struct {
	root          string
	maxBrightness uint
	bFile         *os.File
	trigger       trigger.Trigger
	lock          sync.Mutex
}

func New(name string) (*LED, error) {
	root := filepath.Join(LED_BASE, name)

	mbPath := filepath.Join(root, "max_brightness")
	mbFile, err := os.Open(mbPath)
	if err != nil {
		return nil, err
	}
	defer mbFile.Close()

	mbBytes, err := io.ReadAll(mbFile)
	if err != nil {
		return nil, err
	}

	maxBrightness, err := strconv.ParseUint(strings.TrimSpace(string(mbBytes)), 10, strconv.IntSize)
	if err != nil {
		return nil, err
	}

	bPath := filepath.Join(root, "brightness")
	bFile, err := os.OpenFile(bPath, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	led := &LED{
		root:          root,
		maxBrightness: uint(maxBrightness),
		bFile:         bFile,
		trigger:       nil,
		lock:          sync.Mutex{},
	}

	return led, nil
}

func (l *LED) Close() error {
	if !l.lock.TryLock() {
		return errors.New("LED still in use, can't close")
	}
	defer l.lock.Unlock()

	l.ClearTrigger()

	return l.bFile.Close()
}

func (l *LED) MaxBrightness() uint {
	return l.maxBrightness
}

func (l *LED) GetBrightness() (uint, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, err := l.bFile.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	bytes, err := io.ReadAll(l.bFile)
	if err != nil {
		return 0, err
	}

	brightness, err := strconv.ParseUint(strings.TrimSpace(string(bytes)), 10, strconv.IntSize)
	if err != nil {
		return 0, err
	}

	return uint(brightness), nil
}

func (l *LED) SetBrightness(b uint) error {
	if b > l.maxBrightness {
		return errors.New(fmt.Sprintf("Bad brightness value: %v > %v", b, l.maxBrightness))
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	if _, err := l.bFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(l.bFile, "%d", b); err != nil {
		return err
	}

	return nil
}

func (l *LED) GetTrigger() trigger.Trigger {
	return l.trigger
}

func (l *LED) SetTrigger(t trigger.Trigger) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.trigger != nil {
		l.trigger.Cleanup()
	}

	tPath := filepath.Join(l.root, "trigger")
	tFile, err := os.OpenFile(tPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(tFile, "%s", t.Name()); err != nil {
		return err
	}

	if err := t.Setup(l.root); err != nil {
		return err
	}

	l.trigger = t
	return nil
}

func (l *LED) ClearTrigger() error {
	if l.trigger != nil {
		l.trigger.Cleanup()
	}
	if err := l.SetBrightness(0); err != nil {
		return err
	}

	l.trigger = nil
	return nil
}
