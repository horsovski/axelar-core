package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewLinkRequest - LinkRequest constructor
func NewLinkRequest(sender sdk.AccAddress, recipientAddr string, recipientChain string) *LinkRequest {
	return &LinkRequest{
		Sender:         sender,
		RecipientAddr:  recipientAddr,
		RecipientChain: recipientChain,
	}
}

// Route returns the route for this message
func (m LinkRequest) Route() string {
	return RouterKey
}

// Type returns the type of the message
func (m LinkRequest) Type() string {
	return "Link"
}

// ValidateBasic executes a stateless message validation
func (m LinkRequest) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, sdkerrors.Wrap(err, "sender").Error())
	}

	if m.RecipientAddr == "" {
		return fmt.Errorf("missing recipient address")
	}
	if m.RecipientChain == "" {
		return fmt.Errorf("missing recipient chain")
	}
	return nil
}

// GetSignBytes returns the message bytes that need to be signed
func (m LinkRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the set of signers for this message
func (m LinkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}
