package keeper

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	gogoprototypes "github.com/gogo/protobuf/types"

	"github.com/axelarnetwork/axelar-core/utils"
	"github.com/axelarnetwork/axelar-core/x/evm/exported"
	"github.com/axelarnetwork/axelar-core/x/evm/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	tss "github.com/axelarnetwork/axelar-core/x/tss/exported"
	tsstypes "github.com/axelarnetwork/axelar-core/x/tss/types"
	vote "github.com/axelarnetwork/axelar-core/x/vote/exported"
)

var _ types.MsgServiceServer = msgServer{}

type msgServer struct {
	types.BaseKeeper
	tss         types.TSS
	signer      types.Signer
	nexus       types.Nexus
	voter       types.Voter
	snapshotter types.Snapshotter
}

// NewMsgServerImpl returns an implementation of the bitcoin MsgServiceServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.BaseKeeper, t types.TSS, n types.Nexus, s types.Signer, v types.Voter, snap types.Snapshotter) types.MsgServiceServer {
	return msgServer{
		BaseKeeper:  keeper,
		tss:         t,
		signer:      s,
		nexus:       n,
		voter:       v,
		snapshotter: snap,
	}
}

func validateChainActivated(ctx sdk.Context, n types.Nexus, chain nexus.Chain) error {
	if !n.IsChainActivated(ctx, chain) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("chain %s is not activated yet", chain.Name))
	}

	return nil
}

func (s msgServer) ConfirmGatewayDeployment(c context.Context, req *types.ConfirmGatewayDeploymentRequest) (*types.ConfirmGatewayDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(chain.Name)

	if _, ok := keeper.GetPendingGatewayAddress(ctx); ok {
		return nil, fmt.Errorf("gateway is in the process of confirmation")
	}

	if _, ok := keeper.GetGatewayAddress(ctx); ok {
		return nil, fmt.Errorf("gateway is already confirmed")
	}

	keeper.SetPendingGateway(ctx, common.Address(req.Address))

	period, ok := keeper.GetRevoteLockingPeriod(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve revote locking period")
	}

	votingThreshold, ok := keeper.GetVotingThreshold(ctx)
	if !ok {
		return nil, fmt.Errorf("voting threshold not found")
	}

	minVoterCount, ok := keeper.GetMinVoterCount(ctx)
	if !ok {
		return nil, fmt.Errorf("min voter count not found")
	}

	pollKey := types.GetConfirmGatewayDeploymentPollKey(chain, req.TxID, req.Address)
	if err := s.voter.InitializePoll(
		ctx,
		pollKey,
		s.nexus.GetChainMaintainers(ctx, chain),
		vote.ExpiryAt(ctx.BlockHeight()+period),
		vote.Threshold(votingThreshold),
		vote.MinVoterCount(minVoterCount),
		vote.RewardPool(chain.Name),
	); err != nil {
		return nil, err
	}

	deploymentBytecode, err := getGatewayDeploymentBytecode(ctx, keeper, s.signer, chain)
	if err != nil {
		return nil, err
	}

	height, _ := keeper.GetRequiredConfirmationHeight(ctx)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeGatewayDeploymentConfirmation,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueStart),
			sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
			sdk.NewAttribute(types.AttributeKeyTxID, req.TxID.Hex()),
			sdk.NewAttribute(types.AttributeKeyAddress, req.Address.Hex()),
			sdk.NewAttribute(types.AttributeKeyBytecodeHash, hex.EncodeToString(crypto.Keccak256(deploymentBytecode))),
			sdk.NewAttribute(types.AttributeKeyConfHeight, strconv.FormatUint(height, 10)),
			sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&pollKey))),
		),
	)

	return nil, nil
}

func (s msgServer) VoteConfirmGatewayDeployment(c context.Context, req *types.VoteConfirmGatewayDeploymentRequest) (*types.VoteConfirmGatewayDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(chain.Name)

	if _, ok := keeper.GetGatewayAddress(ctx); ok {
		return &types.VoteConfirmGatewayDeploymentResponse{Log: "gateway is already confirmed"}, nil
	}

	address, ok := keeper.GetPendingGatewayAddress(ctx)
	if !ok {
		return nil, fmt.Errorf("no pending gateway found")
	}

	voter := s.snapshotter.GetOperator(ctx, req.Sender)
	if voter == nil {
		return nil, fmt.Errorf("account %v is not registered as a validator proxy", req.Sender.String())
	}

	poll := s.voter.GetPoll(ctx, req.PollKey)
	voteValue := &gogoprototypes.BoolValue{Value: req.Confirmed}
	if err := poll.Vote(voter, voteValue); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGatewayDeploymentConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueVote),
		sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatBool(voteValue.Value)),
	))

	if poll.Is(vote.Pending) {
		return &types.VoteConfirmGatewayDeploymentResponse{Log: fmt.Sprintf("not enough votes to confirm gateway for chain %s yet", chain.Name)}, nil
	}

	if poll.Is(vote.Failed) {
		if err := keeper.DeletePendingGateway(ctx); err != nil {
			return nil, err
		}

		return &types.VoteConfirmGatewayDeploymentResponse{Log: fmt.Sprintf("poll %s failed", poll.GetKey())}, nil
	}

	confirmed, ok := poll.GetResult().(*gogoprototypes.BoolValue)
	if !ok {
		return nil, fmt.Errorf("result of poll %s has wrong type, expected bool, got %T", req.PollKey.String(), poll.GetResult())
	}

	s.Logger(ctx).Info(fmt.Sprintf("%s gateway confirmation result is %t", chain.Name, confirmed.Value))

	// handle poll result
	event := sdk.NewEvent(
		types.EventTypeGatewayDeploymentConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
		sdk.NewAttribute(types.AttributeKeyAddress, address.Hex()),
		sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&req.PollKey))))
	defer func() { ctx.EventManager().EmitEvent(event) }()

	if !confirmed.Value {
		poll.AllowOverride()
		event = event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueReject))

		if err := keeper.DeletePendingGateway(ctx); err != nil {
			return nil, err
		}

		return &types.VoteConfirmGatewayDeploymentResponse{
			Log: fmt.Sprintf("%s gateway was discarded", chain.Name),
		}, nil
	}

	event = event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueConfirm))
	keeper.ConfirmPendingGateway(ctx)

	return &types.VoteConfirmGatewayDeploymentResponse{}, nil
}

