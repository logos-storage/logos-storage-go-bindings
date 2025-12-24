package storage

/*
   #include "bridge.h"
   #include <stdlib.h>

   static int cGoStorageConnect(void* storageCtx, char* peerId, const char** peerAddresses, uintptr_t peerAddressesSize,  void* resp) {
       return storage_connect(storageCtx, peerId, peerAddresses, peerAddressesSize, (StorageCallback) callback, resp);
   }
*/
import "C"
import (
	"unsafe"
)

// Connect connects to a peer using its peer ID and optional multiaddresses.
// If `peerAddresses` param is supplied, it will be used to  dial the peer,
// otherwise the `peerId` is used to invoke peer discovery, if it succeeds
// the returned addresses will be used to dial.
// `peerAddresses` the listening addresses of the peers to dial,
// eg the one specified with `ListenAddresses` in `StorageConfig`.
func (node StorageNode) Connect(peerId string, peerAddresses []string) error {
	bridge := newBridgeCtx()
	defer bridge.free()

	var cPeerId = C.CString(peerId)
	defer C.free(unsafe.Pointer(cPeerId))

	if len(peerAddresses) > 0 {
		var cAddresses = make([]*C.char, len(peerAddresses))
		for i, addr := range peerAddresses {
			cAddresses[i] = C.CString(addr)
			defer C.free(unsafe.Pointer(cAddresses[i]))
		}

		if C.cGoStorageConnect(node.ctx, cPeerId, &cAddresses[0], C.uintptr_t(len(peerAddresses)), bridge.resp) != C.RET_OK {
			return bridge.callError("cGoStorageConnect")
		}
	} else {
		if C.cGoStorageConnect(node.ctx, cPeerId, nil, 0, bridge.resp) != C.RET_OK {
			return bridge.callError("cGoStorageConnect")
		}
	}

	_, err := bridge.wait()
	return err
}
