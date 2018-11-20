package track

import "github.com/ShyftNetwork/go-empyrean/common"

type InternalTracker interface {
	TraceTransaction(hash common.Hash, bHash common.Hash) (interface{}, error)
}
