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

type DefaultExecutorEvents struct {
}

func (d *DefaultExecutorEvents) Subscribed() {

}

func (d *DefaultExecutorEvents) Launch() {

}

func (d *DefaultExecutorEvents) LaunchGroup() {

}
func (d *DefaultExecutorEvents) Kill() {

}
func (d *DefaultExecutorEvents) Acknowledged() {

}
func (d *DefaultExecutorEvents) Message() {

}
func (d *DefaultExecutorEvents) Shutdown() {

}
func (d *DefaultExecutorEvents) Error() {

}
