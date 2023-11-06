package v1

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = &MsgInit{}
	_ sdk.Msg = &MsgExecute{}
)

func (m *MsgInit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

func (m *MsgInit) ValidateBasic() error {
	return nil
}

func (m *MsgExecute) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

func (m *MsgExecute) ValidateBasic() error {
	return nil
}
