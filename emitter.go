package podio

import (
	"sync"
	event "github.com/wauio/event-emitter"
)

// singleton implementation
var Emitter EventEmitterPodioWrapper
var once sync.Once

func GetPodioEmitter() EventEmitterPodioWrapper {
	once.Do(func() {
		Emitter = &podioEventEmitter{
			emitter: event.New(),
		}
	})

	return Emitter
}

// encapsulate event_emitter.EventEmitter in a singleton
// https://travix.io/encapsulating-dependencies-in-go-b0fd74021f5a
type EventEmitterPodioWrapper interface {
	On(event string, listener interface{}) error
	Fire(event string, payload ...interface{}) ([][]interface{}, error)
	FireBackground(event string, payload ...interface{}) (chan []interface{}, error)
}

type podioEventEmitter struct {
	emitter *event.EventEmitter
}

func (e *podioEventEmitter) On(event string, listener interface{}) error {
	return e.emitter.On(event, listener)
}

func (e *podioEventEmitter) Fire(event string, payload ...interface{}) ([][]interface{}, error) {
	return e.emitter.Fire(event, payload...)
}

func (e *podioEventEmitter) FireBackground(event string, payload ...interface{}) (chan []interface{}, error) {
	return e.emitter.FireBackground(event, payload...)
}