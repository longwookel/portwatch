// Package sampler implements adaptive scan interval adjustment for portwatch.
//
// The Sampler observes whether recent scans detected port changes and
// adjusts the recommended interval accordingly:
//
//   - When a change is recorded the interval shrinks by StepDown, down to Min.
//   - When no change is recorded the interval grows by StepUp, up to Max.
//
// This allows portwatch to react quickly during volatile periods while
// reducing unnecessary CPU and I/O overhead on stable hosts.
//
// Typical usage:
//
//	sampler := sampler.New(sampler.Default())
//
//	for {
//		time.Sleep(sampler.Current())
//		diff := scan()
//		if len(diff.Opened)+len(diff.Closed) > 0 {
//			sampler.RecordChange()
//		}
//		_ = sampler.Next()
//	}
package sampler
