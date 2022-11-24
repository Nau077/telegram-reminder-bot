package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

// список событий, который будем ипользовать

const (
	// неизвестный тип, чтобы понять те случаи, когда не смогли понять тип события
	Unknown Type = iota
)

type Event struct {
	Type Type
	Text string
}
