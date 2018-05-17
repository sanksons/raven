package raven

import "time"

var _ RavenManager = (*RedisSimple)(nil)
var _ RavenManager = (*RedisCluster)(nil)

//Time to wait incase Q is empty.
const BLOCK_FOR_DURATION = 10 * time.Second

//No. of times to try incase of failure.
const MAX_TRY_LIMIT = 3

//
// An interface to be implemented by all Raven Managers.
// Raven Managers contains the exact logic and implementation details of message delivery
// and retrieval.
//
type RavenManager interface {

	//get Sequence to be allocated to message.
	//GetMsgSeq(mtype string, destination Destination) (int, error)

	// Message to be sent, Destination name
	Send(message Message, destination Destination) error

	// Source from which message is to be received.
	// Q in which message is to be stored for temporary basis.
	Receive(source Source, processingQ Q) (*Message, error)

	// Reliable version of above method.
	// Source from which message is to be received.
	// Q in which message is to be stored for temporary basis.
	//ReceiveReliable(source Source, processingQ Q) (*Message, error)

	// Mark the supplied message as processed.
	MarkProcessed(message *Message, processingQ Q) error

	// Mark the supplied message as failed.
	MarkFailed(message *Message, deadQ Q, processingQ Q) error
}
