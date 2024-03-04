package tracing

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
)

const (
	TeamKey      = "team"
	FlowKey      = "flow"
	TierKey      = "tier"
	RequestIDKey = xray.RequestIDKey
)

type Annotation struct {
	Key   string
	Value interface{}
}

type Metadata struct {
	Key   string
	Value interface{}
}

// AddAnnotation adds an annotation to the span in the underlying context.
func AddAnnotation(ctx context.Context, key string, value interface{}) {
	_ = xray.AddAnnotation(ctx, key, value)
}

// AddAnnotations adds multiple annotations to the span in the underlying context.
func AddAnnotations(ctx context.Context, annotations ...Annotation) {
	for _, a := range annotations {
		AddAnnotation(ctx, a.Key, a.Value)
	}
}

// AddMetadata adds metadata to the span in the underlying context.
func AddMetadata(ctx context.Context, key string, value interface{}) {
	_ = xray.AddMetadata(ctx, key, value)
}

// AddMetadatas adds multiple metadata to the span in the underlying context.
func AddMetadatas(ctx context.Context, metadatas ...Metadata) {
	for _, m := range metadatas {
		AddMetadata(ctx, m.Key, m.Value)
	}
}

// Span is a representation of a component being traced.
// Currently this wraps an xray segment.
type Span struct {
	xraySegment *xray.Segment
}

// StartSpan creates a new span that can be used to trace some
// component of a system.
// NOTE: the returned span needs to be closed: span.Close(err)
func StartSpan(ctx context.Context, name string) (context.Context, *Span) {
	newCtx, span := xrayNewSpan(ctx, name)
	return newCtx, span
}

// AddAnnotation adds an annotation to the span.
// With xray this will create key values that are indexed
// and can be searched on.
func (s *Span) AddAnnotation(key string, value interface{}) {
	_ = s.xraySegment.AddAnnotation(key, value)
}

// AddAnnotations adds multiple annotations to the span.
// With xray this will create key values that are indexed
// and can be searched on.
func (s *Span) AddAnnotations(annotations ...Annotation) {
	for _, a := range annotations {
		s.AddAnnotation(a.Key, a.Value)
	}
}

// AddMetadata adds metadata to the span.
func (s *Span) AddMetadata(key string, value interface{}) {
	_ = s.xraySegment.AddMetadata(key, value)
}

// AddMetadatas adds multiple metadata to the span.
func (s *Span) AddMetadatas(metadatas ...Metadata) {
	for _, m := range metadatas {
		s.AddMetadata(m.Key, m.Value)
	}
}

// TraceID returns the trace id of the underlying xray segment.
func (s *Span) TraceID() string {
	return s.xraySegment.DownstreamHeader().String()
}

// Close takes and error and will close off the current span
// adding the error to the span if one is passed in.
func (s *Span) Close(err error) {
	s.xraySegment.Close(err)
}

// xrayNewSpan creates a new X-Ray segment or subsegment. The segment
// is derived from the passed in context. If the segment exsists then
// a subsegment is created from that segment, otherwise if a segment
// doesn't already exist a new one is created. The segment or subsegment
// will be given the passed in name and the metadata will be added
// to the segment or subsegment.
//
// The new segment or subsegment will be attached to the context. Any
// subsequent tracing of that context will created a new subsegment using
// the segment on the context as the parent.
//
// The trace is sent to the X-Ray daemon when the segment is closed.
func xrayNewSpan(ctx context.Context, name string) (context.Context, *Span) {
	var newCtx context.Context
	var newSeg *xray.Segment

	// If no parent segment exists, create a new one.
	if seg := xray.GetSegment(ctx); seg == nil {
		newCtx, newSeg = xray.BeginSegment(ctx, name)
	} else {
		newCtx, newSeg = xray.BeginSubsegment(ctx, name)
	}
	return newCtx, &Span{xraySegment: newSeg}
}
