package core

import (
	"fmt"

	"github.com/jimbot9k/norman/internal/core/dbobjects"
)

// Adapter defines the interface that all database adapters must implement.
// It provides methods for connection management and schema introspection.
type Adapter interface {
	// Name returns the human-readable name of the adapter.
	Name() string
	// Version returns the adapter version string.
	Version() string
	// UniqueSignature returns a unique identifier for this adapter type.
	UniqueSignature() string
	// IsConnectionStringCompatible reports whether the adapter can handle the given connection string.
	IsConnectionStringCompatible(connString string) bool
	// Connect establishes a connection using the provided connection string.
	Connect(connString string) error
	// Close terminates the active database connection.
	Close() error
	// IsConnected reports whether the adapter has an active connection.
	IsConnected() bool
	// MapDatabase maps the database structure into dbobjects.Database.
	MapDatabase() (*dbobjects.Database, []error)
}

// AdapterManager manages database adapter registration and connections.
// It maintains a registry of available adapters and tracks the currently active connection.
type AdapterManager struct {
	adapters      map[string]Adapter
	activeAdapter Adapter
}

// NewAdapterManager creates a new AdapterManager and registers the provided adapters.
func NewAdapterManager(adapters []Adapter) *AdapterManager {
	n := &AdapterManager{
		adapters: make(map[string]Adapter),
	}
	for _, a := range adapters {
		n.registerAdapter(a)
	}
	return n
}

// RegisterAdapter adds an adapter to the manager's registry.
// The adapter is keyed by its UniqueSignature.
func (n *AdapterManager) registerAdapter(a Adapter) {
	if n.adapters == nil {
		n.adapters = make(map[string]Adapter)
	}
	n.adapters[a.UniqueSignature()] = a
}

// findCompatibleAdapter returns the first registered adapter that can handle
// the given connection string, or nil if no compatible adapter is found.
func (n *AdapterManager) findCompatibleAdapter(connString string) Adapter {
	for _, a := range n.adapters {
		if a.IsConnectionStringCompatible(connString) {
			return a
		}
	}
	return nil
}

// Connect establishes a database connection using a compatible adapter.
// It returns an error if no compatible adapter is found or if the connection fails.
func (n *AdapterManager) Connect(connString string) (Adapter, error) {

	if n.activeAdapter != nil && n.activeAdapter.IsConnected() {
		return nil, fmt.Errorf("an adapter is already connected")
	}

	adapter := n.findCompatibleAdapter(connString)
	if adapter == nil {
		return nil, fmt.Errorf("no compatible adapter found for connection string")
	}

	err := adapter.Connect(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect using adapter %s: %v", adapter.UniqueSignature(), err)
	}
	n.activeAdapter = adapter
	return n.GetActiveAdapter(), nil
}

// Close terminates the active database connection, if any.
func (n *AdapterManager) Close() error {
	if n.activeAdapter == nil || !n.activeAdapter.IsConnected() {
		return fmt.Errorf("no active connection to close")
	}

	err := n.activeAdapter.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	n.activeAdapter = nil
	return nil
}

// GetActiveAdapter returns the currently connected adapter, or nil if no connection is active.
func (n *AdapterManager) GetActiveAdapter() Adapter {
	if n.activeAdapter == nil {
		return nil
	}

	return n.activeAdapter
}
