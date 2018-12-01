package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CallbackServer interface {
	HTML() string
	CSS() string
	Javascript() string
	Port() int
	Log() Logger
	Hangup(chan url.Values, url.Values)
	Start(chan url.Values)
}

type AuthCallbackServer struct {
	html       string
	css        string
	javascript string
	port       int
	log        Logger
	hangupFunc func(chan url.Values, url.Values)
}

func NewAuthCallbackServer(html, css, js string, log Logger, port int) AuthCallbackServer {
	acs := AuthCallbackServer{html: html, css: css, javascript: js, log: log, port: port}
	acs.SetHangupFunc(func(done chan url.Values, vals url.Values) {})
	return acs
}

func (acs AuthCallbackServer) HTML() string {
	return acs.html
}
func (acs AuthCallbackServer) CSS() string {
	return acs.css
}
func (acs AuthCallbackServer) Javascript() string {
	return acs.javascript
}
func (acs AuthCallbackServer) Port() int {
	return acs.port
}
func (acs AuthCallbackServer) Log() Logger {
	return acs.log
}
func (acs AuthCallbackServer) Hangup(done chan url.Values, values url.Values) {
	acs.hangupFunc(done, values)
}
func (acs *AuthCallbackServer) SetHangupFunc(hangupFunc func(chan url.Values, url.Values)) {
	acs.hangupFunc = hangupFunc
}

func (acs AuthCallbackServer) Start(done chan url.Values) {
	callbackValues := make(chan url.Values)
	serveMux := http.NewServeMux()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", acs.port),
		Handler: serveMux,
	}

	go func() {
		value := <-callbackValues
		close(callbackValues)
		srv.Close()
		done <- value
	}()

	attemptHangup := func(queryParams url.Values) {
		time.Sleep(10 * time.Millisecond)
		acs.Hangup(callbackValues, queryParams)
	}

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, acs.css)
		io.WriteString(w, acs.html)
		io.WriteString(w, acs.javascript)
		acs.log.Infof("Local server received request to %v %v", r.Method, r.RequestURI)

		// This is a goroutine because we want this handleFunc to complete before
		// Server.Close is invoked by listeners on the callbackValues channel.
		go attemptHangup(r.URL.Query())
	})

	go func() {
		acs.log.Infof("Starting local HTTP server on port %v", acs.port)
		acs.log.Info("Waiting for authorization redirect from Identity Provider...")
		if err := srv.ListenAndServe(); err != nil {
			acs.log.Infof("Stopping local HTTP server on port %v", acs.port)
		}
	}()
}
