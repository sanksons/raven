package raven

var _ RavenManager = (*RedisCluster)(nil)

type RavenManager interface {
	// Message to be sent, Destination name
	Send(message Message, destination Destination) error

	// Source from which message is to be received.
	// Q in which message is to be stored for temporary basis.
	Receive(source Source, processingQ Q) (*Message, error)

	// Mark the supplied message as processed.
	MarkProcessed(message *Message, q Q) error

	// Mark the supplied message as processed.
	MarkFailed(message *Message, deadQ Q, processingQ Q) error
}
