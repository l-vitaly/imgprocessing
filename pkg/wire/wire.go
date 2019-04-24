//+build wireinject

package wire

import (
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	stdwire "github.com/google/wire"
	"github.com/l-vitaly/imgprocessing/pkg/service"
)

func SetupHTTPHandler(savedPath string, logger log.Logger, opts ...kithttp.ServerOption) (http.Handler, func(), error) {
	stdwire.Build(
		service.NewHTTPHandler,
		service.MakeEndpoints,
		service.MakeEncodeDecodeSet,

		ProvideService,
	)
	return nil, nil, nil
}

func ProvideService(
	savePath string,
	logger log.Logger,
) service.Interface {
	return service.NewLoggingService(
		service.NewService(savePath),
		logger,
	)
}
