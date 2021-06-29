package tss

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/axelarnetwork/axelar-core/x/tss/keeper"
	"github.com/axelarnetwork/axelar-core/x/tss/types"
)

// NewHandler returns the handler for the tss module
func NewHandler(k keeper.Keeper, s types.Snapshotter, n types.Nexus, v types.Voter, staker types.StakingKeeper) sdk.Handler {
	server := keeper.NewMsgServerImpl(k, s, staker, v, n)
	h := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.ProcessKeygenTrafficRequest:
			res, err := server.ProcessKeygenTraffic(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.ProcessSignTrafficRequest:
			res, err := server.ProcessSignTraffic(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.StartKeygenRequest:
			res, err := server.StartKeygen(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.RotateKeyRequest:
			res, err := server.RotateKey(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.VotePubKeyRequest:
			res, err := server.VotePubKey(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.VoteSigRequest:
			res, err := server.VoteSig(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,
				fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg))
		}
	}

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		res, err := h(ctx, msg)
		if err != nil {
			k.Logger(ctx).Debug(err.Error())
			return nil, sdkerrors.Wrap(types.ErrTss, err.Error())
		}
		if res.Log != "" {
			k.Logger(ctx).Debug(res.Log)
		}
		return res, nil
	}
}
