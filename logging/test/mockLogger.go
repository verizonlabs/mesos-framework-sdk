package test

type MockLogger struct{}

func (m MockLogger) Emit(severity uint8, template string, args ...interface{}) {

}
