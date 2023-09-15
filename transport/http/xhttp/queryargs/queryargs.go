package queryargs

import (
	"context"
	"fmt"
)

// QueryArgs is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type QA struct {
	Key   string
	Value string
}
type QueryArgs []*QA

// New creates an MD from a given key-values slice.
func New(kvs ...string) QueryArgs {
	l := len(kvs)
	if l%2 == 1 {
		panic(fmt.Sprintf("QueryArgs: NewQueryArgs got an odd number of input pairs for QueryArgs: %d", len(kvs)))
	}
	qas := make([]*QA, 0, l)
	for i := 0; i < len(kvs); i += 2 {
		qas = append(qas, &QA{kvs[i], kvs[i+1]})
	}
	return qas
}

// Add adds the key, value pair to the QueryArgs.
func (q QueryArgs) Add(key, value string) QueryArgs {
	if len(key) == 0 {
		return q
	}
	return append(q, &QA{key, value})
}

// Range iterate over element in QueryArgs.
func (q QueryArgs) Range(f func(k, v string) bool) {
	for _, v := range q {
		if !f(v.Key, v.Value) {
			break
		}
	}
}

// Clone returns a deep copy of QueryArgs
func (q QueryArgs) Clone() QueryArgs {
	qa := make(QueryArgs, 0, len(q))
	for _, v := range q {
		qa = append(qa, &QA{v.Key, v.Value})
	}
	return qa
}

type queryArgsKey struct{}

// NewContext creates a new context with qa attached.
func NewContext(ctx context.Context, qa QueryArgs) context.Context {
	return context.WithValue(ctx, queryArgsKey{}, qa)
}

// FromContext returns the QueryArgs in ctx if it exists.
func FromContext(ctx context.Context) (QueryArgs, bool) {
	qa, ok := ctx.Value(queryArgsKey{}).(QueryArgs)
	return qa, ok
}

// AppendToContext returns a new context with the provided kv merged
// with any existing QueryArgs in the context.
func AppendToContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("QueryArgs: AppendToContext got an odd number of input pairs for QueryArgs: %d", len(kv)))
	}
	qa, _ := FromContext(ctx)
	qa = qa.Clone()
	for i := 0; i < len(kv); i += 2 {
		qa = qa.Add(kv[i], kv[i+1])
	}
	return NewContext(ctx, qa)
}

// MergeToContext merge new QueryArgs into ctx.
func MergeToContext(ctx context.Context, cqa QueryArgs) context.Context {
	qa, _ := FromContext(ctx)
	qa = qa.Clone()
	qa = append(qa, cqa...)
	return NewContext(ctx, qa)
}
