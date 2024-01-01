package trigger

type Trigger interface {
	Setup(root string) error
	Name() string
	Cleanup()
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
