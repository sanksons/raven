package raven

import (
	"errors"
)

var ErrNoMessage error = errors.New("A Message Needs to be defined")
var ErrNoDestination error = errors.New("A Destination Needs to be defined")
var ErrInvalidDestination error = errors.New("Invalid Destination defined")
