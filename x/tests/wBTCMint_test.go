package tests

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	sdk "github.com/cosmos/cosmos-sdk/types"
	goEth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/axelarnetwork/axelar-core/testutils"
	btc "github.com/axelarnetwork/axelar-core/x/bitcoin/exported"
	btcKeeper "github.com/axelarnetwork/axelar-core/x/bitcoin/keeper"
	btcTypes "github.com/axelarnetwork/axelar-core/x/bitcoin/types"
	broadcastTypes "github.com/axelarnetwork/axelar-core/x/broadcast/types"
	eth "github.com/axelarnetwork/axelar-core/x/ethereum/exported"
	ethKeeper "github.com/axelarnetwork/axelar-core/x/ethereum/keeper"
	"github.com/axelarnetwork/axelar-core/x/ethereum/types"
	ethTypes "github.com/axelarnetwork/axelar-core/x/ethereum/types"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	tssTypes "github.com/axelarnetwork/axelar-core/x/tss/types"
)

// 0. Create and start a chain
// 1. Get a deposit address for the given Ethereum recipient address
// 2. Send BTC to the deposit address and wait until confirmed
// 3. Collect all information that needs to be verified about the deposit
// 4. Verify the previously received information
// 5. Wait until verification is complete
// 6. Sign all pending transfers to Ethereum
// 7. Submit the minting command from an externally controlled address to AxelarGateway