func (s msgServer) Link(c context.Context, req *types.LinkRequest) (*types.LinkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	senderChain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, senderChain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(senderChain.Name)
	gatewayAddr, ok := keeper.GetGatewayAddress(ctx)
	if !ok {
		return nil, fmt.Errorf("axelar gateway address not set")
	}

	recipientChain, ok := s.nexus.GetChain(ctx, req.RecipientChain)
	if !ok {
		return nil, fmt.Errorf("unknown recipient chain")
	}

	token := keeper.GetERC20TokenByAsset(ctx, req.Asset)
	found := s.nexus.IsAssetRegistered(ctx, recipientChain.Name, req.Asset)
	if !found || !token.Is(types.Confirmed) {
		return nil, fmt.Errorf("asset '%s' not registered for chain '%s'", req.Asset, recipientChain.Name)
	}

	tokenAddr := token.GetAddress()

	burnerAddr, salt, err := keeper.GetBurnerAddressAndSalt(ctx, tokenAddr, req.RecipientAddr, gatewayAddr)
	if err != nil {
		return nil, err
	}

	symbol := token.GetDetails().Symbol

	s.nexus.LinkAddresses(ctx,
		nexus.CrossChainAddress{Chain: senderChain, Address: burnerAddr.String()},
		nexus.CrossChainAddress{Chain: recipientChain, Address: req.RecipientAddr})

	burnerInfo := types.BurnerInfo{
		TokenAddress:     types.Address(tokenAddr),
		DestinationChain: req.RecipientChain,
		Symbol:           symbol,
		Asset:            req.Asset,
		Salt:             types.Hash(salt),
	}
	keeper.SetBurnerInfo(ctx, burnerAddr, &burnerInfo)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLink,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyChain, req.Chain),
			sdk.NewAttribute(types.AttributeKeyBurnAddress, burnerAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, req.RecipientAddr),
			sdk.NewAttribute(types.AttributeKeyDestinationChain, req.RecipientChain),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, tokenAddr.Hex()),
		),
	)

	return &types.LinkResponse{DepositAddr: burnerAddr.Hex()}, nil
}

// ConfirmToken handles token deployment confirmation
func (s msgServer) ConfirmToken(c context.Context, req *types.ConfirmTokenRequest) (*types.ConfirmTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	_, ok = s.nexus.GetChain(ctx, req.Asset.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Asset.Chain)
	}

	keeper := s.ForChain(chain.Name)
	token := keeper.GetERC20TokenByAsset(ctx, req.Asset.Name)

	err := token.RecordDeployment(req.TxID)
	if err != nil {
		return nil, err
	}

	period, ok := keeper.GetRevoteLockingPeriod(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve revote locking period")
	}

	votingThreshold, ok := keeper.GetVotingThreshold(ctx)
	if !ok {
		return nil, fmt.Errorf("voting threshold not found")
	}

	minVoterCount, ok := keeper.GetMinVoterCount(ctx)
	if !ok {
		return nil, fmt.Errorf("min voter count not found")
	}

	pollKey := types.GetConfirmTokenKey(req.TxID, req.Asset.Name)
	if err := s.voter.InitializePoll(
		ctx,
		pollKey,
		s.nexus.GetChainMaintainers(ctx, chain),
		vote.ExpiryAt(ctx.BlockHeight()+period),
		vote.Threshold(votingThreshold),
		vote.MinVoterCount(minVoterCount),
		vote.RewardPool(chain.Name),
	); err != nil {
		return nil, err
	}

	// if token was initialized, both token and gateway addresses are available
	tokenAddr := token.GetAddress()
	gatewayAddr, _ := keeper.GetGatewayAddress(ctx)
	height, _ := keeper.GetRequiredConfirmationHeight(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeTokenConfirmation,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueStart),
			sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
			sdk.NewAttribute(types.AttributeKeyTxID, req.TxID.Hex()),
			sdk.NewAttribute(types.AttributeKeyGatewayAddress, gatewayAddr.Hex()),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, tokenAddr.Hex()),
			sdk.NewAttribute(types.AttributeKeySymbol, token.GetDetails().Symbol),
			sdk.NewAttribute(types.AttributeKeyAsset, req.Asset.Name),
			sdk.NewAttribute(types.AttributeKeyConfHeight, strconv.FormatUint(height, 10)),
			sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&pollKey))),
		),
	)

	return &types.ConfirmTokenResponse{}, nil
}

func (s msgServer) ConfirmChain(c context.Context, req *types.ConfirmChainRequest) (*types.ConfirmChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if _, found := s.nexus.GetChain(ctx, req.Name); found {
		return nil, fmt.Errorf("chain '%s' is already confirmed", req.Name)
	}

	if _, ok := s.GetPendingChain(ctx, req.Name); !ok {
		return nil, fmt.Errorf("'%s' has not been added yet", req.Name)
	}

	seqNo := s.snapshotter.GetLatestCounter(ctx)
	if seqNo < 0 {
		keyRequirement, ok := s.tss.GetKeyRequirement(ctx, tss.MasterKey, exported.Ethereum.KeyType)
		if !ok {
			return nil, fmt.Errorf("key requirement for key role %s type %s not found", tss.MasterKey.SimpleString(), exported.Ethereum.KeyType)
		}

		snapshot, err := s.snapshotter.TakeSnapshot(ctx, keyRequirement)
		if err != nil {
			return nil, fmt.Errorf("unable to take snapshot: %v", err)
		}

		seqNo = snapshot.Counter
	}
	keeper := s.ForChain(req.Name)

	period, ok := keeper.GetRevoteLockingPeriod(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve revote locking period for chain %s", req.Name)
	}

	votingThreshold, ok := keeper.GetVotingThreshold(ctx)
	if !ok {
		return nil, fmt.Errorf("voting threshold for chain %s not found", req.Name)
	}

	minVoterCount, ok := keeper.GetMinVoterCount(ctx)
	if !ok {
		return nil, fmt.Errorf("min voter count for chain %s not found", req.Name)
	}

	pollKey := vote.NewPollKey(types.ModuleName, req.Name)
	if err := s.voter.InitializePollWithSnapshot(
		ctx,
		pollKey,
		seqNo,
		vote.ExpiryAt(ctx.BlockHeight()+period),
		vote.Threshold(votingThreshold),
		vote.MinVoterCount(minVoterCount),
	); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeChainConfirmation,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueStart),
			sdk.NewAttribute(types.AttributeKeyChain, req.Name),
			sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&pollKey))),
		),
	)

	return &types.ConfirmChainResponse{}, nil
}

