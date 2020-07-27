package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
)

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/registered/{source_port}/{source_channel}/{account}", QueryRegisteredRequestHandlerFn(cliCtx)).Methods("GET")
}
