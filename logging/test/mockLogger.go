package mockLogging

type mockLogger struct{}

func (m *mockLogger) Emit(severity uint8, template string, args ...interface{}) {

}