package clistats

import "sync/atomic"

// counterStatistic is a counter stats field
type counterStatistic struct {
	value       uint64
	description string
}

// staticStatistic is a static stats field
type staticStatistic struct {
	value       interface{}
	description string
}

// dynamicStatistic is a dynamic stats field
type dynamicStatistic struct {
	value       DynamicCallback
	description string
}

// AddCounter adds a uint64 counter field to the statistics client.
//
// A counter is used to track an increasing quantity, like requests,
// errors etc.
func (s *Statistics) AddCounter(id, description string) error {
	if s.hasStarted() {
		return ErrEventLoopStarted
	}
	s.counters[id] = &counterStatistic{
		value:       0,
		description: description,
	}
	return nil
}

// GetCounter returns the current value of a counter.
func (s *Statistics) GetCounter(id string) (uint64, bool) {
	counter, ok := s.counters[id]
	if !ok {
		return 0, false
	}
	return atomic.LoadUint64(&counter.value), true
}

// IncrementCounter increments the value of a counter by a count.
func (s *Statistics) IncrementCounter(id string, count int) {
	counter, ok := s.counters[id]
	if !ok {
		return
	}
	atomic.AddUint64(&counter.value, uint64(count))
}

// AddStatic adds a static information field to the statistics.
//
// The value for these metrics will remain constant throughout the
// lifecycle of the statistics client. All the values will be
// converted into string and displayed as such.
func (s *Statistics) AddStatic(id, description string, value interface{}) error {
	if s.hasStarted() {
		return ErrEventLoopStarted
	}
	s.static[id] = &staticStatistic{
		value:       value,
		description: description,
	}
	return nil
}

// GetStatic returns the original value for a static field.
func (s *Statistics) GetStatic(id string) (interface{}, bool) {
	static, ok := s.static[id]
	if !ok {
		return nil, false
	}
	return static.value, true
}

// AddDynamic adds a dynamic field to display whose value
// is retrieved by running a callback function.
//
// The callback function performs some actions and returns the value
// to display. Generally this is used for calculating requests per
// seconds, elapsed time, etc.
func (s *Statistics) AddDynamic(id, description string, Callback DynamicCallback) error {
	if s.hasStarted() {
		return ErrEventLoopStarted
	}
	s.dynamic[id] = &dynamicStatistic{
		value:       Callback,
		description: description,
	}
	return nil
}

// GetDynamic returns the dynamic field callback for data retrieval.
func (s *Statistics) GetDynamic(id string) (DynamicCallback, bool) {
	dynamic, ok := s.dynamic[id]
	if !ok {
		return nil, false
	}
	return dynamic.value, true
}
