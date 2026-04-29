// Package digest provides a lightweight fingerprinting mechanism for sets of
// open ports.
//
// # Overview
//
// Each scan produces a slice of open port numbers. Converting that slice into
// a [Digest] (a hex-encoded SHA-256 hash of the sorted, comma-separated port
// list) lets the daemon skip expensive diffing and alerting when nothing has
// changed between two consecutive scans.
//
// # Usage
//
//	d := digest.Compute(openPorts)
//	if cache.Changed("localhost", d) {
//	    // run full diff and alert pipeline
//	}
//
// The [Cache] type is safe for concurrent use by multiple goroutines.
package digest
