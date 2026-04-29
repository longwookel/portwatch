// Package circuitbreaker provides a thread-safe circuit breaker for use in
// portwatch components that interact with external systems (e.g. webhook
// senders, audit log writers).
//
// # States
//
// The breaker moves through three states:
//
//   - Closed  — normal operation; all calls are allowed through.
//   - Open    — the failure threshold was exceeded; calls return ErrOpen until
//     the reset timeout elapses.
//   - HalfOpen — one probe call is permitted; a success closes the breaker
//     while another failure re-opens it and restarts the timeout.
//
// # Usage
//
//	b := circuitbreaker.New(5, 30*time.Second)
//
//	if err := b.Allow(); err != nil {
//		// circuit is open — skip the operation
//		return err
//	}
//	if err := doWork(); err != nil {
//		b.RecordFailure()
//		return err
//	}
//	b.RecordSuccess()
package circuitbreaker
