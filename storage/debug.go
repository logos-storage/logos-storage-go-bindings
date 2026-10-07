package storage

/*
   #include "bridge.h"
   #include <stdlib.h>

   static int cGoStorageDebug(void* storageCtx, void* resp) {
       return storage_debug(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageLogLevel(void* storageCtx, char* logLevel, void* resp) {
       return storage_log_level(storageCtx, logLevel, (StorageCallback) callback, resp);
   }

   static int cGoStoragePeerDebug(void* storageCtx, char* peerId, void* resp) {
       return storage_peer_debug(storageCtx, peerId, (StorageCallback) callback, resp);
   }
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

type Node struct {
	NodeId  string  `json:"nodeId"`
	PeerId  string  `json:"peerId"`
	Record  string  `json:"record"`
	Address *string `json:"address"`
	Seen    bool    `json:"seen"`
}

type RoutingTable struct {
	LocalNode Node   `json:"localNode"`
	Nodes     []Node `json:"nodes"`
}

type NatInfo struct {
	// Reachability determined by AutoNAT: Unknown, NotReachable or Reachable
	Reachability string `json:"reachability"`

	// True when the node runs the discovery protocol in client mode,
	// meaning it is not reachable and does not participate to DHT queries.
	ClientMode bool `json:"clientMode"`

	// True when the relay service is started, which happens when the node is
	// not reachable. It does not mean a reservation was obtained.
	RelayRunning bool `json:"relayRunning"`

	// Port mapping protocol in use: none, upnp, pmp, pcp or direct.
	PortMapping string `json:"portMapping"`
}

type Connection struct {
	PeerId string `json:"peerId"`

	// False when the connection goes through a circuit relay
	Direct bool `json:"direct"`
}

type DebugInfo struct {
	// Peer ID
	ID string `json:"id"`

	// Peer info addresses
	// Specified with `ListenIp` and `ListenPort` in `Config`
	Addrs []string `json:"addrs"`

	Spr string `json:"spr"`

	// Addresses contained in the provider record, the one announced
	// for the content the node provides
	ProviderAddresses []string `json:"providerAddresses"`

	// UDP addresses announced in the DHT
	DiscoveryAddresses []string `json:"discoveryAddresses"`

	PeersTable  RoutingTable `json:"table"`
	Nat         NatInfo      `json:"nat"`
	Connections []Connection `json:"connections"`
}

type PeerRecord struct {
	PeerId    string   `json:"peerId"`
	SeqNo     int      `json:"seqNo"`
	Addresses []string `json:"addresses,omitempty"`
}

// Debug retrieves debugging information from the Logos Storage node.
func (node StorageNode) Debug() (DebugInfo, error) {
	var info DebugInfo

	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageDebug(node.ctx, bridge.resp) != C.RET_OK {
		return info, bridge.callError("cGoStorageDebug")
	}

	value, err := bridge.wait()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(value), &info)
	return info, err
}

// UpdateLogLevel updates Chronicles’ runtime logging configuration.
// You can pass a plain level: TRACE, DEBUG, INFO, NOTICE, WARN, ERROR, FATAL.
// The default level is TRACE.
// You can also use Chronicles topic directives. So for example if you want
// to update the general level to INFO but want to see TRACE logs for the libstorage
// topic, you can pass "INFO,libstorage:TRACE".
func (node StorageNode) UpdateLogLevel(logLevel string) error {
	bridge := newBridgeCtx()
	defer bridge.free()

	var cLogLevel = C.CString(string(logLevel))
	defer C.free(unsafe.Pointer(cLogLevel))

	if C.cGoStorageLogLevel(node.ctx, cLogLevel, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoStorageLogLevel")
	}

	_, err := bridge.wait()
	return err
}

// StoragePeerDebug retrieves the peer record for a given peer ID.
// This function is available only if the flag
// -d:storage_enable_api_debug_peers=true was set at build time.
func (node StorageNode) StoragePeerDebug(peerId string) (PeerRecord, error) {
	var record PeerRecord

	bridge := newBridgeCtx()
	defer bridge.free()

	var cPeerId = C.CString(peerId)
	defer C.free(unsafe.Pointer(cPeerId))

	if C.cGoStoragePeerDebug(node.ctx, cPeerId, bridge.resp) != C.RET_OK {
		return record, bridge.callError("cGoStoragePeerDebug")
	}

	value, err := bridge.wait()
	if err != nil {
		return record, err
	}

	err = json.Unmarshal([]byte(value), &record)
	return record, err
}
