package util

import (
	"io"
	"io/ioutil"
	"net/http"
)

func DrainResponse(r *http.Response) {
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
}
