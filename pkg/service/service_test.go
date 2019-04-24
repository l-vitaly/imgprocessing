package service_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/l-vitaly/imgprocessing/pkg/iohttp"
	"github.com/l-vitaly/imgprocessing/pkg/service"
	. "github.com/smartystreets/goconvey/convey"
)

func createServer(savedPath string) http.Handler {
	nopLogger := log.NewNopLogger()
	edSet := service.MakeEncodeDecodeSet()
	s := service.NewLoggingService(
		service.NewService(savedPath),
		nopLogger,
	)
	set := service.MakeEndpoints(s)
	return service.NewHTTPHandler(s, set, edSet, nopLogger)
}

func TestHTTP_Resize(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	savedPath := wd + "/fixtures/.tmp"

	err = os.MkdirAll(savedPath, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(savedPath)

	h := createServer(savedPath)

	filePath := wd + "/fixtures/intro.jpg"

	imageData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	Convey("Send request for resize image", t, func() {
		rr := httptest.NewRecorder()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("data", filepath.Base(filePath))

		So(err, ShouldBeNil)

		_, err = io.Copy(part, bytes.NewBuffer(imageData))

		So(err, ShouldBeNil)

		err = writer.Close()

		So(err, ShouldBeNil)

		v := url.Values{}
		v.Add("width", "100")
		v.Add("height", "100")

		req, err := http.NewRequest("POST", "/resize", body)

		So(err, ShouldBeNil)

		req.Header.Set("Content-Type", writer.FormDataContentType())

		req.URL.RawQuery = v.Encode()

		h.ServeHTTP(rr, req)

		expected := `{}`

		So(strings.TrimSpace(rr.Body.String()), ShouldEqual, expected)

		So(rr.Code, ShouldEqual, 200)

		imageHash := getImageHash(imageData)

		dstImage, err := os.Open(savedPath + "/" + fmt.Sprintf("%s_%dx%d.jpg", imageHash, 100, 100))
		defer dstImage.Close()

		So(err, ShouldBeNil)

		dstImageData, err := ioutil.ReadAll(dstImage)

		So(err, ShouldBeNil)

		cfg, _, err := image.DecodeConfig(bytes.NewBuffer(dstImageData))

		So(err, ShouldBeNil)

		So(cfg.Width, ShouldEqual, 100)
		So(cfg.Height, ShouldEqual, 100)
	})

	Convey("Send fail request without multipart form-data ", t, func() {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/resize", nil)

		v := url.Values{}
		v.Add("width", "100")
		v.Add("height", "100")

		req.URL.RawQuery = v.Encode()

		So(err, ShouldBeNil)

		h.ServeHTTP(rr, req)

		expected := `{"error":"missing form body"}`

		So(strings.TrimSpace(rr.Body.String()), ShouldEqual, expected)
		So(rr.Code, ShouldEqual, 500)
	})

	Convey("Send fail request with empty multipart form-data", t, func() {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/resize", nil)

		v := url.Values{}
		v.Add("width", "100")
		v.Add("height", "100")

		req.URL.RawQuery = v.Encode()

		So(err, ShouldBeNil)

		h.ServeHTTP(rr, req)

		expected := `{"error":"missing form body"}`

		So(strings.TrimSpace(rr.Body.String()), ShouldEqual, expected)
		So(rr.Code, ShouldEqual, 500)
	})
}

func TestHTTP_ResizeByURL(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	savedPath := wd + "/fixtures/.tmp"

	err = os.MkdirAll(savedPath, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(savedPath)

	h := createServer(savedPath)

	rr := httptest.NewRecorder()

	urlStr := "http://agilefusion.com/build/images/pages/aboutus/intro.jpg"

	imageData, err := iohttp.GetContentByURL(urlStr)
	if err != nil {
		t.Fatal(err)
	}

	Convey("Send request for resize by URL", t, func() {

		v := url.Values{}
		v.Add("url", urlStr)
		v.Add("width", "100")
		v.Add("height", "100")

		req, err := http.NewRequest("GET", "/resize", nil)

		So(err, ShouldBeNil)

		req.URL.RawQuery = v.Encode()

		h.ServeHTTP(rr, req)

		expected := `{}`

		So(strings.TrimSpace(rr.Body.String()), ShouldEqual, expected)

		So(rr.Code, ShouldEqual, 200)

		imageHash := getImageHash(imageData)

		dstImage, err := os.Open(savedPath + "/" + fmt.Sprintf("%s_%dx%d.jpg", imageHash, 100, 100))
		defer dstImage.Close()

		So(err, ShouldBeNil)

		dstImageData, err := ioutil.ReadAll(dstImage)

		So(err, ShouldBeNil)

		cfg, _, err := image.DecodeConfig(bytes.NewBuffer(dstImageData))

		So(err, ShouldBeNil)

		So(cfg.Width, ShouldEqual, 100)
		So(cfg.Height, ShouldEqual, 100)
	})
}

func getImageHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
