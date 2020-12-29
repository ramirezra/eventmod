package eventmod

import "context"

// Version for EventMod
const Version = "0.0.1"

// Record contains event data and metadata.
type Record struct {
	AggregateType        string
	AggregateID          string
	GlobalSequenceNumber uint64
	StreamSequenceNumber uint64
	InsertTimeStamp      uint64
	EventID              string
	EventType            string
	Data                 interface{}
	Metadata             interface{}
}

// Store is the interface that stores and retrieves event records.
type Store interface {
	EventBinder
	EventStartingWith(ctx context.Context, globalSequenceNumber uint64) <-chan *Record
	EventsByAggregateTypeStartingWith(ctx context.Context, globalSequenceNumber uint64, aggregateTypes ...string) <-chan *Record
	EventsByStreamStartingWith(ctx context.Context, streamSequenceNumber uint64, streamName string) <-chan *Record
	OptimisticSave(expectedStreamSequenceNumber uint64, eventRecords ...*EventRecord) error
	Save(eventRecords ...*EventRecord) error
	Subscribe(subscribers ...RecordSubscriber)
	SubscribeStartingWith(ctx context.Context, globalSequenceNumber uint64, subscribers ...RecordSubscriber)
	TotalEventsInStream(streamName string) uint64
}

// Event is the interface that defines the required event methods.
type Event interface {
	AggregateMessage
	EventType() string
}

// AggregateMessage is the interfact that support building an event stream name.
type AggregateMessage interface {
	AggregateID() string
	AggregateType() string
}

// EventRecord stores the event and metadata to be used for persisting data.
type EventRecord struct {
	Event    Event
	Metadata interface{}
}

// EventBinder defines how to bind events for serialization.
type EventBinder interface {
	Bind(events ...Event)
}

// RecordSubscriber is the interface that defines how a projection receives Records.
type RecordSubscriber interface {
	Accept(record *Record)
}

// RecordSubscriberFunc is an adapter type that allows the use of ordinary functions as record subscribers.
// If f is a function with the appropriate signature, RecordSubsriberFunc(f) is a Hanlder that calls f.
type RecordSubscriberFunc func(*Record)

// Accept receives a record.
func (f RecordSubscriberFunc) Accept(record *Record) {
	f(record)
}
