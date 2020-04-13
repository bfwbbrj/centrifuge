package recovery

import (
	"math"
	"sort"

	"github.com/centrifugal/protocol"
)

const (
	maxSeq uint32 = math.MaxUint32 // maximum uint32 value
)

// uniquePublications returns slice of unique Publications.
func uniquePublications(s []*protocol.Publication) []*protocol.Publication {
	keys := make(map[uint64]struct{})
	var list []*protocol.Publication
	for _, entry := range s {
		val := (uint64(entry.Seq))<<32 | uint64(entry.Gen)
		if _, value := keys[val]; !value {
			keys[val] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}

// MergePublications ...
func MergePublications(recoveredPubs []*protocol.Publication, bufferedPubs []*protocol.Publication) ([]*protocol.Publication, bool) {
	if len(bufferedPubs) > 0 {
		recoveredPubs = append(recoveredPubs, bufferedPubs...)
	}
	sort.Slice(recoveredPubs, func(i, j int) bool {
		if recoveredPubs[i].Gen != recoveredPubs[j].Gen {
			return recoveredPubs[i].Gen > recoveredPubs[j].Gen
		}
		return recoveredPubs[i].Seq > recoveredPubs[j].Seq
	})
	if len(bufferedPubs) > 0 {
		recoveredPubs = uniquePublications(recoveredPubs)
		prevSeq := PackUint64(recoveredPubs[0].Seq, recoveredPubs[0].Gen)
		for _, p := range recoveredPubs[1:] {
			pubSequence := PackUint64(p.Seq, p.Gen)
			if pubSequence != prevSeq+1 {
				return nil, false
			}
			prevSeq = pubSequence
		}
	}
	return recoveredPubs, true
}

// NextSeqGen ...
func NextSeqGen(currentSeq, currentGen uint32) (uint32, uint32) {
	var nextSeq uint32
	nextGen := currentGen
	if currentSeq == maxSeq {
		nextSeq = 0
		nextGen++
	} else {
		nextSeq = currentSeq + 1
	}
	return nextSeq, nextGen
}

// PackUint64 ...
func PackUint64(seq, gen uint32) uint64 {
	return uint64(gen)*uint64(math.MaxUint32) + uint64(seq)
}

// UnpackUint64 ...
func UnpackUint64(val uint64) (uint32, uint32) {
	return uint32(val), uint32(val >> 32)
}