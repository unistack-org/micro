//
// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"context"
	"crypto/tls"
	"net/http/httptrace"
	"net/textproto"
	"strings"
	"sync"

	"go.unistack.org/micro/v4/tracer"
)

const (
	httpStatus     = "http.status"
	httpHeaderMIME = "http.mime"
	httpRemoteAddr = "http.remote"
	httpLocalAddr  = "http.local"
	httpHost       = "http.host"
)

var hookMap = map[string]string{
	"http.dns":     "http.getconn",
	"http.connect": "http.getconn",
	"http.tls":     "http.getconn",
}

func parentHook(hook string) string {
	if strings.HasPrefix(hook, "http.connect") {
		return hookMap["http.connect"]
	}
	return hookMap[hook]
}

type clientTracer struct {
	context.Context
	tr          tracer.Tracer
	activeHooks map[string]context.Context
	root        tracer.Span
	mu          sync.Mutex
}

func NewClientTrace(ctx context.Context, tr tracer.Tracer) *httptrace.ClientTrace {
	ct := &clientTracer{
		Context:     ctx,
		activeHooks: make(map[string]context.Context),
		tr:          tr,
	}

	return &httptrace.ClientTrace{
		GetConn:              ct.getConn,
		GotConn:              ct.gotConn,
		PutIdleConn:          ct.putIdleConn,
		GotFirstResponseByte: ct.gotFirstResponseByte,
		Got100Continue:       ct.got100Continue,
		Got1xxResponse:       ct.got1xxResponse,
		DNSStart:             ct.dnsStart,
		DNSDone:              ct.dnsDone,
		ConnectStart:         ct.connectStart,
		ConnectDone:          ct.connectDone,
		TLSHandshakeStart:    ct.tlsHandshakeStart,
		TLSHandshakeDone:     ct.tlsHandshakeDone,
		WroteHeaderField:     ct.wroteHeaderField,
		WroteHeaders:         ct.wroteHeaders,
		Wait100Continue:      ct.wait100Continue,
		WroteRequest:         ct.wroteRequest,
	}
}

func (ct *clientTracer) start(hook, spanName string, attrs ...interface{}) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if hookCtx, found := ct.activeHooks[hook]; !found {
		var sp tracer.Span
		ct.activeHooks[hook], sp = ct.tr.Start(ct.getParentContext(hook), spanName, tracer.WithSpanLabels(attrs...), tracer.WithSpanKind(tracer.SpanKindClient))
		if ct.root == nil {
			ct.root = sp
		}
	} else {
		// end was called before start finished, add the start attributes and end the span here
		if span, ok := tracer.SpanFromContext(hookCtx); ok {
			span.AddLabels(attrs...)
			span.Finish()
		}

		delete(ct.activeHooks, hook)
	}
}

func (ct *clientTracer) end(hook string, err error, attrs ...interface{}) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if ctx, ok := ct.activeHooks[hook]; ok { // nolint:nestif
		if span, ok := tracer.SpanFromContext(ctx); ok {
			if err != nil {
				span.SetStatus(tracer.SpanStatusError, err.Error())
			}
			span.AddLabels(attrs...)
			span.Finish()
		}
		delete(ct.activeHooks, hook)
	} else {
		// start is not finished before end is called.
		// Start a span here with the ending attributes that will be finished when start finishes.
		// Yes, it's backwards. v0v
		ctx, span := ct.tr.Start(ct.getParentContext(hook), hook, tracer.WithSpanLabels(attrs...), tracer.WithSpanKind(tracer.SpanKindClient))
		if err != nil {
			span.SetStatus(tracer.SpanStatusError, err.Error())
		}
		ct.activeHooks[hook] = ctx
	}
}

func (ct *clientTracer) getParentContext(hook string) context.Context {
	ctx, ok := ct.activeHooks[parentHook(hook)]
	if !ok {
		return ct.Context
	}
	return ctx
}

func (ct *clientTracer) span(hook string) (tracer.Span, bool) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if ctx, ok := ct.activeHooks[hook]; ok {
		return tracer.SpanFromContext(ctx)
	}
	return nil, false
}

func (ct *clientTracer) getConn(host string) {
	ct.start("http.getconn", "http.getconn", httpHost, host)
}

func (ct *clientTracer) gotConn(info httptrace.GotConnInfo) {
	ct.end("http.getconn",
		nil,
		httpRemoteAddr, info.Conn.RemoteAddr().String(),
		httpLocalAddr, info.Conn.LocalAddr().String(),
	)
}

func (ct *clientTracer) putIdleConn(err error) {
	ct.end("http.receive", err)
}

func (ct *clientTracer) gotFirstResponseByte() {
	ct.start("http.receive", "http.receive")
}

func (ct *clientTracer) dnsStart(info httptrace.DNSStartInfo) {
	ct.start("http.dns", "http.dns", httpHost, info.Host)
}

func (ct *clientTracer) dnsDone(info httptrace.DNSDoneInfo) {
	ct.end("http.dns", info.Err)
}

func (ct *clientTracer) connectStart(network, addr string) {
	_ = network
	ct.start("http.connect."+addr, "http.connect", httpRemoteAddr, addr)
}

func (ct *clientTracer) connectDone(network, addr string, err error) {
	_ = network
	ct.end("http.connect."+addr, err)
}

func (ct *clientTracer) tlsHandshakeStart() {
	ct.start("http.tls", "http.tls")
}

func (ct *clientTracer) tlsHandshakeDone(_ tls.ConnectionState, err error) {
	ct.end("http.tls", err)
}

func (ct *clientTracer) wroteHeaderField(k string, v []string) {
	if sp, ok := ct.span("http.headers"); !ok || sp == nil {
		ct.start("http.headers", "http.headers")
	}
	ct.root.AddLabels("http."+strings.ToLower(k), sliceToString(v))
}

func (ct *clientTracer) wroteHeaders() {
	ct.start("http.send", "http.send")
}

func (ct *clientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		ct.root.SetStatus(tracer.SpanStatusError, info.Err.Error())
	}
	ct.end("http.send", info.Err)
}

func (ct *clientTracer) got100Continue() {
	if sp, ok := ct.span("http.receive"); ok && sp != nil {
		sp.AddEvent("GOT 100 - Continue")
	}
}

func (ct *clientTracer) wait100Continue() {
	if sp, ok := ct.span("http.receive"); ok && sp != nil {
		sp.AddEvent("GOT 100 - Wait")
	}
}

func (ct *clientTracer) got1xxResponse(code int, header textproto.MIMEHeader) error {
	if sp, ok := ct.span("http.receive"); ok && sp != nil {
		sp.AddEvent("GOT 1xx",
			tracer.WithEventLabels(
				httpStatus, code,
				httpHeaderMIME, sm2s(header),
			),
		)
	}
	return nil
}

func sliceToString(value []string) string {
	if len(value) == 0 {
		return "undefined"
	}
	return strings.Join(value, ",")
}

func sm2s(value map[string][]string) string {
	var buf strings.Builder
	for k, v := range value {
		if buf.Len() != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(sliceToString(v))
	}
	return buf.String()
}
