package server

import (
	"bytes"
	"net/http"

	"git.inke.cn/inkelogic/daenerys/config/encoder/json"
)

type Responser interface {
	http.ResponseWriter
	Size() int
	Status() int
	Writer() http.ResponseWriter
	WriteString(string) (int, error)
	WriteJSON(interface{}) (int, error)
	ByteBody() []byte
	StringBody() string
	DoFlush()
}

const (
	noWritten = -1
)

var _ Responser = &responseWriter{}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
	buff   bytes.Buffer
}

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = http.StatusOK
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Writer() http.ResponseWriter {
	return w
}

func (w *responseWriter) ByteBody() []byte {
	return w.buff.Bytes()
}

func (w *responseWriter) StringBody() string {
	return w.buff.String()
}

func (w *responseWriter) written() bool {
	return w.size != noWritten
}

//header status code should write only once
func (w *responseWriter) writeHeaderOnce() {
	if !w.written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

//override
func (w *responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

//override
func (w *responseWriter) Write(data []byte) (n int, err error) {
	w.writeHeaderOnce()
	n, err = w.buff.Write(data)
	w.size += n
	return
}

//override
func (w *responseWriter) WriteHeader(statusCode int) {
	if statusCode > 0 && w.status != statusCode {
		w.status = statusCode
	}
}

func (w *responseWriter) WriteString(s string) (n int, err error) {
	w.writeHeaderOnce()
	n, err = w.buff.WriteString(s)
	w.size += n
	return
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (w *responseWriter) WriteJSON(data interface{}) (n int, err error) {
	b, err := json.NewEncoder().Encode(data)
	if err != nil {
		return
	}

	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentType
	}

	w.writeHeaderOnce()
	n, err = w.buff.Write(b)
	w.size += n
	return
}

func (w *responseWriter) DoFlush() {
	w.writeHeaderOnce()
	if b := w.buff.Bytes(); len(b) > 0 {
		w.ResponseWriter.Write(b)
		w.buff.Reset()
	}
}