// ConfirmDeposit handles deposit confirmations
func (s msgServer) ConfirmDeposit(c context.Context, req *types.ConfirmDepositRequest) (*types.ConfirmDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(chain.Name)

	_, state, ok := keeper.GetDeposit(ctx, common.Hash(req.TxID), common.Address(req.BurnerAddress))
	switch {
	case !ok:
		break
	case state == types.CONFIRMED:
		return nil, fmt.Errorf("already confirmed")
	case state == types.BURNED:
		return nil, fmt.Errorf("already burned")
	}

	burnerInfo := keeper.GetBurnerInfo(ctx, common.Address(req.BurnerAddress))
	if burnerInfo == nil {
		return nil, fmt.Errorf("no burner info found for address %s", req.BurnerAddress)
	}

	period, ok := keeper.GetRevoteLockingPeriod(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve revote locking period for chain %s", req.Chain)
	}

	votingThreshold, ok := keeper.GetVotingThreshold(ctx)
	if !ok {
		return nil, fmt.Errorf("voting threshold for chain %s not found", chain.Name)
	}

	minVoterCount, ok := keeper.GetMinVoterCount(ctx)
	if !ok {
		return nil, fmt.Errorf("min voter count for chain %s not found", chain.Name)
	}

	pollKey := vote.NewPollKey(types.ModuleName, fmt.Sprintf("%s_%s_%d", req.TxID.Hex(), req.BurnerAddress.Hex(), req.Amount.Uint64()))
	if err := s.voter.InitializePoll(
		ctx,
		pollKey,
		s.nexus.GetChainMaintainers(ctx, chain),
		vote.ExpiryAt(ctx.BlockHeight()+period),
		vote.Threshold(votingThreshold),
		vote.MinVoterCount(minVoterCount),
		vote.RewardPool(chain.Name),
	); err != nil {
		return nil, err
	}

	erc20Deposit := types.ERC20Deposit{
		TxID:             req.TxID,
		Amount:           req.Amount,
		Asset:            burnerInfo.Asset,
		DestinationChain: burnerInfo.DestinationChain,
		BurnerAddress:    req.BurnerAddress,
	}
	keeper.SetPendingDeposit(ctx, pollKey, &erc20Deposit)

	height, _ := keeper.GetRequiredConfirmationHeight(ctx)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeDepositConfirmation,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueStart),
			sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
			sdk.NewAttribute(types.AttributeKeyTxID, req.TxID.Hex()),
			sdk.NewAttribute(types.AttributeKeyAmount, req.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyBurnAddress, req.BurnerAddress.Hex()),
			sdk.NewAttribute(types.AttributeKeyTokenAddress, burnerInfo.TokenAddress.Hex()),
			sdk.NewAttribute(types.AttributeKeyConfHeight, strconv.FormatUint(height, 10)),
			sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&pollKey))),
		),
	)

	return &types.ConfirmDepositResponse{}, nil
}

// ConfirmTransferKey handles transfer ownership/operatorship confirmations
func (s msgServer) ConfirmTransferKey(c context.Context, req *types.ConfirmTransferKeyRequest) (*types.ConfirmTransferKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	var keyRole tss.KeyRole
	switch req.TransferType {
	case types.Ownership:
		keyRole = tss.MasterKey
	case types.Operatorship:
		keyRole = tss.SecondaryKey
	default:
		return nil, fmt.Errorf("invalid transfer type %s", req.TransferType.SimpleString())
	}

	_, ok = s.signer.GetNextKeyID(ctx, chain, keyRole)
	if !ok {
		return nil, fmt.Errorf("next %s key for chain %s not set yet", keyRole.SimpleString(), chain.Name)
	}

	keeper := s.ForChain(chain.Name)

	gatewayAddr, ok := keeper.GetGatewayAddress(ctx)
	if !ok {
		return nil, fmt.Errorf("axelar gateway address not set")
	}

	period, ok := keeper.GetRevoteLockingPeriod(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve revote locking period for chain %s", req.Chain)
	}

	votingThreshold, ok := keeper.GetVotingThreshold(ctx)
	if !ok {
		return nil, fmt.Errorf("voting threshold for chain %s not found", chain.Name)
	}

	minVoterCount, ok := keeper.GetMinVoterCount(ctx)
	if !ok {
		return nil, fmt.Errorf("min voter count for chain %s not found", chain.Name)
	}

	pollKey := vote.NewPollKey(types.ModuleName, fmt.Sprintf("%s_%s_%s", req.TxID.Hex(), req.TransferType.SimpleString(), req.KeyID))
	if err := s.voter.InitializePoll(
		ctx,
		pollKey,
		s.nexus.GetChainMaintainers(ctx, chain),
		vote.ExpiryAt(ctx.BlockHeight()+period),
		vote.Threshold(votingThreshold),
		vote.MinVoterCount(minVoterCount),
		vote.RewardPool(chain.Name),
	); err != nil {
		return nil, err
	}

	transferKey := types.TransferKey{
		TxID:      req.TxID,
		Type:      req.TransferType,
		NextKeyID: req.KeyID,
	}
	keeper.SetPendingTransferKey(ctx, pollKey, &transferKey)

	height, _ := keeper.GetRequiredConfirmationHeight(ctx)

	event := sdk.NewEvent(types.EventTypeTransferKeyConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueStart),
		sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
		sdk.NewAttribute(types.AttributeKeyTxID, req.TxID.Hex()),
		sdk.NewAttribute(types.AttributeKeyTransferKeyType, req.TransferType.SimpleString()),
		sdk.NewAttribute(types.AttributeKeyKeyType, chain.KeyType.SimpleString()),
		sdk.NewAttribute(types.AttributeKeyGatewayAddress, gatewayAddr.Hex()),
		sdk.NewAttribute(types.AttributeKeyConfHeight, strconv.FormatUint(height, 10)),
		sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&pollKey))),
	)
	defer func() { ctx.EventManager().EmitEvent(event) }()

	key, ok := s.signer.GetKey(ctx, req.KeyID)
	if !ok {
		return nil, fmt.Errorf("key %s does not exist", req.KeyID)
	}

	switch chain.KeyType {
	case tss.Threshold:
		pk, err := key.GetECDSAPubKey()
		if err != nil {
			return nil, err
		}

		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyAddress, crypto.PubkeyToAddress(pk).Hex()),
			sdk.NewAttribute(types.AttributeKeyThreshold, ""),
		)
	case tss.Multisig:
		addresses, threshold, err := getMultisigAddresses(key)
		if err != nil {
			return nil, err
		}

		addressStrs := make([]string, len(addresses))
		for i, address := range addresses {
			addressStrs[i] = address.Hex()
		}

		event = event.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyAddress, strings.Join(addressStrs, ",")),
			sdk.NewAttribute(types.AttributeKeyThreshold, strconv.FormatUint(uint64(threshold), 10)),
		)
	default:
		return nil, fmt.Errorf("uknown key type for chain %s", chain.Name)
	}

	return &types.ConfirmTransferKeyResponse{}, nil
}

