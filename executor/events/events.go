package events

// Interface for all events sent to a custom executor.
type ExecutorEvents interface {
	Subscribed()
	Launch()
	LaunchGroup()
	Kill()
	Acknowledged()
	Message()
	Shutdown()
	Error()
}
