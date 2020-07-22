package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

func RegisterHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/registered/{source_port}/{source_channel}/{account}", queryRequestHandlerFn(clientCtx)).Methods("GET")
}
