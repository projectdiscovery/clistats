package clistats

// DefaultOptions for clistats
var DefaultOptions = Options{
	ListenAddress: "127.0.0.1:63636",
	Web:           true,
}

// Options to customize behavior
type Options struct {
	ListenAddress string
	Web           bool
}
