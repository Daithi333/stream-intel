package pipeline

import (
	"context"

	"github.com/dmcelhill/stream-intel/pkg/model"
)

type Pipeline struct {
	chEvents chan model.TaxiTrip
}

func New(bufferSize int) *Pipeline {
	return &Pipeline{
		chEvents: make(chan model.TaxiTrip, bufferSize),
	}
}

func (p *Pipeline) Send(ctx context.Context, event model.TaxiTrip) error {
	select {
	case p.chEvents <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Pipeline) Receive() <-chan model.TaxiTrip {
	return p.chEvents
}

func (p *Pipeline) Len() int {
	return len(p.chEvents)
}

func (p *Pipeline) Close() {
	close(p.chEvents)
}
