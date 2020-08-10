package util

import (
	"io"
	"io/ioutil"
	"net/http"
)

func DrainResponse(r *http.Response) {
	if r == nil || r.Body == nil {
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
}
