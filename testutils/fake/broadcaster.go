package fake

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axelarnetwork/axelar-core/x/broadcast/exported"
)

var _ exported.Broadcaster = Broadcaster{}

// Broadcaster for module integration tests
type Broadcaster struct {
	submitMsg      func(msg sdk.Msg) (result <-chan *Result)
	val2Proxy      map[string]sdk.AccAddress
	proxy2Val      map[string]sdk.ValAddress
	LocalPrincipal sdk.ValAddress
	cdc            *codec.Codec
}

// NewBroadcaster creates a new broadcaster fake that sends messages to the blockchainIn channel.
func NewBroadcaster(cdc *codec.Codec, localPrincipal sdk.ValAddress, submitMsg func(msg sdk.Msg) (result <-chan *Result)) Broadcaster {
	return Broadcaster{
		cdc:            cdc,
		submitMsg:      submitMsg,
		val2Proxy:      make(map[string]sdk.AccAddress),
		proxy2Val:      make(map[string]sdk.ValAddress),
		LocalPrincipal: localPrincipal,
	}
}

// Broadcast submits the given messages to the blockchain
func (b Broadcaster) Broadcast(ctx sdk.Context, msgs []exported.MsgWithSenderSetter) error {
	for _, msg := range msgs {
		proxy := b.GetProxy(ctx, b.LocalPrincipal)
		if proxy == nil {
			return fmt.Errorf("proxy not set")
		}
		msg.SetSender(proxy)

		/*
			exported.MsgWithSenderSetter is usually implemented by a pointer.
			However, handlers expect to receive the message by value and do a switch on the message type.
			If they receive the pointer they won't recognize the correct message type.
			By marshalling and unmarshalling into sdk.Msg we get the message by value.
		*/
		bz := b.cdc.MustMarshalBinaryLengthPrefixed(msg)
		var sentMsg sdk.Msg
		b.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &sentMsg)

		_ = b.submitMsg(sentMsg)
	}
	return nil
}

// RegisterProxy registers the given proxy for the given validator principal
func (b Broadcaster) RegisterProxy(_ sdk.Context, principal sdk.ValAddress, proxy sdk.AccAddress) error {
	b.val2Proxy[principal.String()] = proxy
	b.proxy2Val[proxy.String()] = principal
	return nil
}

// GetPrincipal returns the principal validator associated with the given proxy
func (b Broadcaster) GetPrincipal(_ sdk.Context, proxy sdk.AccAddress) sdk.ValAddress {
	return b.proxy2Val[proxy.String()]
}

// GetProxy returns the proxy associated with the given principal validator
func (b Broadcaster) GetProxy(_ sdk.Context, principal sdk.ValAddress) sdk.AccAddress {
	return b.val2Proxy[principal.String()]
}

// GetProxyCount returns the proxy count
func (b Broadcaster) GetProxyCount(_ sdk.Context) uint32 {
	return uint32(len(b.val2Proxy))
}

// GetLocalPrincipal returns the principal validator associated with this broadcaster instance
func (b Broadcaster) GetLocalPrincipal(_ sdk.Context) sdk.ValAddress {
	return b.LocalPrincipal
}
