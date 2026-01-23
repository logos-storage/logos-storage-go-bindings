package storage

import (
	"encoding/json"
	"unsafe"
)

/*
   #include "bridge.h"
   #include <stdlib.h>

   static int cGoStorageStorageList(void* storageCtx, void* resp) {
       return storage_list(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageStorageFetch(void* storageCtx, char* cid, void* resp) {
       return storage_fetch(storageCtx, cid, (StorageCallback) callback, resp);
   }

   static int cGoStorageStorageSpace(void* storageCtx, void* resp) {
       return storage_space(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageStorageDelete(void* storageCtx, char* cid, void* resp) {
       return storage_delete(storageCtx, cid, (StorageCallback) callback, resp);
   }

   static int cGoStorageStorageExists(void* storageCtx, char* cid, void* resp) {
       return storage_exists(storageCtx, cid, (StorageCallback) callback, resp);
   }
*/
import "C"

type manifestWithCid struct {
	Cid      string   `json:"cid"`
	Manifest Manifest `json:"manifest"`
}

type Space struct {
	// TotalBlocks is the number of blocks stored by the node
	TotalBlocks int `json:"totalBlocks"`

	// QuotaMaxBytes is the maximum storage space (in bytes) available
	// for the node in Logos Storage's local repository.
	QuotaMaxBytes int64 `json:"quotaMaxBytes"`

	// QuotaUsedBytes is the mount of storage space (in bytes) currently used
	// for storing files in Logos Storage's local repository.
	QuotaUsedBytes int64 `json:"quotaUsedBytes"`

	// QuotaReservedBytes is the amount of storage reserved (in bytes) in the
	// Logos Storage's local repository for future use when storage requests will be picked
	// up and hosted by the node using node's availabilities.
	// This does not include the storage currently in use.
	QuotaReservedBytes int64 `json:"quotaReservedBytes"`
}

// Manifests returns the list of all manifests stored by the Logos Storage node.
func (node StorageNode) Manifests() ([]Manifest, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageStorageList(node.ctx, bridge.resp) != C.RET_OK {
		return nil, bridge.callError("cGoStorageStorageList")
	}
	value, err := bridge.wait()
	if err != nil {
		return nil, err
	}

	var items []manifestWithCid
	err = json.Unmarshal([]byte(value), &items)
	if err != nil {
		return nil, err
	}

	var list []Manifest
	for _, item := range items {
		item.Manifest.Cid = item.Cid
		list = append(list, item.Manifest)
	}

	return list, err
}

// Fetch download a file from the network and store it to the local node.
func (node StorageNode) Fetch(cid string) (Manifest, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	var cCid = C.CString(cid)
	defer C.free(unsafe.Pointer(cCid))

	if C.cGoStorageStorageFetch(node.ctx, cCid, bridge.resp) != C.RET_OK {
		return Manifest{}, bridge.callError("cGoStorageStorageFetch")
	}

	value, err := bridge.wait()
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = json.Unmarshal([]byte(value), &manifest)
	if err != nil {
		return Manifest{}, err
	}

	manifest.Cid = cid
	return manifest, nil
}

// Space returns information about the storage space used and available.
func (node StorageNode) Space() (Space, error) {
	var space Space

	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageStorageSpace(node.ctx, bridge.resp) != C.RET_OK {
		return space, bridge.callError("cGoStorageStorageSpace")
	}

	value, err := bridge.wait()
	if err != nil {
		return space, err
	}

	err = json.Unmarshal([]byte(value), &space)
	return space, err
}

// Deletes either a single block or an entire dataset
// from the local node. Does nothing if the dataset is not locally available.
func (node StorageNode) Delete(cid string) error {
	bridge := newBridgeCtx()
	defer bridge.free()

	var cCid = C.CString(cid)
	defer C.free(unsafe.Pointer(cCid))

	if C.cGoStorageStorageDelete(node.ctx, cCid, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoStorageStorageDelete")
	}

	_, err := bridge.wait()
	return err
}

// Exists checks if a given cid exists in the local storage.
func (node StorageNode) Exists(cid string) (bool, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	var cCid = C.CString(cid)
	defer C.free(unsafe.Pointer(cCid))

	if C.cGoStorageStorageExists(node.ctx, cCid, bridge.resp) != C.RET_OK {
		return false, bridge.callError("cGoStorageStorageExists")
	}

	result, err := bridge.wait()
	return result == "true", err
}
