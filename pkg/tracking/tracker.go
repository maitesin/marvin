package tracking

type Event struct {
	Timestamp   string
	Information string
}

type Tracker interface {
	Track(id string) ([]Event, error)
}