func (s msgServer) VoteConfirmChain(c context.Context, req *types.VoteConfirmChainRequest) (*types.VoteConfirmChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	registeredChain, registered := s.nexus.GetChain(ctx, req.Name)
	if registered {
		return &types.VoteConfirmChainResponse{Log: fmt.Sprintf("chain %s already confirmed", registeredChain.Name)}, nil
	}
	chain, ok := s.GetPendingChain(ctx, req.Name)
	if !ok {
		return nil, fmt.Errorf("unknown chain %s", req.Name)
	}

	voter := s.snapshotter.GetOperator(ctx, req.Sender)
	if voter == nil {
		return nil, fmt.Errorf("account %v is not registered as a validator proxy", req.Sender.String())
	}

	poll := s.voter.GetPoll(ctx, req.PollKey)
	voteValue := &gogoprototypes.BoolValue{Value: req.Confirmed}
	if err := poll.Vote(voter, voteValue); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeChainConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueVote),
		sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatBool(voteValue.Value)),
	))

	if poll.Is(vote.Pending) {
		return &types.VoteConfirmChainResponse{Log: fmt.Sprintf("not enough votes to confirm chain in %s yet", req.Name)}, nil
	}

	if poll.Is(vote.Failed) {
		s.DeletePendingChain(ctx, req.Name)
		return &types.VoteConfirmChainResponse{Log: fmt.Sprintf("poll %s failed", poll.GetKey())}, nil
	}

	confirmed, ok := poll.GetResult().(*gogoprototypes.BoolValue)
	if !ok {
		return nil, fmt.Errorf("result of poll %s has wrong type, expected bool, got %T", req.PollKey.String(), poll.GetResult())
	}

	s.Logger(ctx).Info(fmt.Sprintf("EVM chain confirmation result is %s", poll.GetResult()))
	s.DeletePendingChain(ctx, req.Name)

	// handle poll result
	event := sdk.NewEvent(types.EventTypeChainConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyChain, req.Name),
		sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&req.PollKey))))

	if !confirmed.Value {
		poll.AllowOverride()
		ctx.EventManager().EmitEvent(
			event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueReject)))
		return &types.VoteConfirmChainResponse{
			Log: fmt.Sprintf("chain %s was rejected", req.Name),
		}, nil
	}
	ctx.EventManager().EmitEvent(
		event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueConfirm)))

	s.nexus.SetChain(ctx, chain)
	s.nexus.RegisterAsset(ctx, chain.Name, chain.NativeAsset)

	return &types.VoteConfirmChainResponse{}, nil
}

// VoteConfirmDeposit handles votes for deposit confirmations
func (s msgServer) VoteConfirmDeposit(c context.Context, req *types.VoteConfirmDepositRequest) (*types.VoteConfirmDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(chain.Name)
	pendingDeposit, pollFound := keeper.GetPendingDeposit(ctx, req.PollKey)
	confirmedDeposit, state, depositFound := keeper.GetDeposit(ctx, common.Hash(req.TxID), common.Address(req.BurnAddress))

	switch {
	// a malicious user could try to delete an ongoing poll by providing an already confirmed deposit,
	// so we need to check that it matches the poll before deleting
	case depositFound && pollFound && confirmedDeposit == pendingDeposit:
		keeper.DeletePendingDeposit(ctx, req.PollKey)
		fallthrough
	// If the voting threshold has been met and additional votes are received they should not return an error
	case depositFound:
		switch state {
		case types.CONFIRMED:
			return &types.VoteConfirmDepositResponse{Log: fmt.Sprintf("deposit in %s to address %s already confirmed", confirmedDeposit.TxID.Hex(), confirmedDeposit.BurnerAddress.Hex())}, nil
		case types.BURNED:
			return &types.VoteConfirmDepositResponse{Log: fmt.Sprintf("deposit in %s to address %s already spent", confirmedDeposit.TxID.Hex(), confirmedDeposit.BurnerAddress.Hex())}, nil
		}
	case !pollFound:
		return nil, fmt.Errorf("no deposit found for poll %s", req.PollKey.String())
	case pendingDeposit.BurnerAddress != req.BurnAddress || pendingDeposit.TxID != req.TxID:
		return nil, fmt.Errorf("deposit in %s to address %s does not match poll %s", req.TxID.Hex(), req.BurnAddress.Hex(), req.PollKey.String())
	default:
		// assert: the deposit is known and has not been confirmed before
	}

	_, ok = s.nexus.GetChain(ctx, pendingDeposit.DestinationChain)
	if !ok {
		return nil, fmt.Errorf("destination chain %s is not a registered chain", pendingDeposit.DestinationChain)
	}

	voter := s.snapshotter.GetOperator(ctx, req.Sender)
	if voter == nil {
		return nil, fmt.Errorf("account %v is not registered as a validator proxy", req.Sender.String())
	}

	poll := s.voter.GetPoll(ctx, req.PollKey)
	voteValue := &gogoprototypes.BoolValue{Value: req.Confirmed}
	if err := poll.Vote(voter, voteValue); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDepositConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueVote),
		sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatBool(voteValue.Value)),
	))

	if poll.Is(vote.Pending) {
		return &types.VoteConfirmDepositResponse{Log: fmt.Sprintf("not enough votes to confirm deposit in %s to %s yet", req.TxID.Hex(), req.BurnAddress.Hex())}, nil
	}

	if poll.Is(vote.Failed) {
		keeper.DeletePendingDeposit(ctx, req.PollKey)
		return &types.VoteConfirmDepositResponse{Log: fmt.Sprintf("poll %s failed", poll.GetKey())}, nil
	}

	confirmed, ok := poll.GetResult().(*gogoprototypes.BoolValue)
	if !ok {
		return nil, fmt.Errorf("result of poll %s has wrong type, expected bool, got %T", req.PollKey.String(), poll.GetResult())
	}

	s.Logger(ctx).Info(fmt.Sprintf("%s deposit confirmation result is %s", chain.Name, poll.GetResult()))
	keeper.DeletePendingDeposit(ctx, req.PollKey)

	depositAddr := nexus.CrossChainAddress{Address: pendingDeposit.BurnerAddress.Hex(), Chain: chain}
	recipient, ok := s.nexus.GetRecipient(ctx, depositAddr)
	if !ok {
		return nil, fmt.Errorf("cross-chain sender has no recipient")
	}

	// handle poll result
	event := sdk.NewEvent(types.EventTypeDepositConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
		sdk.NewAttribute(types.AttributeKeyDestinationChain, recipient.Chain.Name),
		sdk.NewAttribute(types.AttributeKeyDestinationAddress, recipient.Address),
		sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&req.PollKey))))
	defer func() { ctx.EventManager().EmitEvent(event) }()

	if !confirmed.Value {
		poll.AllowOverride()
		event = event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueReject))
		return &types.VoteConfirmDepositResponse{
			Log: fmt.Sprintf("deposit in %s to %s was discarded", req.TxID.Hex(), req.BurnAddress.Hex()),
		}, nil
	}
	event = event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueConfirm))

	amount := sdk.NewInt64Coin(pendingDeposit.Asset, pendingDeposit.Amount.BigInt().Int64())

	feeRate, ok := keeper.GetTransactionFeeRate(ctx)
	if !ok {
		return nil, fmt.Errorf("could not retrieve transaction fee rate")
	}

	if err := s.nexus.EnqueueForTransfer(ctx, depositAddr, amount, feeRate); err != nil {
		return nil, err
	}
	keeper.SetDeposit(ctx, pendingDeposit, types.CONFIRMED)

	return &types.VoteConfirmDepositResponse{}, nil
}

