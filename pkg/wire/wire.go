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
		stdwire.Bind(new(service.Interface), (*service.Service)(nil)),

		service.NewHTTPHandler,
		service.MakeEndpoints,
		service.MakeEncodeDecodeSet,

		service.NewService,
	)
	return nil, nil, nil
}
