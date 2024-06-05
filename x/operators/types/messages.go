package types

func NewMsgRegisterOperator(creator string, moniker string, website string, pictureURL string) *MsgRegisterOperator {
	return &MsgRegisterOperator{
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
		Sender:     creator,
	}
}
