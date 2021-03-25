package keeper

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/axelarnetwork/axelar-core/testutils"
	"github.com/axelarnetwork/axelar-core/testutils/rand"
	btc "github.com/axelarnetwork/axelar-core/x/bitcoin/exported"
	btcTypes "github.com/axelarnetwork/axelar-core/x/bitcoin/types"
	eth "github.com/axelarnetwork/axelar-core/x/ethereum/exported"
	"github.com/axelarnetwork/axelar-core/x/nexus/exported"
	"github.com/axelarnetwork/axelar-core/x/nexus/types"

	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/axelar-core/testutils/fake"
)

const (
	addrMaxLength int   = 20
	maxAmount     int64 = 100000000000
	linkedAddr    int   = 50
)

var keeper Keeper

func init() {
	cdc := testutils.Codec()
	nexusSubspace := params.NewSubspace(testutils.Codec(), sdk.NewKVStoreKey("nexusKey"), sdk.NewKVStoreKey("tNexusKey"), "nexus")
	keeper = NewKeeper(cdc, sdk.NewKVStoreKey("nexus"), nexusSubspace)
}

func TestLinkNoForeignAssetSupport(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())

	sender, recipient := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
	keeper.LinkAddresses(ctx, sender, recipient)
	err := keeper.EnqueueForTransfer(ctx, sender, makeRandAmount(makeRandomDenom()))
	assert.Error(t, err)
}

func TestLinkSuccess(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())

	sender, recipient := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
	keeper.LinkAddresses(ctx, sender, recipient)
	err := keeper.EnqueueForTransfer(ctx, sender, makeRandAmount(btcTypes.Satoshi))
	assert.NoError(t, err)
	recp, ok := keeper.GetRecipient(ctx, sender)
	assert.True(t, ok)
	assert.Equal(t, recipient, recp)

	sender.Address = rand.Str(20)
	err = keeper.EnqueueForTransfer(ctx, sender, makeRandAmount(btcTypes.Satoshi))
	assert.Error(t, err)
	recp, ok = keeper.GetRecipient(ctx, sender)
	assert.False(t, ok)
	assert.NotEqual(t, recipient, recp)
}

func TestPrepareNoLink(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())

	sender, _ := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
	err := keeper.EnqueueForTransfer(ctx, sender, makeRandAmount(btcTypes.Satoshi))
	assert.Error(t, err)
}

func TestPrepareSuccess(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())

	amounts := make(map[exported.CrossChainAddress]sdk.Coin)
	for i := 0; i < linkedAddr; i++ {
		sender, recipient := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
		amounts[recipient] = makeRandAmount(btcTypes.Satoshi)
		keeper.LinkAddresses(ctx, sender, recipient)
		err := keeper.EnqueueForTransfer(ctx, sender, amounts[recipient])
		assert.NoError(t, err)
	}

	transfers := keeper.GetPendingTransfersForChain(ctx, eth.Ethereum)
	assert.Equal(t, len(transfers), len(amounts))
	assert.Equal(t, linkedAddr, len(transfers))

	count := 0
	for _, transfer := range transfers {
		amount, ok := amounts[transfer.Recipient]
		if ok {
			count++
			assert.Equal(t, transfer.Asset, amount)
		}
	}
	assert.Equal(t, linkedAddr, count)
}

func TestArchive(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())

	recipients := make([]exported.CrossChainAddress, 0)
	var total uint64 = 0

	for i := 0; i < linkedAddr; i++ {
		sender, recipient := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
		recipients = append(recipients, recipient)
		keeper.LinkAddresses(ctx, sender, recipient)
		amount := makeRandAmount(btcTypes.Satoshi)
		err := keeper.EnqueueForTransfer(ctx, sender, amount)
		assert.NoError(t, err)
		total += amount.Amount.Uint64()
	}

	transfers := keeper.GetPendingTransfersForChain(ctx, eth.Ethereum)

	for _, transfer := range transfers {
		keeper.ArchivePendingTransfer(ctx, transfer)
	}

	archived := keeper.GetArchivedTransfersForChain(ctx, eth.Ethereum)
	assert.Equal(t, linkedAddr, len(archived))

	count := 0
	for _, archive := range archived {
		for _, transfer := range transfers {
			if transfer.Recipient.Address == archive.Recipient.Address {
				count++
				assert.Equal(t, archive.Asset, transfer.Asset)
			}
		}
	}
	assert.Equal(t, linkedAddr, count)
	assert.Equal(t, 0, len(keeper.GetPendingTransfersForChain(ctx, eth.Ethereum)))
}

