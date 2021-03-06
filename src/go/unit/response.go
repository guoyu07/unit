/*
 * Copyright (C) Max Romanov
 * Copyright (C) NGINX, Inc.
 */

package unit

/*
#include "nxt_go_lib.h"
*/
import "C"

import (
	"fmt"
	"net/http"
	"os"
)

type response struct {
	header     http.Header
	headerSent bool
	req        *http.Request
	c_req      C.nxt_go_request_t
}

func new_response(c_req C.nxt_go_request_t, req *http.Request) *response {
	resp := &response{
		header: http.Header{},
		req:    req,
		c_req:  c_req,
	}

	return resp
}

func (r *response) Header() http.Header {
	return r.header
}

func (r *response) Write(p []byte) (n int, err error) {
	if !r.headerSent {
		r.WriteHeader(http.StatusOK)
	}

	l := C.size_t(len(p))
	b := getCBytes(p)
	res := C.nxt_go_response_write(r.c_req, b, l)
	C.free(b)
	return int(res), nil
}

func (r *response) WriteHeader(code int) {
	if r.headerSent {
		// Note: explicitly using Stderr, as Stdout is our HTTP output.
		fmt.Fprintf(os.Stderr, "CGI attempted to write header twice")
		return
	}
	r.headerSent = true
	fmt.Fprintf(r, "%s %d %s\r\n", r.req.Proto, code, http.StatusText(code))

	// Set a default Content-Type
	if _, hasType := r.header["Content-Type"]; !hasType {
		r.header.Add("Content-Type", "text/html; charset=utf-8")
	}

	r.header.Write(r)

	r.Write([]byte("\r\n"))
}
