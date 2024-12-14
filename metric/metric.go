package metric

import "sync/atomic"

func Incr(addr *uint64) {
	atomic.AddUint64(addr, 1)
}

var (
	entryCount             uint64
	probeCount             uint64
	probeSuccessCount      uint64
	probeFailedCount       uint64
	selectCount            uint64
	selectSuccessCount     uint64
	selectFailedCount      uint64
	selectDataMissingCount uint64
	acceptCount            uint64
)

func IncrEntryCount()             { Incr(&entryCount) }
func IncrProbeCount()             { Incr(&probeCount) }
func IncrProbeSuccessCount()      { Incr(&probeSuccessCount) }
func IncrProbeFailedCount()       { Incr(&probeFailedCount) }
func IncrSelectCount()            { Incr(&selectCount) }
func IncrSelectSuccessCount()     { Incr(&selectSuccessCount) }
func IncrSelectFailedCount()      { Incr(&selectFailedCount) }
func IncrSelectDataMissingCount() { Incr(&selectDataMissingCount) }
func IncrAcceptCount()            { Incr(&acceptCount) }

type Metrics struct {
	EntryCount             uint64
	ProbeCount             uint64
	ProbeSuccessCount      uint64
	ProbeFailedCount       uint64
	SelectCount            uint64
	SelectSuccessCount     uint64
	SelectFailedCount      uint64
	SelectDataMissingCount uint64
	AcceptCount            uint64
}

func Get() *Metrics {
	return &Metrics{
		EntryCount:             entryCount,
		ProbeCount:             probeCount,
		ProbeSuccessCount:      probeSuccessCount,
		ProbeFailedCount:       probeFailedCount,
		SelectCount:            selectCount,
		SelectSuccessCount:     selectSuccessCount,
		SelectFailedCount:      selectFailedCount,
		SelectDataMissingCount: selectDataMissingCount,
		AcceptCount:            acceptCount,
	}
}
