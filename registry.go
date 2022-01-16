package eventbus

import (
	"sort"
)

type item struct {
	Priority   Priority
	Subscriber Subscriber
}

type registry struct {
	sorted bool
	items  []item
}

func (a *registry) Add(p Priority, s Subscriber) {
	a.items = append(a.items, item{
		Priority:   p,
		Subscriber: s,
	})
	a.sorted = false
}

func (a *registry) Remove(p Priority, s Subscriber) {
	var newItems = make([]item, 0, len(a.items))
	for _, t := range a.items {
		if !(t.Subscriber == s && t.Priority == p) {
			newItems = append(newItems, t)
		}
	}
	if len(newItems) != len(a.items) {
		a.items = newItems
		a.sorted = false
	}
}

func (a *registry) Queue() []item {
	if a.sorted {
		return a.items
	}
	sort.Slice(a.items, func(i, j int) bool {
		return a.items[i].Priority < a.items[j].Priority
	})
	a.sorted = true
	return a.items
}

type agency map[Key]*registry

func (a agency) Register(subscriber Subscriber) (err error) {
	if a == nil || subscriber == nil {
		return
	}
	if initializer, ok := subscriber.(Initializer); ok {
		if err = initializer.Init(); err != nil {
			return err
		}
	}
	for key, priority := range subscriber.Events() {
		if _, ok := a[key]; !ok {
			a[key] = &registry{}
		}
		a[key].Add(priority, subscriber)
	}
	return nil
}

func (a agency) Unregister(subscriber Subscriber) {
	if a == nil || subscriber == nil {
		return
	}
	for key, priority := range subscriber.Events() {
		if _, ok := a[key]; !ok {
			continue
		}
		a[key].Remove(priority, subscriber)
	}
	return
}

func (a agency) UnregisterEvent(e Event) {
	if a == nil || e == nil {
		return
	}
	delete(a, e.Key())
}
