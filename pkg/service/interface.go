//go:generate gotool gk -o ./endpoint.go ep ./$GOFILE Interface
//go:generate gotool gk -o ./http.go tr ht -gresp -greq ./$GOFILE Interface
//go:generate gotool gk -o ./logging.go lg ./$GOFILE Interface
//go:generate gotool gk -o ./instrumenting.go it ./$GOFILE Interface
//go:generate gotool gk -o ./../../docs/swagger.yaml sg ./$GOFILE Interface

package service

import (
	"context"
)

type (
	// Interface image processing service.
	Interface interface {
		/*
			@GoKitEndpoint(
				Transport=@Transport(HTTP=@HTTP(
					Method='POST', Path='/resize', ResponseError=['400', '500'])
				)
			)
		*/
		Resize(ctx context.Context, data []byte, width, height int) (err error)
		/*
			@GoKitEndpoint(
				Transport=@Transport(HTTP=@HTTP(
					Method='POST', Path='/resize/url', ResponseError=['400', '500'])
				)
			)
		*/
		ResizeByURL(ctx context.Context, url string, width, height int) (err error)
	}
)
