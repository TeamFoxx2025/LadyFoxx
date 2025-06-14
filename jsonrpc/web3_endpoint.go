package jsonrpc

import (
	"fmt"
	"runtime"

	"github.com/TeamFoxx2025/LadyFoxx/helper/keccak"
	"github.com/TeamFoxx2025/LadyFoxx/versioning"
)

// Web3 is the web3 jsonrpc endpoint
type Web3 struct {
	chainName string
}

var clientVersionTemplate = "%s/%s/%s-%s/%s"

// ClientVersion returns the version of the web3 client (web3_clientVersion)
// Spec: https://ethereum.org/en/developers/docs/apis/json-rpc/#web3_clientversion
func (w *Web3) ClientVersion() (interface{}, error) {
	var version string
	if versioning.Version != "" {
		version = versioning.Version
	} else if versioning.Commit != "" {
		version = versioning.Commit[:8]
	}

	return fmt.Sprintf(
		clientVersionTemplate,
		w.chainName,
		version,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version(),
	), nil
}

// Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data
func (w *Web3) Sha3(v argBytes) (interface{}, error) {
	dst := keccak.Keccak256(nil, v)

	return argBytes(dst), nil
}
