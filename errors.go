package raven

import (
	"errors"
)

var ErrNoMessage error = errors.New("A Message Needs to be defined")
var ErrNoDestination error = errors.New("A Destination Needs to be defined")
var ErrInvalidDestination error = errors.New("Invalid Destination defined")

var ErrEmptyQueue error = errors.New("Empty Queue")

var ErrNotImplemented error = errors.New("Feature Not Implemented")

var ErrTmpFailure error = errors.New("Temporary Failure")
