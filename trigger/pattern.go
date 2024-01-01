package trigger

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	repeatTrue  = "-1"
	repeatFalse = "1"
)

type Pattern struct {
	pattern [][2]uint
	repeat  bool
}

func NewPattern(pattern [][2]uint, repeat bool) *Pattern {
	return &Pattern{
		pattern,
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
	for _, t := range p.pattern {
		pStr = fmt.Sprintf("%s %d %d", pStr, t[0], t[1])
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
