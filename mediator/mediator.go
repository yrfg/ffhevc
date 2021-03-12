package mediator

// M is a universal public variable of Mediator
var M *Mediator = &Mediator{}

// Mediator is a toolset of media processing
type Mediator struct {
}

// VideoTranscode return a VideoTranscoder
func (m *Mediator) VideoTranscode() *VideoTranscoder {
	return NewVideoTransCoder()
}