// VoteConfirmToken handles votes for token deployment confirmations
func (s msgServer) VoteConfirmToken(c context.Context, req *types.VoteConfirmTokenRequest) (*types.VoteConfirmTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(chain.Name)
	token := keeper.GetERC20TokenByAsset(ctx, req.Asset)
	switch {
	case token.Is(types.Confirmed):
		return &types.VoteConfirmTokenResponse{
			Log: fmt.Sprintf("token %s deployment already confirmed", req.Asset)}, nil
	case !token.Is(types.Pending):
		return nil, fmt.Errorf("no open poll for token '%s'", token.GetAsset())
	case types.GetConfirmTokenKey(token.GetTxID(), token.GetAsset()) != req.PollKey:
		return nil, fmt.Errorf("poll key mismatch (expected %s, got %s)", types.GetConfirmTokenKey(token.GetTxID(), token.GetAsset()).String(), req.PollKey.String())
	default:
		// assert: the token is known and has not been confirmed before
	}

	voter := s.snapshotter.GetOperator(ctx, req.Sender)
	if voter == nil {
		return nil, fmt.Errorf("account %v is not registered as a validator proxy", req.Sender.String())
	}

	poll := s.voter.GetPoll(ctx, req.PollKey)
	voteValue := &gogoprototypes.BoolValue{Value: req.Confirmed}
	if err := poll.Vote(voter, voteValue); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeTokenConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueVote),
		sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatBool(voteValue.Value)),
	))

	if poll.Is(vote.Pending) {
		return &types.VoteConfirmTokenResponse{Log: fmt.Sprintf("not enough votes to confirm token %s yet", req.Asset)}, nil
	}

	if poll.Is(vote.Failed) {
		token.RejectDeployment()
		return &types.VoteConfirmTokenResponse{Log: fmt.Sprintf("poll %s failed", poll.GetKey())}, nil
	}

	confirmed, ok := poll.GetResult().(*gogoprototypes.BoolValue)
	if !ok {
		return nil, fmt.Errorf("result of poll %s has wrong type, expected bool, got %T", req.PollKey.String(), poll.GetResult())
	}

	s.Logger(ctx).Info(fmt.Sprintf("token deployment confirmation result is %s", poll.GetResult()))

	// handle poll result
	event := sdk.NewEvent(types.EventTypeTokenConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
		sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&req.PollKey))))

	if !confirmed.Value {
		poll.AllowOverride()
		token.RejectDeployment()
		ctx.EventManager().EmitEvent(
			event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueReject)))
		return &types.VoteConfirmTokenResponse{
			Log: fmt.Sprintf("token %s was discarded", req.Asset),
		}, nil
	}

	ctx.EventManager().EmitEvent(
		event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueConfirm)))

	s.nexus.RegisterAsset(ctx, chain.Name, req.Asset)
	token.ConfirmDeployment()

	return &types.VoteConfirmTokenResponse{
		Log: fmt.Sprintf("token %s deployment confirmed", req.Asset)}, nil
}

// VoteConfirmTransferKey handles votes for transfer ownership/operatorship confirmations
func (s msgServer) VoteConfirmTransferKey(c context.Context, req *types.VoteConfirmTransferKeyRequest) (*types.VoteConfirmTransferKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(chain.Name)

	pendingTransfer, pendingTransferFound := keeper.GetPendingTransferKey(ctx, req.PollKey)
	archivedTransfer, archivedTransferFound := keeper.GetArchivedTransferKey(ctx, req.PollKey)

	var keyRole tss.KeyRole
	switch {
	case !pendingTransferFound && !archivedTransferFound:
		return nil, fmt.Errorf("no transfer key found for poll %s", req.PollKey.String())
	// If the voting threshold has been met and additional votes are received they should not return an error
	case archivedTransferFound:
		return &types.VoteConfirmTransferKeyResponse{Log: fmt.Sprintf("%s in %s to keyID %s already confirmed", archivedTransfer.Type.SimpleString(), archivedTransfer.TxID.Hex(), archivedTransfer.NextKeyID)}, nil
	case pendingTransferFound:
		keyRole = s.signer.GetKeyRole(ctx, pendingTransfer.NextKeyID)
		if keyRole == tss.Unknown {
			return nil, fmt.Errorf("key %s cannot be found", pendingTransfer.NextKeyID)
		}
	default:
		// assert: the transfer ownership/operatorship is known and has not been confirmed before
	}

	voter := s.snapshotter.GetOperator(ctx, req.Sender)
	if voter == nil {
		return nil, fmt.Errorf("account %v is not registered as a validator proxy", req.Sender.String())
	}

	poll := s.voter.GetPoll(ctx, req.PollKey)
	voteValue := &gogoprototypes.BoolValue{Value: req.Confirmed}
	if err := poll.Vote(voter, voteValue); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeTransferKeyConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueVote),
		sdk.NewAttribute(types.AttributeKeyValue, strconv.FormatBool(voteValue.Value)),
	))

	if poll.Is(vote.Pending) {
		return &types.VoteConfirmTransferKeyResponse{Log: fmt.Sprintf("not enough votes to confirm transfer key in poll %s yet", req.PollKey.String())}, nil
	}

	if poll.Is(vote.Failed) {
		keeper.DeletePendingTransferKey(ctx, req.PollKey)
		return &types.VoteConfirmTransferKeyResponse{Log: fmt.Sprintf("poll %s failed", poll.GetKey())}, nil
	}

	confirmed, ok := poll.GetResult().(*gogoprototypes.BoolValue)
	if !ok {
		return nil, fmt.Errorf("result of poll %s has wrong type, expected bool, got %T", req.PollKey.String(), poll.GetResult())
	}

	// handle poll result
	event := sdk.NewEvent(types.EventTypeTransferKeyConfirmation,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyChain, chain.Name),
		sdk.NewAttribute(types.AttributeKeyTransferKeyType, pendingTransfer.Type.SimpleString()),
		sdk.NewAttribute(types.AttributeKeyPoll, string(types.ModuleCdc.MustMarshalJSON(&req.PollKey))))

	if !confirmed.Value {
		poll.AllowOverride()
		ctx.EventManager().EmitEvent(
			event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueReject)))

		msg := fmt.Sprintf("failed to confirmed %s key transfer for chain %s", keyRole.SimpleString(), chain.Name)
		s.Logger(ctx).Error(msg)
		return &types.VoteConfirmTransferKeyResponse{Log: msg}, nil

	}

	keeper.ArchiveTransferKey(ctx, req.PollKey)
	ctx.EventManager().EmitEvent(
		event.AppendAttributes(sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueConfirm)))

	if err := s.signer.RotateKey(ctx, chain, keyRole); err != nil {
		return nil, err
	}

	s.Logger(ctx).Info(fmt.Sprintf("successfully confirmed %s key transfer for chain %s", keyRole.SimpleString(), chain.Name))
	return &types.VoteConfirmTransferKeyResponse{}, nil
}

