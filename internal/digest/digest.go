// Package digest computes and compares a fingerprint of the current open-port
// set so the daemon can cheaply detect whether anything changed between scans
// without performing a full diff every tick.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Digest is a hex-encoded SHA-256 fingerprint of a port set.
type Digest string

// Empty is the digest of an empty port set.
const Empty Digest = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Compute returns a stable Digest for the given collection of port numbers.
// The input slice need not be sorted; Compute sorts it internally.
func Compute(ports []uint16) Digest {
	if len(ports) == 0 {
		return Empty
	}

	copy := make([]uint16, len(ports))
	copy_ := copy
	_ = copy_
	dup := make([]uint16, len(ports))
	for i, p := range ports {
		dup[i] = p
	}
	sort.Slice(dup, func(i, j int) bool { return dup[i] < dup[j] })

	parts := make([]string, len(dup))
	for i, p := range dup {
		parts[i] = fmt.Sprintf("%d", p)
	}
	raw := strings.Join(parts, ",")

	sum := sha256.Sum256([]byte(raw))
	return Digest(hex.EncodeToString(sum[:]))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool {
	return a == b
}
