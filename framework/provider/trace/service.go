package trace

import (
	"context"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"devops-http/framework/gin"
	"net/http"
	"time"
)

type KeyTrace string

var ContextKey = KeyTrace("trace-key")

type NiceTraceService struct {
	idService contract2.IDService

	traceIDGenerator contract2.IDService
	spanIDGenerator  contract2.IDService
}

func NewNiceTraceService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	idService := c.MustMake(contract2.IDKey).(contract2.IDService)
	return &NiceTraceService{idService: idService}, nil
}

// WithTrace register new trace to context
func (t *NiceTraceService) WithTrace(c context.Context, trace *contract2.TraceContext) context.Context {
	if ginC, ok := c.(*gin.Context); ok {
		ginC.Set(string(ContextKey), trace)
		return ginC
	} else {
		newC := context.WithValue(c, ContextKey, trace)
		return newC
	}
}

// GetTrace From trace context
func (t *NiceTraceService) GetTrace(c context.Context) *contract2.TraceContext {
	if ginC, ok := c.(*gin.Context); ok {
		if val, ok2 := ginC.Get(string(ContextKey)); ok2 {
			return val.(*contract2.TraceContext)
		}
	}

	if tc, ok := c.Value(ContextKey).(*contract2.TraceContext); ok {
		return tc
	}
	return nil
}

// NewTrace generate a new trace
func (t *NiceTraceService) NewTrace() *contract2.TraceContext {
	var traceID, spanID string
	if t.traceIDGenerator != nil {
		traceID = t.traceIDGenerator.NewID()
	} else {
		traceID = t.idService.NewID()
	}

	if t.spanIDGenerator != nil {
		spanID = t.spanIDGenerator.NewID()
	} else {
		spanID = t.idService.NewID()
	}
	tc := &contract2.TraceContext{
		TraceID:    traceID,
		ParentID:   "",
		SpanID:     spanID,
		CspanID:    "",
		Annotation: map[string]string{},
	}
	return tc
}

// StartSpan ChildSpan instance a sub trace with new span id
func (t *NiceTraceService) StartSpan(tc *contract2.TraceContext) *contract2.TraceContext {
	var childSpanID string
	if t.spanIDGenerator != nil {
		childSpanID = t.spanIDGenerator.NewID()
	} else {
		childSpanID = t.idService.NewID()
	}
	childSpan := &contract2.TraceContext{
		TraceID:  tc.TraceID,
		ParentID: "",
		SpanID:   tc.SpanID,
		CspanID:  childSpanID,
		Annotation: map[string]string{
			contract2.TraceKeyTime: time.Now().String(),
		},
	}
	return childSpan
}

// ExtractHTTP GetTrace By Http
func (t *NiceTraceService) ExtractHTTP(req *http.Request) *contract2.TraceContext {
	tc := &contract2.TraceContext{}
	tc.TraceID = req.Header.Get(contract2.TraceKeyTraceID)
	tc.ParentID = req.Header.Get(contract2.TraceKeySpanID)
	tc.SpanID = req.Header.Get(contract2.TraceKeyCspanID)
	tc.CspanID = ""

	if tc.TraceID == "" {
		tc.TraceID = t.idService.NewID()
	}

	if tc.SpanID == "" {
		tc.SpanID = t.idService.NewID()
	}

	return tc
}

// InjectHTTP Set Trace to Http
func (t *NiceTraceService) InjectHTTP(req *http.Request, tc *contract2.TraceContext) *http.Request {
	req.Header.Add(contract2.TraceKeyTraceID, tc.TraceID)
	req.Header.Add(contract2.TraceKeySpanID, tc.SpanID)
	req.Header.Add(contract2.TraceKeyCspanID, tc.CspanID)
	req.Header.Add(contract2.TraceKeyParentID, tc.ParentID)
	return req
}

func (t *NiceTraceService) ToMap(tc *contract2.TraceContext) map[string]string {
	m := map[string]string{}
	if tc == nil {
		return m
	}
	m[contract2.TraceKeyTraceID] = tc.TraceID
	m[contract2.TraceKeySpanID] = tc.SpanID
	m[contract2.TraceKeyCspanID] = tc.CspanID
	m[contract2.TraceKeyParentID] = tc.ParentID

	if tc.Annotation != nil {
		for k, v := range tc.Annotation {
			m[k] = v
		}
	}
	return m
}

// func (t *NiceTraceService) SetTraceIDService(service contract.IDService) {
// 	t.traceIDGenerator = service
// }

// func (t *NiceTraceService) SetSpanIDService(service contract.IDService) {
// 	t.spanIDGenerator = service
// }
