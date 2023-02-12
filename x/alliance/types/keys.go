package types

const (
	// ModuleName defines the module name
	ModuleName = "alliance"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_alliance"
)

// IBC events
const (
	AttributeKeyAck = "acknowledgement"
)

// Events
const (
	EventAlliance = "alliance_event"
)

var (
	ICAKey = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