func (s msgServer) CreateDeployToken(c context.Context, req *types.CreateDeployTokenRequest) (*types.CreateDeployTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	keeper := s.ForChain(req.Chain)

	originChain, found := s.nexus.GetChain(ctx, req.Asset.Chain)
	if !found {
		return nil, fmt.Errorf("%s is not a registered chain", req.Asset.Chain)
	}

	if !s.nexus.IsAssetRegistered(ctx, originChain.Name, req.Asset.Name) {
		return nil, fmt.Errorf("asset %s is not registered on the origin chain %s", originChain.NativeAsset, originChain.Name)
	}

	if _, nextMasterKeyAssigned := s.signer.GetNextKeyID(ctx, chain, tss.MasterKey); nextMasterKeyAssigned {
		return nil, fmt.Errorf("next %s key already assigned for chain %s, rotate key first", tss.MasterKey.SimpleString(), chain.Name)
	}

	masterKeyID, ok := s.signer.GetCurrentKeyID(ctx, chain, tss.MasterKey)
	if !ok {
		return nil, fmt.Errorf("no master key for chain %s found", chain.Name)
	}

	token, err := keeper.CreateERC20Token(ctx, req.Asset.Name, req.TokenDetails)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to initialize token %s(%s) for chain %s", req.TokenDetails.TokenName, req.TokenDetails.Symbol, chain.Name)
	}

	cmd, err := token.CreateDeployCommand(masterKeyID)
	if err != nil {
		return nil, err
	}

	if err := keeper.EnqueueCommand(ctx, cmd); err != nil {
		return nil, err
	}

	return &types.CreateDeployTokenResponse{}, nil
}

func (s msgServer) CreateBurnTokens(c context.Context, req *types.CreateBurnTokensRequest) (*types.CreateBurnTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keeper := s.ForChain(req.Chain)

	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	deposits := keeper.GetConfirmedDeposits(ctx)
	if len(deposits) == 0 {
		return &types.CreateBurnTokensResponse{}, nil
	}

	chainID := s.getChainID(ctx, req.Chain)
	if chainID == nil {
		return nil, fmt.Errorf("could not find chain ID for '%s'", req.Chain)
	}

	if _, nextSecondaryKeyAssigned := s.signer.GetNextKeyID(ctx, chain, tss.SecondaryKey); nextSecondaryKeyAssigned {
		return nil, fmt.Errorf("next %s key already assigned for chain %s, rotate key first", tss.SecondaryKey.SimpleString(), chain.Name)
	}

	secondaryKeyID, ok := s.signer.GetCurrentKeyID(ctx, chain, tss.SecondaryKey)
	if !ok {
		return nil, fmt.Errorf("no %s key for chain %s found", tss.SecondaryKey.SimpleString(), chain.Name)
	}

	seen := map[string]bool{}
	for _, deposit := range deposits {
		keeper.DeleteDeposit(ctx, deposit)
		keeper.SetDeposit(ctx, deposit, types.BURNED)

		burnerAddressHex := deposit.BurnerAddress.Hex()

		if seen[burnerAddressHex] {
			continue
		}

		burnerInfo := keeper.GetBurnerInfo(ctx, common.Address(deposit.BurnerAddress))
		if burnerInfo == nil {
			return nil, fmt.Errorf("no burner info found for address %s", burnerAddressHex)
		}

		cmd, err := types.CreateBurnTokenCommand(chainID, secondaryKeyID, ctx.BlockHeight(), *burnerInfo)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to create burn-token command to burn token at address %s for chain %s", burnerAddressHex, chain.Name)
		}

		if err := keeper.EnqueueCommand(ctx, cmd); err != nil {
			return nil, err
		}

		seen[burnerAddressHex] = true
	}

	return &types.CreateBurnTokensResponse{}, nil
}

func getMultisigThreshold(keyCount int, threshold utils.Threshold) uint8 {
	return uint8(
		sdk.NewDec(int64(keyCount)).
			MulInt64(threshold.Numerator).
			QuoInt64(threshold.Denominator).
			Ceil().
			RoundInt64(),
	)
}

func getMultisigAddresses(key tss.Key) ([]common.Address, uint8, error) {
	multisigPubKeys, err := key.GetMultisigPubKey()
	if err != nil {
		return nil, 0, sdkerrors.Wrapf(types.ErrEVM, err.Error())
	}

	threshold := uint8(key.GetMultisigKey().Threshold)
	return types.KeysToAddresses(multisigPubKeys...), threshold, nil
}

