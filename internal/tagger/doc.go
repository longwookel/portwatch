// Package tagger maps port numbers to human-readable service labels.
//
// Built-in mappings cover the most common well-known ports (SSH, HTTP, HTTPS,
// common databases, etc.). Callers may supply additional custom mappings at
// construction time; custom labels always take precedence over the built-ins.
//
// Usage:
//
//	t := tagger.New(map[uint16]string{9200: "elasticsearch"})
//	fmt.Println(t.Tag(22))   // "ssh"
//	fmt.Println(t.Tag(9200)) // "elasticsearch"
//	fmt.Println(t.Tag(1234)) // "unknown/1234"
package tagger
