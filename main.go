package main

import (
	"fmt"
	"os"
	"flag"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

const (
	serviceName = "opentracing-proxy"
	hostPort = "0.0.0.0:0"
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"
	debug = false
	sameSpan = true
	traceID128Bit = true
)

func main() {

	// ********************************************************
	// Here are some of the necessary components that an application reporting to Zipkin will need
	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v", err)
		os.Exit(-1)
	}
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v", err)
		os.Exit(-1)
	}
	// We will use this as the way to tracking traces between requests and responses
	cache := make(map[int64]opentracing.Span)
	// ********************************************************

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.Logf("%v", "We can see what APIs are being called!")

		// ********************************************************
		// Here is where we create the start of our "span"
		// You can see that there isn't anything Zipkin-specific below
		span := tracer.StartSpan("GotRequest")
		cache[ctx.Session] = span
		tracer.Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
		ctx.Logf("%v", req.Header)
		// This tag is viewable by clicking on the trace *and* clicking
		// on the span in the trace
		span.SetTag("Host", req.Host)
		span.LogEvent("Injection")
		// ********************************************************

		return req, ctx.Resp
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

		// ********************************************************
		// Finally, we can finish our "span" from earlier
		span := cache[ctx.Session]
		if span != nil {
			span.LogEvent("Finished")
			defer span.Finish()
		}
		// ********************************************************

		ctx.Logf("%v", "We can modify some data coming back!")
		return resp
	})
	verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
