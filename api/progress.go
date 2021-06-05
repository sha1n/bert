package api

// Listener a listener for benchmark progress events
type Listener interface {
	OnBenchmarkStart()
	OnBenchmarkEnd()
	OnScenarioStart(id ID)
	OnScenarioEnd(id ID)
	OnMessagef(id ID, format string, args ...interface{})
	OnMessage(id ID, message string)
	OnError(id ID, err error)
}

