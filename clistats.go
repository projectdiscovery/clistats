package clistats

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/pkg/errors"
)

// StatisticsClient is an interface implemented by a statistics client.
//
// A unique ID is to be provided along with a description for the field to be
// displayed as output.
//
// Multiple types of statistics are provided like Counters as well as static
// fields which display static information only.
//
// A metric cannot be added once the client has been started. An
// error will be returned if the metric cannot be added. Already existing fields
// of same names are overwritten.
type StatisticsClient interface {
	// Start starts the event loop of the stats client.
	Start(printer PrintCallback, tickDuration time.Duration) error

	// Stop stops the event loop of the stats client
	Stop() error

	// AddCounter adds a uint64 counter field to the statistics client.
	//
	// A counter is used to track an increasing quantity, like requests,
	// errors etc.
	AddCounter(id, description string) error

	// GetCounter returns the current value of a counter.
	GetCounter(id string) (uint64, bool)

	// IncrementCounter increments the value of a counter by a count.
	IncrementCounter(id string, count int)

	// AddStatic adds a static information field to the statistics.
	//
	// The value for these metrics will remain constant throughout the
	// lifecycle of the statistics client. All the values will be
	// converted into string and displayed as such.
	AddStatic(id, description string, value interface{}) error

	// GetStatic returns the original value for a static field.
	GetStatic(id string) (interface{}, bool)

	// AddDynamic adds a dynamic field to display whose value
	// is retrieved by running a callback function.
	//
	// The callback function performs some actions and returns the value
	// to display. Generally this is used for calculating requests per
	// seconds, elapsed time, etc.
	AddDynamic(id, description string, Callback DynamicCallback) error

	// GetDynamic returns the dynamic field callback for data retrieval.
	GetDynamic(id string) (DynamicCallback, bool)
}

// DynamicCallback is called during statistics calculation for a dynamic
// field.
//
// The value returned from this callback is displayed as the current value
// of a dynamic field. This can be utilised to calculated things like elapsed
// time, requests per seconds, etc.
type DynamicCallback func(client StatisticsClient) interface{}

// Statistics is a client for showing statistics on the stdout.
type Statistics struct {
	ctx    context.Context
	cancel context.CancelFunc
	ticker tickerInterface
	events <-chan keyboard.KeyEvent

	// started indicates if the client has started.
	started uint32

	// counters is a list of counters for the client. These can only
	// be accessed concurrently via atomic operations and once the main
	// event loop has started must not be modified.
	counters map[string]*counterStatistic

	// static contains a list of static counters for the client.
	static map[string]*staticStatistic

	// dynamic contains a lsit of dynamic metrics for the client.
	dynamic map[string]*dynamicStatistic

	// printer is the printing callback for data display
	printer PrintCallback
}

var (
	// ErrEventLoopStarted is returned when stats event loop has already started
	ErrEventLoopStarted = errors.New("stats event loop started")
)

func (s *Statistics) hasStarted() bool {
	return atomic.LoadUint32(&s.started) == 1
}

// PrintCallback is used by clients to build and display a string on the screen.
type PrintCallback func(client StatisticsClient)

var _ StatisticsClient = (*Statistics)(nil)

// New creates a new statistics client for cli stats printing.
func New() *Statistics {
	ctx, cancel := context.WithCancel(context.Background())

	return &Statistics{
		ctx:      ctx,
		cancel:   cancel,
		started:  0,
		counters: make(map[string]*counterStatistic),
		static:   make(map[string]*staticStatistic),
		dynamic:  make(map[string]*dynamicStatistic),
	}
}

// Start starts the event loop of the stats client.
func (s *Statistics) Start(printer PrintCallback, tickDuration time.Duration) error {
	atomic.StoreUint32(&s.started, 1)
	s.printer = printer

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		return errors.Wrap(err, "could not get keyboard events")
	}
	s.events = keysEvents

	go s.eventLoop(tickDuration)
	return nil
}

// eventLoop is the event loop listening for keyboard events as well as
// looking out for cancellation attempts.
func (s *Statistics) eventLoop(tickDuration time.Duration) {
	if tickDuration != -1 {
		s.ticker = &ticker{t: time.NewTicker(tickDuration)}
	} else {
		s.ticker = &noopTicker{tick: make(chan time.Time)}
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.Tick():
			s.printer(s)
		case <-s.events:
			s.printer(s)
		}
	}
}

// Stop stops the event loop of the stats client
func (s *Statistics) Stop() error {
	s.cancel()
	keyboard.Close()
	s.ticker.Stop()
	return nil
}
