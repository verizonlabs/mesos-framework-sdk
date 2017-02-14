package executor

/*
The executor interface should only be made concrete by a custom executor.
*/

// All calls a custom executor should make.
type Executor interface {
	Subscribe()
	Update()
	Message()
}
