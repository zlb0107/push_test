package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Response struct {
	err    error
	code   int            // http status code, value from resp.StatusCode
	rsp    *http.Response // raw http.Response
	req    *http.Request  // raw http request
	buffer *bytes.Buffer  // rsp body buffer copy
}

func BuildResp(req *http.Request, resp *http.Response) (*Response, error) {
	res := &Response{
		rsp:    resp,
		req:    req,
		buffer: bytes.NewBuffer([]byte{}),
	}
	if resp != nil {
		res.code = resp.StatusCode
		if res.code != http.StatusOK {
			res.err = &url.Error{
				Err: fmt.Errorf("%s", resp.Status),
			}
		}
	} else {
		res.code = http.StatusInternalServerError
	}
	res.makeRspByteBuffer()
	return res, res.err
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) Code() int {
	return r.code
}

func (r *Response) RawRequest() *http.Request {
	return r.req
}

func (r *Response) RawResponse() *http.Response {
	return r.rsp
}

func (r *Response) ResetBody(body []byte) {
	r.buffer.Reset()
	r.rsp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

func (r *Response) Bytes() []byte {
	if r.err != nil {
		return nil
	}
	r.makeRspByteBuffer()
	return r.buffer.Bytes()
}

func (r *Response) String() string {
	if r.err != nil {
		return ""
	}
	r.makeRspByteBuffer()
	return r.buffer.String()
}

func (r *Response) JSON(obj interface{}) error {
	if r.err != nil {
		return r.err
	}
	r.makeRspByteBuffer()
	jsonDecoder := json.NewDecoder(r.buffer)
	return jsonDecoder.Decode(&obj)
}

func (r *Response) Save(fileName string) error {
	if r.err != nil {
		return r.err
	}

	fd, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	r.makeRspByteBuffer()
	_, err = io.Copy(fd, r.buffer)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (r *Response) ClearBuffer() {
	if r.err != nil {
		return
	}
	r.buffer.Reset()
}

func (r *Response) makeRspByteBuffer() {
	//reuse buffer
	if r.buffer.Len() != 0 || r.rsp == nil {
		return
	}

	_, err := io.Copy(r.buffer, r.rsp.Body)
	if err != nil {
		r.err = err
	}
	r.rsp.Body.Close()
}
