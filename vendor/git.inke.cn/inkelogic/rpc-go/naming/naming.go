// Package naming defines the naming API and related data structures for gRPC.
// The interface is EXPERIMENTAL and may be suject to change.
package naming

type WatchType uint8

const (
	SERVICE WatchType = iota
	KV      WatchType = iota
)

// Operation defines the corresponding operations for a name resolution change.
type Operation uint8

const (
	// Add indicates a new address is added.
	Add Operation = iota
	// Delete indicates an exisiting address is deleted.
	Delete
)

// Update defines a name resolution update. Notice that it is not valid having both
// empty string Addr and nil Metadata in an Update.
type Update struct {
	// Op indicates the operation of the update.
	Op Operation
	// Addr is the updated address. It is empty string if there is no address update.
	Addr string
	// Metadata is the updated metadata. It is nil if there is no metadata update.
	// Metadata is not required for a custom naming implementation.
	Metadata interface{}
}

// Resolver creates a Watcher for a target to track its resolution changes.
type Resolver interface {
	// Resolve creates a Watcher for target.
	Resolve(target string) (Watcher, error)
}

// Watcher watches for the updates on the specified target.
type Watcher interface {
	// Next blocks until an update or error happens. It may return one or more
	// updates. The first call should get the full set of the results. It should
	// return an error if and only if Watcher cannot recover.
	Next() ([]*Update, error)
	// Close closes the Watcher.
	Close()
}

type KvMessage struct {
	OriginPath string
	Path       string
	PathValue  string
}

type Message struct {
	Type WatchType

	Target string
	Proto  string
	Addrs  []string

	KValue string
	Key    string
}

type Manager interface {
	Start() error

	UnRegister() error
	RegisterService(targets []string, proto string, tag []string, address string, port int) error
	GetService(target, proto, tag, dc string) ([]string, error)
	Watch(target, proto, tag, dc string)

	WatchKV(path string)
	GetKvValue(path string) (string, error)

	GetValues(keys []string) (map[string]string, error)

	Next() (*Message, error)
}
