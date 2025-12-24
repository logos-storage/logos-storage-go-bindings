package storage

/*
   #include "bridge.h"
   #include <stdlib.h>

   void libstorageNimMain(void);

   static void storage_host_init_once(void){
       static int done;
       if (!__atomic_exchange_n(&done, 1, __ATOMIC_SEQ_CST)) libstorageNimMain();
   }

   // resp must be set != NULL in case interest on retrieving data from the callback
   void callback(int ret, char* msg, size_t len, void* resp);

   static void* cGoStorageNew(const char* configJson, void* resp) {
       void* ret = storage_new(configJson, (StorageCallback) callback, resp);
       return ret;
   }

   static int cGoStorageStart(void* storageCtx, void* resp) {
       return storage_start(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageStop(void* storageCtx, void* resp) {
       return storage_stop(storageCtx, (StorageCallback) callback, resp);
   }

	static int cGoStorageClose(void* storageCtx, void* resp) {
		return storage_close(storageCtx, (StorageCallback) callback, resp);
	}

   static int cGoStorageDestroy(void* storageCtx, void* resp) {
       return storage_destroy(storageCtx, (StorageCallback) callback, resp);
   }

    static int cGoStorageVersion(void* storageCtx, void* resp) {
       return storage_version(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageRevision(void* storageCtx, void* resp) {
       return storage_revision(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageRepo(void* storageCtx, void* resp) {
       return storage_repo(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStorageSpr(void* storageCtx, void* resp) {
       return storage_spr(storageCtx, (StorageCallback) callback, resp);
   }

   static int cGoStoragePeerId(void* storageCtx, void* resp) {
       return storage_peer_id(storageCtx, (StorageCallback) callback, resp);
   }
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

type LogLevel string

const (
	TRACE  LogLevel = "trace"
	DEBUG  LogLevel = "debug"
	INFO   LogLevel = "info"
	NOTICE LogLevel = "notice"
	WARN   LogLevel = "warn"
	ERROR  LogLevel = "error"
	FATAL  LogLevel = "fatal"
)

type LogFormat string

const (
	LogFormatAuto     LogFormat = "auto"
	LogFormatColors   LogFormat = "colors"
	LogFormatNoColors LogFormat = "nocolors"
	LogFormatJSON     LogFormat = "json"
)

type RepoKind string

const (
	FS      RepoKind = "fs"
	SQLite  RepoKind = "sqlite"
	LevelDb RepoKind = "leveldb"
)

type Config struct {
	// Default: INFO
	LogLevel string `json:"log-level,omitempty"`

	// Specifies what kind of logs should be written to stdout
	// Default: auto
	LogFormat LogFormat `json:"log-format,omitempty"`

	// Enable the metrics server
	// Default: false
	MetricsEnabled bool `json:"metrics,omitempty"`

	// Listening address of the metrics server
	// Default: 127.0.0.1
	MetricsAddress string `json:"metrics-address,omitempty"`

	// Listening HTTP port of the metrics server
	// Default: 8008
	MetricsPort int `json:"metrics-port,omitempty"`

	// The directory where Logos Storage will store configuration and data
	// Default:
	// $HOME\AppData\Roaming\Storage on Windows
	// $HOME/Library/Application Support/Storage on macOS
	// $HOME/.cache/storage on Linux
	DataDir string `json:"data-dir,omitempty"`

	// Multi Addresses to listen on
	// Default: ["/ip4/0.0.0.0/tcp/0"]
	ListenAddrs []string `json:"listen-addrs,omitempty"`

	// Specify method to use for determining public address.
	// Must be one of: any, none, upnp, pmp, extip:<IP>
	// Default: any
	Nat string `json:"nat,omitempty"`

	// Discovery (UDP) port
	// Default: 8090
	DiscoveryPort int `json:"disc-port,omitempty"`

	// Source of network (secp256k1) private key file path or name
	// Default: "key"
	NetPrivKeyFile string `json:"net-privkey,omitempty"`

	// Specifies one or more bootstrap nodes to use when connecting to the network.
	BootstrapNodes []string `json:"bootstrap-node,omitempty"`

	// The maximum number of peers to connect to.
	// Default: 160
	MaxPeers int `json:"max-peers,omitempty"`

	// Number of worker threads (\"0\" = use as many threads as there are CPU cores available)
	// Default: 0
	NumThreads int `json:"num-threads,omitempty"`

	// Node agent string which is used as identifier in network
	// Default: "Logos Storage"
	AgentString string `json:"agent-string,omitempty"`

	// Backend for main repo store (fs, sqlite, leveldb)
	// Default: fs
	RepoKind RepoKind `json:"repo-kind,omitempty"`

	// The size of the total storage quota dedicated to the node
	// Default: 20 GiBs
	StorageQuota int `json:"storage-quota,omitempty"`

	// Default block timeout in seconds - 0 disables the ttl
	// Default: 30 days
	BlockTtl string `json:"block-ttl,omitempty"`

	// Time interval in seconds - determines frequency of block
	// maintenance cycle: how often blocks are checked for expiration and cleanup
	// Default: 10 minutes
	BlockMaintenanceInterval string `json:"block-mi,omitempty"`

	// Number of blocks to check every maintenance cycle
	// Default: 1000
	BlockMaintenanceNumberOfBlocks int `json:"block-mn,omitempty"`

	// Number of times to retry fetching a block before giving up
	// Default: 3000
	BlockRetries int `json:"block-retries,omitempty"`

	// The size of the block cache, 0 disables the cache -
	// might help on slow hardrives
	// Default: 0
	CacheSize int `json:"cache-size,omitempty"`

	// Default: "" (no log file)
	LogFile string `json:"log-file,omitempty"`
}

type StorageNode struct {
	ctx unsafe.Pointer
}

type ChunkSize int

func (c ChunkSize) valOrDefault() int {
	if c == 0 {
		return defaultBlockSize
	}

	return int(c)
}

func (c ChunkSize) toSizeT() C.size_t {
	return C.size_t(c.valOrDefault())
}

// New creates a new Logos Storage node with the provided configuration.
// The node is not started automatically; you need to call StorageStart
// to start it.
// It returns a Logos Storage node that can be used to interact
// with the Logos Storage network.
func New(config Config) (*StorageNode, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	cJsonConfig := C.CString(string(jsonConfig))
	defer C.free(unsafe.Pointer(cJsonConfig))

	ctx := C.cGoStorageNew(cJsonConfig, bridge.resp)

	if _, err := bridge.wait(); err != nil {
		return nil, bridge.err
	}

	return &StorageNode{ctx: ctx}, bridge.err
}

// Start starts the Logos Storage node.
func (node StorageNode) Start() error {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageStart(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoStorageStart")
	}

	_, err := bridge.wait()
	return err
}

// StartAsync is the asynchronous version of Start.
func (node StorageNode) StartAsync(onDone func(error)) {
	go func() {
		err := node.Start()
		onDone(err)
	}()
}

// Stop stops the Logos Storage node.
func (node StorageNode) Stop() error {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageStop(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoStorageStop")
	}

	_, err := bridge.wait()
	return err
}

// Destroy destroys the Logos Storage node, freeing all resources.
// The node must be stopped before calling this method.
func (node StorageNode) Destroy() error {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageClose(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoStorageClose")
	}

	_, err := bridge.wait()
	if err != nil {
		return err
	}

	if C.cGoStorageDestroy(node.ctx, bridge.resp) != C.RET_OK {
		return bridge.callError("cGoStorageDestroy")
	}

	// We don't wait for the bridge here.
	// The destroy function does not call the worker thread,
	// it destroys the context directly and return the return
	// value synchronously.

	return nil
}

// Version returns the version of the Logos Storage node.
func (node StorageNode) Version() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageVersion(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoStorageVersion")
	}

	return bridge.wait()
}

func (node StorageNode) Revision() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageRevision(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoStorageRevision")
	}

	return bridge.wait()
}

// Repo returns the path of the data dir folder.
func (node StorageNode) Repo() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageRepo(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoStorageRepo")
	}

	return bridge.wait()
}

func (node StorageNode) Spr() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStorageSpr(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoStorageSpr")
	}

	return bridge.wait()
}

func (node StorageNode) PeerId() (string, error) {
	bridge := newBridgeCtx()
	defer bridge.free()

	if C.cGoStoragePeerId(node.ctx, bridge.resp) != C.RET_OK {
		return "", bridge.callError("cGoStoragePeerId")
	}

	return bridge.wait()
}
