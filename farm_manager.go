package raven

type RavenManager interface {
	Send(string) error
}
