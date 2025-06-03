package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
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

func extendETag(last_item_id uint, err error, etag *uint) error {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			*etag = 0
			err = nil
		} else {
			return err
		}
	} else {
		*etag += last_item_id
	}

	return nil
}

func setETag(last_item_id uint, err error, etag *string) error {
	etag_n := uint(0)
	err = extendETag(last_item_id, err, &etag_n)
	if err != nil {
		return err
	}

	*etag = fmt.Sprint(etag_n)
	return nil
}
