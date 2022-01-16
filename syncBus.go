package eventbus

import "sync"

var defaultSyncBus Bus
var defaultSyncBusInit sync.Once

func Sync() Bus {
	defaultSyncBusInit.Do(func() {
		defaultSyncBus = NewSync()
	})
	return defaultSyncBus
}

type SyncBus struct {
	agency
}

func (b *SyncBus) Post(e Event) (err error) {
	var queue []item
	if evtReg, ok := b.agency[e.Key()]; ok {
		queue = evtReg.Queue()
	} else {
		return nil
	}
	if initializer, ok := e.(Initializer); ok {
		if err = initializer.Init(); err != nil {
			return err
		}
	}
	for _, t := range queue {
		if err = t.Subscriber.OnEvent(e); err != nil {
			return err
		}
	}
	return nil
}

func (b *SyncBus) PostAndClear(e Event) (err error) {
	defer func() {
		if err == nil {
			b.UnregisterEvent(e)
		}
	}()
	return b.Post(e)
}

func NewSync() *SyncBus {
	return &SyncBus{
		agency: agency{},
	}
}
