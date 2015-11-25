package main

import (
	"net/http"
	"os"

	"github.com/dancannon/gorethink"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

var logger log.Logger

func main() {

	// setup a logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("svc", "sensor").With("caller", log.DefaultCaller)

	// setup rethinkdb session
	session := getSession()

	// build the context
	ctx := context.Background()

	// new service
	svc := NewSensorService(session)

	// wrap it in the logging middleware
	svc = loggingMiddleware(logger)(svc)

	// bind the service to HTTP with the context
	// with it's matching encoder/decoder
	recordHandler := httptransport.NewServer(
		ctx,
		makeRecordEndpoint(svc),
		decodeRecordRequest,
		encodeResponse,
	)

	// assign an endpoint route
	http.Handle("/sensor/record", recordHandler)

	// bind the listener
	logger.Log("msg", "HTTP", "addr", ":5000")
	logger.Log("err", http.ListenAndServe(":5000", nil))
}

func getSession() *gorethink.Session {
	session, _ := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "localhost:28015",
		Database: "sensord",
	})

	return session
}
