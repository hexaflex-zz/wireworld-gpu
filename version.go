package main

import (
	"fmt"
)

// Application name and version constants.
const (
	AppVendor  = "hexaflex"
	AppName    = "wireworld-gpu"
	AppVersion = "v0.0.2"
)

// Version returns a string with version information.
func Version() string {
	return fmt.Sprintf("%s %s %s",
		AppVendor, AppName, AppVersion)
}
