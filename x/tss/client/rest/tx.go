package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	clientUtils "github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/x/balance/exported"
	"github.com/axelarnetwork/axelar-core/x/tss/types"
)

type ReqKeygenStart struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	NewKeyId  string       `json:"key_id" yaml:"key_id"`
	Threshold int          `json:"threshold" yaml:"threshold"`
}

type ReqMasterkeyAssign struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	KeyId     string       `json:"key_id" yaml:"key_id"`
	Threshold int          `json:"threshold" yaml:"threshold"`
}

type ReqMasterkeyRotate struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/tx/%s/keygen/start", types.ModuleName), keygenStartHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/tx/%s/masterkey/assign/{chain}", types.ModuleName), masterkeyAssignHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/tx/%s/masterkey/rotate/{chain}", types.ModuleName), masterkeyRotateHandlerFn(cliCtx)).Methods("POST")
}

func keygenStartHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqKeygenStart
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		sender, ok := clientUtils.ExtractReqSender(w, req.BaseReq)
		if !ok {
			return
		}

		msg := types.MsgKeygenStart{
			Sender:    sender,
			NewKeyID:  req.NewKeyId,
			Threshold: req.Threshold,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func masterkeyAssignHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqMasterkeyAssign
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, ok := clientUtils.ExtractReqSender(w, req.BaseReq)
		if !ok {
			return
		}
		chain := mux.Vars(r)["chain"]

		msg := types.MsgAssignNextMasterKey{
			Sender: sender,
			Chain:  exported.ChainFromString(chain),
			KeyID:  req.KeyId,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func masterkeyRotateHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReqMasterkeyRotate
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		sender, ok := clientUtils.ExtractReqSender(w, req.BaseReq)
		if !ok {
			return
		}
		chain := mux.Vars(r)["chain"]

		msg := types.MsgRotateMasterKey{
			Sender: sender,
			Chain:  exported.ChainFromString(chain),
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}