package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/l-vitaly/imgprocessing/pkg/iohttp"

	"github.com/disintegration/imaging"
)

var (
	// ErrFailedLoadURL failed load image.
	ErrFailedLoadURL = errors.New("failed load image url")

	// ErrDataEmpty empty data image.
	ErrDataEmpty = errors.New("empty data image")
)

var _ Interface = &Service{}

// Service image processing.
type Service struct {
	savedPath string
}

func (s *Service) getImageHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Resize image.
func (s *Service) Resize(ctx context.Context, data []byte, width int, height int) (err error) {
	if len(data) == 0 {
		return ErrDataEmpty
	}

	bufImg := bytes.NewBuffer(data)

	srcImg, err := imaging.Decode(bufImg)
	if err != nil {
		return err
	}

	dstImage := imaging.Resize(srcImg, width, height, imaging.Lanczos)

	hashImg := s.getImageHash(data)

	err = imaging.Save(srcImg, s.savedPath+"/"+hashImg+"_original.jpg")
	if err != nil {
		return err
	}
	err = imaging.Save(dstImage, s.savedPath+"/"+hashImg+"_"+fmt.Sprintf("%dx%d", width, height)+".jpg")
	if err != nil {
		return err
	}
	return nil
}

// ResizeByURL resize image by URL.
func (s *Service) ResizeByURL(ctx context.Context, url string, width int, height int) (err error) {
	data, err := iohttp.GetContentByURL(url)
	if err != nil {
		return ErrFailedLoadURL
	}
	return s.Resize(ctx, data, width, height)
}

// NewService creates a service instance.
func NewService(savedPath string) *Service {
	return &Service{
		savedPath: savedPath,
	}
}
