package errchan

type ErrorProducer interface {
	Errors() <-chan error
}

type BufferedErrorChannel struct {
	errors chan error
}

func New(b int) *BufferedErrorChannel {
	return &BufferedErrorChannel{
		errors: make(chan error, b),
	}
}

func (e *BufferedErrorChannel) Errors() <-chan error {
	return e.errors
}
func (e *BufferedErrorChannel) SendError(err error) {
	select {
	case e.errors <- err:
	default:
	}
}

func (e *BufferedErrorChannel) ConsumeErrors(ep ErrorProducer) {
	go func() {
		for err := range ep.Errors() {
			e.SendError(err)
		}
	}()
}
