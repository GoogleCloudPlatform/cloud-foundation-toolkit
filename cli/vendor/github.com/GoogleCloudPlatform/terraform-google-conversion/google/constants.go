package google

import (
	"errors"

	"github.com/hashicorp/terraform/helper/mutexkv"
)

// ErrNoConversion can be returned if a conversion is unable to be returned
// because of the current state of the system.
// Example: The conversion requires that the resource has already been created
// and is now being updated).
var ErrNoConversion = errors.New("no conversion")

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()
