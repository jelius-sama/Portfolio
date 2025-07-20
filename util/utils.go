package util

import (
	"strconv"
)

// AddrOf takes a value and returns its address as a pointer.
func AddrOf[T any](literal T) *T { return &literal }

func IsValidPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false // not a number
	}
	return port >= 1 && port <= 65535
}
