// Package tagger assigns human-readable labels to ports based on well-known
// service mappings and user-defined rules.
package tagger

import "fmt"

// well-known maps common port numbers to service names.
var wellKnown = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger labels ports with service names.
type Tagger struct {
	custom map[uint16]string
}

// New returns a Tagger seeded with the provided custom labels.
// Custom labels take precedence over built-in well-known mappings.
func New(custom map[uint16]string) *Tagger {
	c := make(map[uint16]string, len(custom))
	for k, v := range custom {
		c[k] = v
	}
	return &Tagger{custom: c}
}

// Tag returns the label for port p. It checks custom rules first, then
// well-known mappings. If no match is found it returns "unknown/<port>".
func (t *Tagger) Tag(p uint16) string {
	if label, ok := t.custom[p]; ok {
		return label
	}
	if label, ok := wellKnown[p]; ok {
		return label
	}
	return fmt.Sprintf("unknown/%d", p)
}

// TagAll returns a map of port → label for every port in the slice.
func (t *Tagger) TagAll(ports []uint16) map[uint16]string {
	out := make(map[uint16]string, len(ports))
	for _, p := range ports {
		out[p] = t.Tag(p)
	}
	return out
}
