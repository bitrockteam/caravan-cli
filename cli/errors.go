package cli

import (
	"fmt"
)

type ConfigFileNotFound struct {
	Err error
}

func (e ConfigFileNotFound) Error() string {
	return fmt.Sprintf("config file not found: %s", e.Err.Error())
}

type ConfigFileUnreadable struct {
	Err error
}

func (e ConfigFileUnreadable) Error() string {
	return fmt.Sprintf("config file unreadable: %s", e.Err.Error())
}