func TestTotalInvalid(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())
	btcSender, btcRecipient := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
	keeper.LinkAddresses(ctx, btcSender, btcRecipient)
	ethSender, ethRecipient := makeRandAddressesForChain(eth.Ethereum, btc.Bitcoin)
	keeper.LinkAddresses(ctx, ethSender, ethRecipient)

	err := keeper.EnqueueForTransfer(ctx, btcSender, makeRandAmount(btcTypes.Satoshi))
	assert.NoError(t, err)
	transfer := keeper.GetPendingTransfersForChain(ctx, eth.Ethereum)[0]
	keeper.ArchivePendingTransfer(ctx, transfer)
	total := transfer.Asset.Amount.Int64()
	amount := sdk.NewCoin(btcTypes.Satoshi, sdk.NewInt(total+rand.I64Between(1, 100000)))
	err = keeper.EnqueueForTransfer(ctx, ethSender, amount)
	assert.Error(t, err)
}

func TestTotalSucess(t *testing.T) {
	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetParams(ctx, types.DefaultParams())
	btcSender, btcRecipient := makeRandAddressesForChain(btc.Bitcoin, eth.Ethereum)
	keeper.LinkAddresses(ctx, btcSender, btcRecipient)

	ethSender, ethRecipient := makeRandAddressesForChain(eth.Ethereum, btc.Bitcoin)
	keeper.LinkAddresses(ctx, ethSender, ethRecipient)

	err := keeper.EnqueueForTransfer(ctx, btcSender, makeRandAmount(btcTypes.Satoshi))
	assert.NoError(t, err)
	transfer := keeper.GetPendingTransfersForChain(ctx, eth.Ethereum)[0]
	keeper.ArchivePendingTransfer(ctx, transfer)
	total := transfer.Asset.Amount.Int64()
	amount := sdk.NewCoin(btcTypes.Satoshi, sdk.NewInt(rand.I64Between(1, total)))
	err = keeper.EnqueueForTransfer(ctx, ethSender, amount)
	assert.NoError(t, err)
	amount = sdk.NewCoin(btcTypes.Satoshi, sdk.NewInt(total))
	err = keeper.EnqueueForTransfer(ctx, ethSender, amount)
	assert.Error(t, err)
}

func TestSetChainGetChain_MixCaseChainName(t *testing.T) {
	chainName := strings.ToUpper(rand.StrBetween(5, 10)) + strings.ToLower(rand.StrBetween(5, 10))
	chain := exported.Chain{
		Name:                  chainName,
		NativeAsset:           rand.Str(3),
		SupportsForeignAssets: true,
	}

	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetChain(ctx, chain)

	actual, ok := keeper.GetChain(ctx, strings.ToUpper(chainName))

	assert.True(t, ok)
	assert.Equal(t, chain, actual)

	actual, ok = keeper.GetChain(ctx, strings.ToLower(chainName))

	assert.True(t, ok)
	assert.Equal(t, chain, actual)
}

func TestSetChainGetChain_UpperCaseChainName(t *testing.T) {
	chainName := strings.ToUpper(rand.StrBetween(5, 10))
	chain := exported.Chain{
		Name:                  chainName,
		NativeAsset:           rand.Str(3),
		SupportsForeignAssets: true,
	}

	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetChain(ctx, chain)

	actual, ok := keeper.GetChain(ctx, strings.ToUpper(chainName))

	assert.True(t, ok)
	assert.Equal(t, chain, actual)

	actual, ok = keeper.GetChain(ctx, strings.ToLower(chainName))

	assert.True(t, ok)
	assert.Equal(t, chain, actual)
}

func TestSetChainGetChain_LowerCaseChainName(t *testing.T) {
	chainName := strings.ToLower(rand.StrBetween(5, 10))
	chain := exported.Chain{
		Name:                  chainName,
		NativeAsset:           rand.Str(3),
		SupportsForeignAssets: true,
	}

	ctx := sdk.NewContext(fake.NewMultiStore(), abci.Header{}, false, log.TestingLogger())
	keeper.SetChain(ctx, chain)

	actual, ok := keeper.GetChain(ctx, strings.ToUpper(chainName))

	assert.True(t, ok)
	assert.Equal(t, chain, actual)

	actual, ok = keeper.GetChain(ctx, strings.ToLower(chainName))

	assert.True(t, ok)
	assert.Equal(t, chain, actual)
}

func makeRandomDenom() string {
	d := rand.Strings(3, 3).WithAlphabet([]rune("abcdefghijklmnopqrstuvwxyz")).Take(1)
	return d[0]
}

func makeRandAmount(denom string) sdk.Coin {

	return sdk.NewCoin(denom, sdk.NewInt(rand.I64Between(1, maxAmount)))
}

func makeRandAddressesForChain(origin, distination exported.Chain) (exported.CrossChainAddress, exported.CrossChainAddress) {
	sender := exported.CrossChainAddress{
		Address: rand.Str(addrMaxLength),
		Chain:   origin,
	}
	recipient := exported.CrossChainAddress{
		Address: rand.Str(addrMaxLength),
		Chain:   distination,
	}

	return sender, recipient
}
