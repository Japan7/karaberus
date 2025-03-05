package server

import (
	"fmt"
	"io"
	"net/http"
)

// Close but you can’t have an error returned so you can safely defer it
func Closer(closer io.Closer) {
	if closer == nil {
		return
	}
	err := closer.Close()
	if err != nil {
		Warn(err.Error())
	}
}

// http.Client.Do but you can’t have a nil response
func Do(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("request failed terribly for %s", req.URL)
	}

	return resp, err
}
