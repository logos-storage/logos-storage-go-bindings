package storage

/*
   #include "bridge.h"
   #include <stdlib.h>

   static int cGoStorageTogglePrivateQueries(void* storageCtx, bool enabled, void* resp) {
       return storage_toggle_private_queries(storageCtx, enabled, (StorageCallback) callback, resp);
   }
*/
import "C"

// TogglePrivateQueries toggles routing of DHT queries over the Logos mix network.
//
// When enabled, all subsequent DHT queries are tunnelled over Mix; this
// affects queries only, not advertisements.
//
// Enabling requires Mix to be configured: MixEnabled true and at least one
// DhtMixProxies set (see Config). Otherwise enabling fails with an error.
// Disabling is always allowed.
//
// This is a temporary API and will likely be removed before mainnet.
//
// On success, it returns the previous toggle state (true = private queries
// were already enabled).
func (node StorageNode) TogglePrivateQueries(enabled bool) (bool, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageTogglePrivateQueries(node.ctx, C.bool(enabled), bridge.resp) != C.RET_OK {
		return false, bridge.callError("cGoStorageTogglePrivateQueries")
	}

	previous, err := bridge.wait()
	if err != nil {
		return false, err
	}

	return previous == "true", nil
}
