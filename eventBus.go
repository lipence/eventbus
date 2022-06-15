package eventbus

type Key string

func (k Key) String() string { return string(k) }

type Priority uint8

type Initializer interface {
	Init() error
}

type Event interface {
	Key() Key
}

type Subscriber interface {
	Events() map[Key]Priority
	OnEvent(Event) error
}

type Bus interface {
	Post(Event) error
	PostAndClear(Event) error
	Register(Subscriber) error
	Unregister(Subscriber)
}
