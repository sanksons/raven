package raven

type FarmManager interface {
	Send(string) error
}
