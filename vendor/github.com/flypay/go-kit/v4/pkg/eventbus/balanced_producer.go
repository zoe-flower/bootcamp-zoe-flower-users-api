package eventbus

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"google.golang.org/protobuf/proto"
)

const minRequiredEmitters = 2

// BalancedEmitter is an Emitter used in weighted balancing
type BalancedEmitter struct {
	Emitter     Emitter
	Percentage  int
	accumulated int
}

// NewBalancedEmitter returns a BalancedEmitter to be used in
// weighted routing when emitting events
func NewBalancedEmitter(e Emitter, p int) BalancedEmitter {
	return BalancedEmitter{
		Emitter:     e,
		Percentage:  p,
		accumulated: p,
	}
}

// BalancedProducer allows emitting events across a balanced producer
type BalancedProducer struct {
	emitters []BalancedEmitter
	mu       sync.Mutex
}

// NewBalancedProducer will balance the load across emitters.
// The total percentage of BalancedEmitter.Percentage must be equal to
// 100.
func NewBalancedProducer(emitters ...BalancedEmitter) (*BalancedProducer, error) {
	if len(emitters) < minRequiredEmitters {
		return nil, fmt.Errorf("two or more producers are required")
	}

	accumulatedTotal := 0
	for _, e := range emitters {
		accumulatedTotal += e.Percentage
	}

	if accumulatedTotal != 100 {
		return nil, fmt.Errorf("accumulated percentage of the emitters must be 100. Total given: %d", accumulatedTotal)
	}

	// Make a copy of emitters to work with internally so we're not messing with anything outside
	sortedEmitters := append([]BalancedEmitter(nil), emitters...)

	sort.Slice(sortedEmitters, func(i, j int) bool { return sortedEmitters[i].accumulated > sortedEmitters[j].accumulated })
	return &BalancedProducer{emitters: sortedEmitters}, nil
}

func (bp *BalancedProducer) rebalance() BalancedEmitter {
	/*
		choose the host with the largest percentageAccum
		subtract 100 from the percentageAccum for the chosen host
		add percentageLoad to percentageAccum for all hosts, including the chosen host
	*/
	bp.mu.Lock()
	defer bp.mu.Unlock()

	bp.emitters[0].accumulated -= 100

	for i := 0; i < len(bp.emitters); i++ {
		bp.emitters[i].accumulated += bp.emitters[i].Percentage
	}

	current := bp.emitters[0]
	sort.Slice(bp.emitters, func(i, j int) bool { return bp.emitters[i].accumulated > bp.emitters[j].accumulated })

	return current
}

// Emit will emit across the producers properly
func (bp *BalancedProducer) Emit(ctx context.Context, event proto.Message, headers ...Header) error {
	current := bp.rebalance()

	return current.Emitter.Emit(ctx, event, headers...)
}

// Use currently does nothing
// implements eventbus.producer
func (*BalancedProducer) Use(ProducerMiddleware) {}

// Close currently does nothing
// implements eventbus.producer
func (*BalancedProducer) Close() error { return nil }
