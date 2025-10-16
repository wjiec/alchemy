package alchemy

import (
	"context"
)

// Addr represents a network end point address.
type Addr interface {
	// Network returns the name of the network (for example, "tcp", "udp").
	Network(context.Context) string

	// String returns the string format of the address (for example, "192.0.2.1:25", "[2001:db8::1]:80").
	String(context.Context) string
}

// NewAddr creates and returns a new Addr implementation for the given network and address.
func NewAddr(network, address string) Addr {
	return &networkAddress{network: network, address: address}
}

// networkAddress is an implementation of net.Addr that holds a network name and address string.
type networkAddress struct {
	network, address string
}

// Network returns the network name of the address.
func (a *networkAddress) Network(_ context.Context) string { return a.network }

// String returns the address as a string.
func (a *networkAddress) String(_ context.Context) string { return a.address }

// TCP creates a net.Addr implementation with the network set to "tcp" and the specified address.
func TCP(address string) Addr { return &networkAddress{network: "tcp", address: address} }
