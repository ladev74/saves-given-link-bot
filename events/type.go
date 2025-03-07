package events

type Fetcher interface {
	Fetch(linit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknow Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
