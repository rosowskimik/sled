package trigger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	repeatTrue  = "-1"
	repeatFalse = "1"
)

type Step struct {
	Brightness uint
	Duration   time.Duration
}

func NewStep(brightness uint, duration time.Duration) Step {
	return Step{
		Brightness: brightness,
		Duration:   duration,
	}
}

type Pattern struct {
	steps  []Step
	repeat bool
}

func NewPattern(repeat bool, steps ...Step) *Pattern {
	return &Pattern{
		steps,
		repeat,
	}
}

func (p *Pattern) Setup(root string) error {
	patternPath := filepath.Join(root, "pattern")
	patternFile, err := os.OpenFile(patternPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer patternFile.Close()

	repeatPath := filepath.Join(root, "repeat")
	repeatFile, err := os.OpenFile(repeatPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer repeatFile.Close()

	var rStr string
	if p.repeat {
		rStr = repeatTrue
	} else {
		rStr = repeatFalse
	}
	if _, err := repeatFile.WriteString(rStr); err != nil {
		return err
	}

	var pStr string
	for _, step := range p.steps {
		pStr = fmt.Sprintf("%s %d %d", pStr, step.Brightness, ledDelay(step.Duration))
	}
	if _, err := patternFile.WriteString(pStr); err != nil {
		return err
	}

	return nil
}

func (*Pattern) Name() string {
	return "pattern"
}

func (*Pattern) Cleanup() {}
