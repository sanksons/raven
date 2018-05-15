package raven

var _ RavenManager = (*RedisCluster)(nil)

type RavenManager interface {
	// Message to be sent, Destination name
	Send(message string, destination string) error

	Receive(dest string, tempQ string) (string, error)
}
