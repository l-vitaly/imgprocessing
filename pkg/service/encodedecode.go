package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// MakeEncodeDecodeSet create encode/decode sets.
func MakeEncodeDecodeSet() (set EncodeDecodeSet) {
	set.Server.DecodeHTTPResizeRequest = func(_ context.Context, r *http.Request) (request interface{}, err error) {
		var req ResizeRequest

		err = r.ParseMultipartForm(1024 * 1024 * 10)

		if err != nil && r.Body != nil {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return nil, err
			}
			return req, nil
		}
		if err != nil {
			return nil, err
		}

		if r.MultipartForm != nil {
			formFiles := r.MultipartForm.File["data"]
			if len(formFiles) > 0 {
				f, err := formFiles[0].Open()
				if err != nil {
					return nil, err
				}
				defer f.Close()

				req.Data, err = ioutil.ReadAll(f)

				if err != nil {
					return nil, err
				}
			}
		}
		vars := mux.Vars(r)
		if req.Width, err = strconv.Atoi(vars["width"]); err != nil {
			return nil, err
		}
		if req.Height, err = strconv.Atoi(vars["height"]); err != nil {
			return nil, err
		}
		return req, nil
	}
	set.Server.DecodeHTTPResizeByURLRequest = func(_ context.Context, r *http.Request) (request interface{}, err error) {
		var req ResizeByURLRequest

		vars := mux.Vars(r)

		req.URL = vars["url"]

		if req.Width, err = strconv.Atoi(vars["width"]); err != nil {
			return nil, err
		}
		if req.Height, err = strconv.Atoi(vars["height"]); err != nil {
			return nil, err
		}
		return req, nil
	}
	return
}
