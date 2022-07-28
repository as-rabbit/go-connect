package mysql

type Logger interface {
	Info()
	
	Error()
	
	Warn()
	
	Trace()
}

type Log interface {
	Info()
	
	Error()
}

type DefaultLogger struct {
}

func (d *DefaultLogger) Info() {
	panic("implement me")
}

func (d *DefaultLogger) Error() {
	panic("implement me")
}

func (d *DefaultLogger) Warn() {
	panic("implement me")
}

func (d *DefaultLogger) Trace() {
	panic("implement me")
}