func getGatewayDeploymentBytecode(ctx sdk.Context, k types.ChainKeeper, s types.Signer, chain nexus.Chain) ([]byte, error) {
	externalKeyIDs, ok := s.GetExternalKeyIDs(ctx, chain)
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrEVM, fmt.Sprintf("no %s keys for chain %s found", tss.ExternalKey.SimpleString(), chain.Name))
	}

	externalPubKeys := make([]ecdsa.PublicKey, len(externalKeyIDs))
	for i, externalKeyID := range externalKeyIDs {
		externalKey, ok := s.GetKey(ctx, externalKeyID)
		if !ok {
			return nil, sdkerrors.Wrap(types.ErrEVM, fmt.Sprintf("%s key %s for chain %s not found", tss.ExternalKey.SimpleString(), externalKeyID, chain.Name))
		}

		pk, err := externalKey.GetECDSAPubKey()
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrEVM, err.Error())
		}

		externalPubKeys[i] = pk
	}
	externalKeyAddresses := types.KeysToAddresses(externalPubKeys...)
	externalKeyThreshold := getMultisigThreshold(len(externalKeyAddresses), s.GetExternalMultisigThreshold(ctx))

	bz, _ := k.GetGatewayByteCodes(ctx)

	masterKey, ok := s.GetCurrentKey(ctx, chain, tss.MasterKey)
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrEVM, fmt.Sprintf("no %s key for chain %s found", tss.MasterKey.SimpleString(), chain.Name))
	}

	secondaryKey, ok := s.GetCurrentKey(ctx, chain, tss.SecondaryKey)
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrEVM, fmt.Sprintf("no %s key for chain %s found", tss.SecondaryKey.SimpleString(), chain.Name))
	}

	switch chain.KeyType {
	case tss.Threshold:
		masterPubKey, err := masterKey.GetECDSAPubKey()
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrEVM, err.Error())
		}

		secondaryPubKey, err := secondaryKey.GetECDSAPubKey()
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrEVM, err.Error())
		}

		return types.GetSinglesigGatewayDeploymentBytecode(
			bz,
			externalKeyAddresses,
			uint8(s.GetExternalMultisigThreshold(ctx).Numerator),
			crypto.PubkeyToAddress(masterPubKey),
			crypto.PubkeyToAddress(secondaryPubKey),
		)
	case tss.Multisig:
		masterMultisigAddresses, masterMultisigThreshold, err := getMultisigAddresses(masterKey)
		if err != nil {
			return nil, err
		}

		secondaryMultisigAddresses, secondaryMultisigThreshold, err := getMultisigAddresses(secondaryKey)
		if err != nil {
			return nil, err
		}

		return types.GetMultisigGatewayDeploymentBytecode(
			bz,
			externalKeyAddresses,
			externalKeyThreshold,
			masterMultisigAddresses,
			masterMultisigThreshold,
			secondaryMultisigAddresses,
			secondaryMultisigThreshold,
		)
	default:
		return nil, fmt.Errorf("unknown key type set for chain %s", chain.Name)
	}
}

func (s msgServer) CreatePendingTransfers(c context.Context, req *types.CreatePendingTransfersRequest) (*types.CreatePendingTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keeper := s.ForChain(req.Chain)

	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	pendingTransfers := s.nexus.GetTransfersForChain(ctx, chain, nexus.Pending)
	if len(pendingTransfers) == 0 {
		return &types.CreatePendingTransfersResponse{}, nil
	}

	if _, nextSecondaryKeyAssigned := s.signer.GetNextKeyID(ctx, chain, tss.SecondaryKey); nextSecondaryKeyAssigned {
		return nil, fmt.Errorf("next %s key already assigned for chain %s, rotate key first", tss.SecondaryKey.SimpleString(), chain.Name)
	}

	secondaryKeyID, ok := s.signer.GetCurrentKeyID(ctx, chain, tss.SecondaryKey)
	if !ok {
		return nil, fmt.Errorf("no %s key for chain %s found", tss.SecondaryKey.SimpleString(), chain.Name)
	}

	getRecipientAndAsset := func(transfer nexus.CrossChainTransfer) string {
		return fmt.Sprintf("%s-%s", transfer.Recipient.Address, transfer.Asset.Denom)
	}
	transfers := nexus.MergeTransfersBy(pendingTransfers, getRecipientAndAsset)

	for _, transfer := range transfers {
		token := keeper.GetERC20TokenByAsset(ctx, transfer.Asset.Denom)
		cmd, err := token.CreateMintCommand(secondaryKeyID, transfer)

		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed create mint-token command for transfer %d", transfer.ID)
		}

		s.Logger(ctx).Info(fmt.Sprintf("storing data for mint command %s", cmd.ID.Hex()))

		if err := keeper.EnqueueCommand(ctx, cmd); err != nil {
			return nil, err
		}
	}

	for _, pendingTransfer := range pendingTransfers {
		s.nexus.ArchivePendingTransfer(ctx, pendingTransfer)
	}

	return &types.CreatePendingTransfersResponse{}, nil
}

