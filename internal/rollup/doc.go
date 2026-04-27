// Package rollup provides an event-coalescing buffer for port-change Diffs.
//
// During rapid port-state churn (e.g. a service restart cycling through
// dozens of ephemeral ports) the scanner may emit many Diff values in quick
// succession. Sending an individual alert for each would overwhelm both the
// operator and any downstream webhook receiver.
//
// A Rollup solves this by collecting every Diff that arrives within a
// configurable quiet-period window and merging them into a single value
// before handing it to the caller's Flusher function.
//
// Typical usage:
//
//	r := rollup.New(2*time.Second, func(d snapshot.Diff) {
//		notifier.Notify(ctx, d)
//	})
//
//	// Inside the daemon scan loop:
//	r.Add(diff)
//
//	// On shutdown, flush any buffered events:
//	r.Flush()
package rollup
