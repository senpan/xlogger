package logtrace

import (
	"context"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	TraceMetaKeyTraceId   = "x_trace_id"
	TraceMetaKeyRpcId     = "x_rpcid"
	TraceMetaKeyName      = "x_name"
	TraceMetaKeyRequestId = "x-request-id"
	TraceMetaKeyB3TraceId = "x-b3-traceid"
)

// GetMetadataKey context key
func GetMetadataKey() string {
	return "__logtrace_metadata_key__"
}

var enableMetadataKey map[string]bool = map[string]bool{TraceMetaKeyTraceId: true, TraceMetaKeyRpcId: true,
	TraceMetaKeyRequestId: true, TraceMetaKeyB3TraceId: true}

// AppendLogTraceMetadataContext metadata转ctx
func AppendLogTraceMetadataContext(ctx context.Context, metadata map[string]string) context.Context {
	if metadata == nil {
		return ctx
	}
	t := NewTraceNode()
	for k, v := range metadata {
		if _, ok := enableMetadataKey[k]; ok {
			t.Set(k, v)
		}
	}
	ctx = context.WithValue(ctx, GetMetadataKey(), t)
	return ctx
}

// GenLogTraceMetadata InitTraceNode
func GenLogTraceMetadata() *TraceNode {
	t := NewTraceNode()
	traceId := NewTraceId()
	t.Set(TraceMetaKeyTraceId, traceId)
	t.Set(TraceMetaKeyRpcId, "0.1")
	t.Set(TraceMetaKeyRequestId, traceId)
	return t
}

func NewTraceId() string {
	return uuid.NewV4().String()
}

// ExtractTraceNodeFromContext Get TraceNode
func ExtractTraceNodeFromContext(ctx context.Context) *TraceNode {
	if ctx == nil {
		return NewTraceNode()
	}
	meta := ctx.Value(GetMetadataKey())
	if meta == nil {
		return NewTraceNode()
	} else {
		if val, ok := meta.(*TraceNode); ok {
			return val
		} else {
			return NewTraceNode()
		}
	}
}

// InjectMetadata TraceNode add other kv
func InjectMetadata(ctx context.Context, mapPtr *map[string]string) bool {
	meta := ExtractTraceNodeFromContext(ctx)
	traceRpcId := meta.Get(TraceMetaKeyRpcId)
	if len(traceRpcId) == 0 {
		return false
	}
	for k, v := range meta.ForkMap() {
		(*mapPtr)[k] = v
	}
	return true
}

// IncrementRpcId Inc RpcId 1.1.1=>1.1.2
func IncrementRpcId(ctx context.Context) bool {
	meta := ExtractTraceNodeFromContext(ctx)
	traceRpcId := meta.Get(TraceMetaKeyRpcId)
	if len(traceRpcId) == 0 {
		return false
	}

	index := strings.LastIndex(traceRpcId, ".")
	if index == -1 {
		return false
	}
	index += 1

	id, err := strconv.Atoi(traceRpcId[index:])
	if err != nil {
		return false
	}

	id += 1
	meta.Set(TraceMetaKeyRpcId, traceRpcId[0:index]+strconv.Itoa(id))
	return true
}

// AppendNewRpcId Append RpcId  1.1.1=>1.1.1.0
func AppendNewRpcId(ctx context.Context) bool {
	meta := ExtractTraceNodeFromContext(ctx)
	traceRpcId := meta.Get(TraceMetaKeyRpcId)
	if len(traceRpcId) == 0 {
		return false
	}

	traceRpcId = traceRpcId[:len(traceRpcId)-1] + ".0"
	meta.Set(TraceMetaKeyRpcId, traceRpcId)
	return true
}

// AppendKeyValue Ctx add k,v
func AppendKeyValue(ctx context.Context, key, value string) bool {
	meta := ExtractTraceNodeFromContext(ctx)
	traceRpcId := meta.Get(TraceMetaKeyRpcId)
	if len(traceRpcId) == 0 {
		return false
	}

	meta.Set(key, value)
	return true
}

// NewTraceMetadataContext create a new context with log trace metadata
func NewTraceMetadataContext(ctx context.Context, name string) context.Context {
	t := GenLogTraceMetadata()
	t.Set(TraceMetaKeyName, name)
	return context.WithValue(ctx, GetMetadataKey(), t)
}

// NewTraceNodeContext create a new context with trace node
func NewTraceNodeContext(ctx context.Context, meta *TraceNode) context.Context {
	return context.WithValue(ctx, GetMetadataKey(), meta)
}
