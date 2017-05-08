package orbitip

import (
	"fmt"
	"net/http"
)

// Default path and file ext the reader sends requests too: /orbit.php
var (
	DefaultRoot = "/orbit"
	DefaultExt  = PHP
)

// Command defines NFC reader command types
type Command string

// String returns the command string
func (c Command) String() string {
	return string(c)
}

// NFC Commands
var (
	PU = Command("PU") // Power Up
	HB = Command("HB") // Heartbeat
	CO = Command("CO") // Card Opperation
	SW = Command("SW") // Level Change
	PG = Command("PG") // Ping
)

// Ext defines NFC path extenstion (.php, .asp etc)
type Ext struct {
	ID   uint8
	Name string
}

// Return the string of the extenstion
func (e Ext) String() string {
	return e.Name
}

// Supported extenstions by the reader
var (
	PHP  = Ext{0, ".php"} // Default
	ASP  = Ext{1, ".asp"}
	CFM  = Ext{2, ".cfm"}
	PL   = Ext{3, ".pl"}
	HTM  = Ext{4, ".htm"}
	HTML = Ext{5, ".html"}
	ASPX = Ext{6, ".aspx"}
	JSP  = Ext{7, ".jsp"}
)

// Params sent by NFC reader
type Params struct {
	Date     string
	Time     string
	ID       string
	ULen     string
	UID      string
	Command  string
	Version  string
	Contact1 string
	Contact2 string
	SID      string
	Data     string
	PSRC     string
	MD5      string
	MAC      string
	Relay    string
	SD       string
}

// The HandlerFunc type allows the use of ordinary functions as NFC command handlers.
type HandlerFunc func(parms Params) ([]byte, error)

// Handlers maps a NFC command to a specific HandlerFunc.
type Handlers map[Command]HandlerFunc

// Set sets a HandlerFunc for the given NFC command
func (h Handlers) Set(cmd Command, fn HandlerFunc) {
	h[cmd] = fn
}

// Del removes a HandleFunc for the given NFC command
func (h Handlers) Del(cmd Command) {
	delete(h, cmd)
}

// The ServeMux type implements the http.Handler interface for serving
// NFC requests. This can be passed directly into a http.Server instance on the
// Handler property.
type ServeMux struct {
	handlers Handlers
}

// Handlers returns the NFC Command to HandlerFunc map
func (m *ServeMux) Handlers() Handlers {
	return m.handlers
}

// ServeHTTP handles an NFC HTTP request, first checking if a HandleFunc
// exists for the NFC command and if so calling that function. Any data returned
// by the HandleFunc is written to the http response.
func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fn, ok := m.handlers[Command(values.Get("cmd"))]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := fn(Params{
		Date:     values.Get("date"),
		Time:     values.Get("time"),
		ID:       values.Get("id"),
		ULen:     values.Get("ulen"),
		UID:      values.Get("uid"),
		Command:  values.Get("cmd"),
		Version:  values.Get("ver"),
		Contact1: values.Get("contact1"),
		Contact2: values.Get("contact2"),
		SID:      values.Get("sid"),
		Data:     values.Get("data"),
		PSRC:     values.Get("psrc"),
		MD5:      values.Get("md5"),
		MAC:      values.Get("mac"),
		Relay:    values.Get("relay"),
		SD:       values.Get("sd"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

// NewServeMux constructs a new ServeMux
func NewServeMux(handlers Handlers) *ServeMux {
	return &ServeMux{
		handlers: handlers,
	}
}

// New constructs a new HTTP server for running the NFC HTTP server
func New(addr string, root string, ext Ext, handlers Handlers) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("%s%s", root, ext), NewServeMux(handlers))
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
