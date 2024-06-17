package rsmt2d

import (
	"errors"
	"fmt"
)

const (
	// Leopard is a codec that was originally implemented in the C++ library
	// https://github.com/catid/leopard. rsmt2d uses a Go port of the C++
	// library in https://github.com/klauspost/reedsolomon. The Leopard codec
	// uses 8-bit leopard for shares less than or equal to 256. The Leopard
	// codec uses 16-bit leopard for shares greater than 256.
	Leopard = "Leopard"
)

// It will be returned from ReconstructSome if there were too few shards
// to reconstruct the missing data.
var ErrTooFewShares = errors.New("too few shards given")

type Codec interface {
	// Encode encodes original data, automatically extracting share size.
	// There must be no missing shares. Only returns parity shares.
	Encode(data [][]byte) ([][]byte, error)
	// Decode decodes sparse original + parity data, automatically extracting share size.
	// Missing shares must be nil. Returns original + parity data.
	Decode(data [][]byte) ([][]byte, error)
	// ReconstructSome will recreate only requested data, if possible.
	//
	// Given a list of shares, some of which contain data, fills in the
	// shards indicated by true values in the "required" parameter.
	// The length of the "required" array must be equal to either original data  or full data (original+parity) .
	// If the length is equal to original, the reconstruction of parity shards will be ignored.
	//
	// The length of data array must be equal to size of original + parity data.
	// You indicate that a shard is missing by setting it to nil or zero-length.
	// If a shard is zero-length but has sufficient capacity, that memory will
	// be used, otherwise a new []byte will be allocated.
	//
	// If there are too few shares to reconstruct the missing
	// ones, ErrTooFewShares will be returned.
	ReconstructSome(data [][]byte, required []bool) error
	// MaxChunks returns the max number of chunks this codec supports in a 2D
	// original data square. Chunk is a synonym of share.
	MaxChunks() int
	// Name returns the name of the codec.
	Name() string
	// ValidateChunkSize returns an error if this codec does not support
	// chunkSize. Returns nil if chunkSize is supported. Chunk is a synonym of
	// share.
	ValidateChunkSize(chunkSize int) error
}

// codecs is a global map used for keeping track of registered codecs for testing and JSON unmarshalling
var codecs = make(map[string]Codec)

func registerCodec(ct string, codec Codec) {
	if codecs[ct] != nil {
		panic(fmt.Sprintf("%v already registered", codec))
	}
	codecs[ct] = codec
}
