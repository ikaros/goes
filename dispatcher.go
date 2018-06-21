package goes

import (
	"strconv"
)

type Reactor = func(Event) error

var reactorRegistry = map[string]registryReactors{}

type registryReactors struct {
	Sync  []Reactor
	Async []Reactor
}

func On(event EventInterface, sync []Reactor, async []Reactor) {
	eventType := event.AggregateType() +
		"." + event.Action() +
		"." + strconv.FormatUint(event.Version(), 10)

	if sync == nil {
		sync = []Reactor{}
	}
	if async == nil {
		async = []Reactor{}
	}

	var newSync []Reactor
	var newAsync []Reactor

	if reactors, ok := reactorRegistry[eventType]; ok == true {
		newSync = reactors.Sync
		newAsync = reactors.Async
	} else {
		newSync = []Reactor{}
		newAsync = []Reactor{}
	}

	newSync = append(newSync, sync...)
	newAsync = append(newAsync, async...)

	reactorRegistry[eventType] = registryReactors{Sync: newSync, Async: newAsync}
}

func Dispatch(event Event) error {
	data := event.Data.(EventInterface)
	eventType := data.AggregateType() +
		"." + data.Action() +
		"." + strconv.FormatUint(data.Version(), 10)

	if reactors, ok := reactorRegistry[eventType]; ok == true {
		// dispatch sync reactor synchronously
		// it can be something like a projection
		for _, syncReactor := range reactors.Sync {
			if err := syncReactor(event); err != nil {
				return nil
			}
		}

		// dispatch async reactors asynchronously
		for _, asyncReactor := range reactors.Async {
			go asyncReactor(event)
		}
	}
	return nil
}
