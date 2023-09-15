package header

import (
	"context"
	"fmt"
	"strings"
)

// Header is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Header map[string][]string

// New creates an h from a given key-values map.
func New(hs ...map[string][]string) Header {
	h := Header{}
	for _, s := range hs {
		for k, vList := range s {
			for _, v := range vList {
				h.Add(k, v)
			}
		}
	}
	return h
}

// Add adds the key, value pair to the header.
func (s Header) Add(key, value string) {
	if len(key) == 0 {
		return
	}

	s[strings.ToLower(key)] = append(s[strings.ToLower(key)], value)
}

// Get returns the value associated with the passed key.
func (s Header) Get(key string) string {
	v := s[strings.ToLower(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Set stores the key-value pair.
func (s Header) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	s[strings.ToLower(key)] = []string{value}
}

// Range iterate over element in Header.
func (s Header) Range(f func(k string, v []string) bool) {
	for k, v := range s {
		if !f(k, v) {
			break
		}
	}
}

// Values returns a slice of values associated with the passed key.
func (s Header) Values(key string) []string {
	return s[strings.ToLower(key)]
}

// Clone returns a deep copy of Header
func (s Header) Clone() Header {
	h := make(Header, len(s))
	for k, v := range s {
		h[k] = v
	}
	return h
}

type headerKey struct{}

// NewContext creates a new context with h attached.
func NewContext(ctx context.Context, h Header) context.Context {
	return context.WithValue(ctx, headerKey{}, h)
}

// FromContext returns the Header in ctx if it exists.
func FromContext(ctx context.Context) (Header, bool) {
	h, ok := ctx.Value(headerKey{}).(Header)
	return h, ok
}

// AppendToContext returns a new context with the provided kv merged
// with any existing Header in the context.
func AppendToContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("Header: AppendToContext got an odd number of input pairs for Header: %d", len(kv)))
	}
	h, _ := FromContext(ctx)
	h = h.Clone()
	for i := 0; i < len(kv); i += 2 {
		h.Set(kv[i], kv[i+1])
	}
	return NewContext(ctx, h)
}

// MergeToContext merge new Header into ctx.
func MergeToContext(ctx context.Context, hd Header) context.Context {
	h, _ := FromContext(ctx)
	h = h.Clone()
	for k, v := range hd {
		h[k] = v
	}
	return NewContext(ctx, h)
}
