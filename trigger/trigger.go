package trigger

import "time"

type Trigger interface {
	Setup(root string) error
	Name() string
	Cleanup()
}

func ledDelay(d time.Duration) int64 {
	if d > 0 && d < time.Millisecond {
		d = time.Millisecond
	}

	return int64(d / time.Millisecond)
}

type simple struct{}

func (*simple) Setup(string) error {
	return nil
}

func (*simple) Cleanup() {}

type DefaultOn struct{ *simple }

func (*DefaultOn) Name() string {
	return "default-on"
}

type Heartbeat struct{ *simple }

func (*Heartbeat) Name() string {
	return "heartbeat"
}

type Panic struct{ *simple }

func (*Panic) Name() string {
	return "panic"
}
