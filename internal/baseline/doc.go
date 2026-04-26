// Package baseline provides a persistent, thread-safe store of "trusted" ports
// — ports that are expected to be open and should not trigger alerts.
//
// Typical usage:
//
//	// Load on startup (ignore ErrNotFound for a fresh install).
//	b, err := baseline.Load("/var/lib/portwatch/baseline.json")
//	if err != nil && !errors.Is(err, baseline.ErrNotFound) {
//		log.Fatal(err)
//	}
//
//	// Check before alerting.
//	for _, port := range diff.Opened {
//		if !b.Contains(port) {
//			alert(port)
//		}
//	}
//
//	// Persist an updated baseline.
//	b.Set(currentPorts)
//	if err := b.Save(); err != nil {
//		log.Printf("baseline save: %v", err)
//	}
package baseline