func (s msgServer) createTransferKeyCommand(ctx sdk.Context, transferKeyType types.TransferKeyType, chainStr string, nextKeyID tss.KeyID) (types.Command, error) {
	chain, ok := s.nexus.GetChain(ctx, chainStr)
	if !ok {
		return types.Command{}, fmt.Errorf("%s is not a registered chain", chainStr)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return types.Command{}, err
	}

	chainID := s.getChainID(ctx, chainStr)
	if chainID == nil {
		return types.Command{}, fmt.Errorf("could not find chain ID for '%s'", chainStr)
	}

	var keyRole tss.KeyRole
	switch transferKeyType {
	case types.Ownership:
		keyRole = tss.MasterKey
	case types.Operatorship:
		keyRole = tss.SecondaryKey
	default:
		return types.Command{}, fmt.Errorf("invalid transfer key type %s", transferKeyType.SimpleString())
	}

	// don't allow any transfer key if the next master/secondary key is already assigned
	if _, nextMasterKeyAssigned := s.signer.GetNextKeyID(ctx, chain, tss.MasterKey); nextMasterKeyAssigned {
		return types.Command{}, fmt.Errorf("next %s key already assigned for chain %s, rotate key first", tss.MasterKey.SimpleString(), chain.Name)
	}
	if _, nextSecondaryKeyAssigned := s.signer.GetNextKeyID(ctx, chain, tss.SecondaryKey); nextSecondaryKeyAssigned {
		return types.Command{}, fmt.Errorf("next %s key already assigned for chain %s, rotate key first", tss.SecondaryKey.SimpleString(), chain.Name)
	}

	if err := s.signer.AssertMatchesRequirements(ctx, s.snapshotter, chain, nextKeyID, keyRole); err != nil {
		return types.Command{}, sdkerrors.Wrapf(err, "key %s does not match requirements for role %s", nextKeyID, keyRole.SimpleString())
	}

	if err := s.signer.AssignNextKey(ctx, chain, keyRole, nextKeyID); err != nil {
		return types.Command{}, err
	}

	currMasterKeyID, ok := s.signer.GetCurrentKeyID(ctx, chain, tss.MasterKey)
	if !ok {
		return types.Command{}, fmt.Errorf("current %s key not set for chain %s", tss.MasterKey, chain.Name)
	}

	nextKey, ok := s.signer.GetKey(ctx, nextKeyID)
	if !ok {
		return types.Command{}, fmt.Errorf("could not find threshold key '%s'", nextKeyID)
	}

	switch chain.KeyType {
	case tss.Threshold:
		pk, err := nextKey.GetECDSAPubKey()
		if err != nil {
			return types.Command{}, err
		}

		address := crypto.PubkeyToAddress(pk)
		s.Logger(ctx).Debug(fmt.Sprintf("creating command %s for chain %s to transfer to address %s", transferKeyType.SimpleString(), chain.Name, address))

		return types.CreateSinglesigTransferCommand(transferKeyType, chainID, currMasterKeyID, crypto.PubkeyToAddress(pk))
	case tss.Multisig:
		addresses, threshold, err := getMultisigAddresses(nextKey)
		if err != nil {
			return types.Command{}, err
		}

		addressStrs := make([]string, len(addresses))
		for i, address := range addresses {
			addressStrs[i] = address.Hex()
		}

		s.Logger(ctx).Debug(fmt.Sprintf("creating command %s for chain %s to transfer to addresses %s", transferKeyType.SimpleString(), chain.Name, strings.Join(addressStrs, ",")))

		return types.CreateMultisigTransferCommand(transferKeyType, chainID, currMasterKeyID, threshold, addresses...)
	default:
		return types.Command{}, fmt.Errorf("invalid key type '%s'", chain.KeyType.SimpleString())
	}
}

func (s msgServer) CreateTransferOwnership(c context.Context, req *types.CreateTransferOwnershipRequest) (*types.CreateTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keeper := s.ForChain(req.Chain)

	if _, ok := keeper.GetGatewayAddress(ctx); !ok {
		return nil, fmt.Errorf("axelar gateway address not set")
	}

	cmd, err := s.createTransferKeyCommand(ctx, types.Ownership, req.Chain, req.KeyID)
	if err != nil {
		return nil, err
	}

	if err := keeper.EnqueueCommand(ctx, cmd); err != nil {
		return nil, err
	}

	return &types.CreateTransferOwnershipResponse{}, nil
}

func (s msgServer) CreateTransferOperatorship(c context.Context, req *types.CreateTransferOperatorshipRequest) (*types.CreateTransferOperatorshipResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keeper := s.ForChain(req.Chain)

	if _, ok := keeper.GetGatewayAddress(ctx); !ok {
		return nil, fmt.Errorf("axelar gateway address not set")
	}

	cmd, err := s.createTransferKeyCommand(ctx, types.Operatorship, req.Chain, req.KeyID)
	if err != nil {
		return nil, err
	}

	if err := keeper.EnqueueCommand(ctx, cmd); err != nil {
		return nil, err
	}

	return &types.CreateTransferOperatorshipResponse{}, nil
}

func (s msgServer) SignCommands(c context.Context, req *types.SignCommandsRequest) (*types.SignCommandsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	chain, ok := s.nexus.GetChain(ctx, req.Chain)
	if !ok {
		return nil, fmt.Errorf("%s is not a registered chain", req.Chain)
	}

	if err := validateChainActivated(ctx, s.nexus, chain); err != nil {
		return nil, err
	}

	chainID := s.getChainID(ctx, req.Chain)
	if chainID == nil {
		return nil, fmt.Errorf("could not find chain ID for '%s'", req.Chain)
	}

	keeper := s.ForChain(chain.Name)
	id, err := keeper.CreateNewBatchToSign(ctx)
	if err != nil {
		return nil, err
	}

	// if no error was thrown above, the batch exists
	batchedCommands := keeper.GetBatchByID(ctx, id)

	counter, ok := s.signer.GetSnapshotCounterForKeyID(ctx, batchedCommands.GetKeyID())
	if !ok {
		return nil, fmt.Errorf("no snapshot counter for key ID %s registered", batchedCommands.GetKeyID())
	}

	sigMetadata := types.SigMetadata{
		Type:  types.SigCommand,
		Chain: chain.Name,
	}

	batchedCommandsIDHex := hex.EncodeToString(batchedCommands.GetID())
	err = s.signer.StartSign(ctx, tss.SignInfo{
		KeyID:           batchedCommands.GetKeyID(),
		SigID:           batchedCommandsIDHex,
		Msg:             batchedCommands.GetSigHash().Bytes(),
		SnapshotCounter: counter,
		RequestModule:   types.ModuleName,
		Metadata:        string(types.ModuleCdc.MustMarshalJSON(&sigMetadata)),
	}, s.snapshotter, s.voter)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyChain, req.Chain),
			sdk.NewAttribute(sdk.AttributeKeySender, req.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyBatchedCommandsID, batchedCommandsIDHex),
		),
	)

	return &types.SignCommandsResponse{BatchedCommandsID: batchedCommands.GetID()}, nil
}

func (s msgServer) AddChain(c context.Context, req *types.AddChainRequest) (*types.AddChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if _, found := s.nexus.GetChain(ctx, req.Name); found {
		return nil, fmt.Errorf("chain '%s' is already registered", req.Name)
	}

	if err := req.Params.Validate(); err != nil {
		return nil, err
	}

	if !tsstypes.TSSEnabled && req.KeyType == tss.Threshold {
		return nil, fmt.Errorf("TSS is disabled")
	}

	s.SetPendingChain(ctx, nexus.Chain{Name: req.Name, NativeAsset: req.NativeAsset, SupportsForeignAssets: true, KeyType: req.KeyType})
	s.SetParams(ctx, req.Params)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeNewChain,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.AttributeValueUpdate),
			sdk.NewAttribute(types.AttributeKeyChain, req.Name),
			sdk.NewAttribute(types.AttributeKeyNativeAsset, req.NativeAsset),
		),
	)

	return &types.AddChainResponse{}, nil
}

func (s msgServer) getChainID(ctx sdk.Context, chain string) (chainID *big.Int) {
	for _, p := range s.GetParams(ctx) {
		if strings.EqualFold(p.Chain, chain) {
			chainID = s.ForChain(chain).GetChainIDByNetwork(ctx, p.Network)
		}
	}

	return
}