func Test_wBTC_mint(t *testing.T) {
	randStrings := testutils.RandStrings(5, 50)
	defer randStrings.Stop()

	// 0. Set up chain
	const nodeCount = 10
	stringGen := testutils.RandStrings(5, 50).Distinct()
	defer stringGen.Stop()

	// create a chain with nodes and assign them as validators
	chain, nodeData := initChain(nodeCount, "mint")
	keygenDone, verifyDone, signDone := registerEventListeners(nodeData[0].Node)

	// register proxies for all validators
	for i, proxy := range randStrings.Take(nodeCount) {
		res := <-chain.Submit(broadcastTypes.MsgRegisterProxy{Principal: nodeData[i].Validator.OperatorAddress, Proxy: sdk.AccAddress(proxy)})
		assert.NoError(t, res.Error)
	}

	// create master key for btc
	btcMasterKeyID := randStrings.Next()
	btcMasterKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// prepare mocks with the btc master key
	var correctBTCKeygens []<-chan bool
	for _, n := range nodeData {
		correctBTCKeygens = append(correctBTCKeygens, prepareKeygen(n.Mocks.Keygen, btcMasterKeyID, btcMasterKey.PublicKey))
	}

	// start keygen
	btcKeygenResult := <-chain.Submit(
		tssTypes.MsgKeygenStart{Sender: randomSender(), NewKeyID: btcMasterKeyID})
	assert.NoError(t, btcKeygenResult.Error)
	for _, isCorrect := range correctBTCKeygens {
		assert.True(t, <-isCorrect)
	}

	// create master key for eth
	ethMasterKeyID := randStrings.Next()
	ethMasterKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// prepare mocks with the eth master key
	var correctETHKeygens []<-chan bool
	for _, n := range nodeData {
		correctETHKeygens = append(correctETHKeygens, prepareKeygen(n.Mocks.Keygen, ethMasterKeyID, ethMasterKey.PublicKey))
	}

	// start keygen
	ethKeygenResult := <-chain.Submit(
		tssTypes.MsgKeygenStart{Sender: randomSender(), NewKeyID: ethMasterKeyID})
	assert.NoError(t, ethKeygenResult.Error)
	for _, isCorrect := range correctETHKeygens {
		assert.True(t, <-isCorrect)
	}

	// wait for voting to be done
	if err := waitFor(keygenDone, 1); err != nil {
		assert.FailNow(t, "keygen", err)
	}

	// assign bitcoin master key
	assignBTCKeyResult := <-chain.Submit(
		tssTypes.MsgAssignNextMasterKey{Sender: randomSender(), Chain: btc.Bitcoin.Name, KeyID: btcMasterKeyID})
	assert.NoError(t, assignBTCKeyResult.Error)

	// assign ethereum master key
	assignETHKeyResult := <-chain.Submit(
		tssTypes.MsgAssignNextMasterKey{Sender: randomSender(), Chain: eth.Ethereum.Name, KeyID: ethMasterKeyID})
	assert.NoError(t, assignETHKeyResult.Error)

	// rotate to the btc master key
	btcRotateResult := <-chain.Submit(
		tssTypes.MsgRotateMasterKey{Sender: randomSender(), Chain: btc.Bitcoin.Name})
	assert.NoError(t, btcRotateResult.Error)

	// rotate to the eth master key
	ethRotateResult := <-chain.Submit(
		tssTypes.MsgRotateMasterKey{Sender: randomSender(), Chain: eth.Ethereum.Name})
	assert.NoError(t, ethRotateResult.Error)

	// prepare caches for upcoming signatures
	totalDepositCount := int(testutils.RandIntBetween(1, 20))
	var correctSigns []<-chan bool
	cache := NewSignatureCache(totalDepositCount + 1)
	for _, n := range nodeData {
		correctSign := prepareSign(n.Mocks.Tofnd, ethMasterKeyID, ethMasterKey, cache)
		correctSigns = append(correctSigns, correctSign)
	}

	// setup axelar gateway
	bz, err := nodeData[0].Node.Query(
		[]string{ethTypes.QuerierRoute, ethKeeper.CreateDeployTx},
		abci.RequestQuery{
			Data: testutils.Codec().MustMarshalJSON(
				ethTypes.DeployParams{
					GasPrice: sdk.NewInt(1),
					GasLimit: 3000000,
				})},
	)
	assert.NoError(t, err)
	var result types.DeployResult
	testutils.Codec().MustUnmarshalJSON(bz, &result)

	deployGatewayResult := <-chain.Submit(
		ethTypes.MsgSignTx{Sender: randomSender(), Tx: testutils.Codec().MustMarshalJSON(result.Tx)})
	assert.NoError(t, deployGatewayResult.Error)

	for _, isCorrect := range correctSigns {
		assert.True(t, <-isCorrect)
	}

	// wait for voting to be done (signing takes longer to tally up)
	if err := waitFor(signDone, 1); err != nil {
		assert.FailNow(t, "signing", err)
	}

	bz, err = nodeData[0].Node.Query(
		[]string{ethTypes.QuerierRoute, ethKeeper.SendTx, string(deployGatewayResult.Data)},
		abci.RequestQuery{Data: nil},
	)

	// steps followed as per https://github.com/axelarnetwork/axelarate#mint-erc20-wrapped-bitcoin-tokens-on-ethereum

	for i := 0; i < totalDepositCount; i++ {
		// 1. Get a deposit address for an Ethereum recipient address
		// we don't provide an actual recipient address, so it is created automatically
		crosschainAddr := nexus.CrossChainAddress{Chain: eth.Ethereum, Address: testutils.RandStringBetween(5, 20)}
		res := <-chain.Submit(btcTypes.NewMsgLink(randomSender(), crosschainAddr.Address, crosschainAddr.Chain.Name))
		assert.NoError(t, res.Error)
		depositAddr := string(res.Data)

		// 2. Send BTC to the deposit address and wait until confirmed
		expectedDepositInfo := randomOutpointInfo(depositAddr)
		for _, n := range nodeData {
			n.Mocks.BTC.GetOutPointInfoFunc = func(bHash *chainhash.Hash, out *wire.OutPoint) (btcTypes.OutPointInfo, error) {
				if bHash.IsEqual(expectedDepositInfo.BlockHash) && out.String() == expectedDepositInfo.OutPoint.String() {
					return expectedDepositInfo, nil
				}
				return btcTypes.OutPointInfo{}, fmt.Errorf("outpoint info not found")
			}
		}

		// 3. Collect all information that needs to be verified about the deposit
		bz, err := nodeData[0].Node.Query(
			[]string{btcTypes.QuerierRoute, btcKeeper.QueryOutInfo, expectedDepositInfo.BlockHash.String()},
			abci.RequestQuery{Data: testutils.Codec().MustMarshalJSON(expectedDepositInfo.OutPoint)})
		assert.NoError(t, err)
		var info btcTypes.OutPointInfo
		testutils.Codec().MustUnmarshalJSON(bz, &info)

		// 4. Verify the previously received information
		res = <-chain.Submit(btcTypes.NewMsgVerifyTx(randomSender(), info))
		assert.NoError(t, res.Error)

	}

	// 5. Wait until verification is complete
	if err := waitFor(verifyDone, totalDepositCount); err != nil {
		assert.FailNow(t, "verification", err)
	}

	// 6. Sign all pending transfers to Ethereum
	res := <-chain.Submit(ethTypes.NewMsgSignPendingTransfers(randomSender()))
	assert.NoError(t, res.Error)
	for _, isCorrect := range correctSigns {
		assert.True(t, <-isCorrect)
	}

	commandID := common.BytesToHash(res.Data)

	// wait for voting to be done (signing takes longer to tally up)
	if err := waitFor(signDone, 1); err != nil {
		assert.FailNow(t, "signing", err)
	}

	// 7. Submit the minting command from an externally controlled address to AxelarGateway
	nodeData[0].Mocks.ETH.SendAndSignTransactionFunc = func(_ context.Context, _ goEth.CallMsg) (string, error) {
		return "", nil
	}

	sender := randomSender()

	_, err = nodeData[0].Node.Query(
		[]string{ethTypes.QuerierRoute, ethKeeper.SendCommand},
		abci.RequestQuery{
			Data: testutils.Codec().MustMarshalJSON(
				ethTypes.CommandParams{
					CommandID: ethTypes.CommandID(commandID),
					Sender:    sender.String(),
				})},
	)
	assert.NoError(t, err)
}
