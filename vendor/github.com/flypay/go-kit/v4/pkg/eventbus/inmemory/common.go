package inmemory

const (
	defaultBufferSize = 1024
)

func NewMessageChannel() chan *Message {
	return make(chan *Message, defaultBufferSize)
}

func deadletterChannel() chan *Message {
	return make(chan *Message, defaultBufferSize)
}
