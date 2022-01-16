package main

import (
	"fmt"
	"log"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/lipence/eventbus"
)

var GlobalBus = eventbus.Sync()

const EvtNameOnStringify = "nameOnStringify"

// define an Event

type eventNameOnStringify struct{ name string }

func (e eventNameOnStringify) Key() eventbus.Key {
	return EvtNameOnStringify
}

// define an Subscriber

type eventNameUppercase struct {
	nameCaser cases.Caser
}

func (eventNameUppercase) Events() map[eventbus.Key]eventbus.Priority {
	return map[eventbus.Key]eventbus.Priority{
		EvtNameOnStringify: 0,
	}
}

func (e *eventNameUppercase) OnEvent(evt eventbus.Event) error {
	if evtNameStringify, ok := evt.(*eventNameOnStringify); ok {
		evtNameStringify.name = e.nameCaser.String(evtNameStringify.name)
	}
	return nil
}

func init() {
	if err := GlobalBus.Register(&eventNameUppercase{
		nameCaser: cases.Title(language.English),
	}); err != nil {
		log.Fatal(err)
	}
}

// trigger Event

type Name struct {
	first, last string
}

func (n Name) String() string {
	var name = fmt.Sprintf("%s %s", n.first, n.last)
	var event = eventNameOnStringify{name: name}
	if err := GlobalBus.Post(&event); err == nil {
		return event.name
	}
	return name
}

func main() {
	var peter = &Name{"peter", "lee"}
	log.Println(peter) // log displays `Peter Lee`
}
