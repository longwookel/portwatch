// Package ratelimit provides a token-bucket rate limiter for controlling
// the frequency of outbound notifications and alert dispatches.
//
// A Limiter is initialised with a burst capacity and a refill rate. Each
// call to Allow consumes one token; tokens are replenished at the given
// rate over time. When the bucket is empty Allow returns false and the
// caller should suppress the action until tokens become available.
//
// Example usage:
//
//	rl := ratelimit.New(10, time.Minute) // 10 alerts per minute
//	if rl.Allow("port-scan") {
//		sender.Send(ctx, msg)
//	}
package ratelimit
