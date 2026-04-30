// Package dedup implements diff-level deduplication for portwatch.
//
// When the port scanner detects a change it produces a [snapshot.Diff].
// If the daemon is configured with a short polling interval the same
// transition (e.g. port 8080 opened) can appear in many consecutive
// diffs before the underlying state settles.  Forwarding every one of
// those diffs to downstream sinks — alerters, audit logs, webhooks —
// creates noise and can trigger rate-limit errors on external services.
//
// [Deduplicator] solves this by computing a compact signature for each
// diff and remembering the last time that signature was forwarded.  Any
// identical diff that arrives within the configured suppression window
// is silently dropped.
//
// # Usage
//
//	dd := dedup.New(30 * time.Second)
//
//	for diff := range diffStream {
//	    if dd.IsDuplicate(diff) {
//	        continue
//	    }
//	    forwardToAlerts(diff)
//	}
//
// Call [Deduplicator.Purge] periodically (e.g. once per minute) to
// reclaim memory used by expired cache entries.
package dedup
