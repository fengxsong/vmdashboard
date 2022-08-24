package router

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/fengxsong/vmdashboard/pkg/router/chains"
)

func defaultBuildHandlerChain(handler http.Handler, o *completeOption, logger *zap.SugaredLogger) http.Handler {
	resolver := chains.NewRequestInfoResolver(APIGroupPrefix, []string{DefaultLegacyAPIPrefix}, logger)
	getter := chains.GetterFunc(chains.DefaultRequestAttributesGetter)

	// CanI authorize must been the innest filter chain
	handler = chains.WithCanI(handler, o.authorizer, getter, chains.Forbidden())
	handler = chains.WithRequestInfo(handler, resolver)
	handler = chains.WithAuthentication(handler, o.authenticator, chains.Unauthorized(), logger)
	handler = chains.WithLogging(handler, logger)
	return handler
}
