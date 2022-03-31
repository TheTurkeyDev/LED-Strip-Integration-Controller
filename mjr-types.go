package main

type MJRRequest struct {
	Type      string   `json:"type"`
	Nonce     string   `json:"nonce"`
	ChannelId int      `json:"channel_id"`
	Topics    []string `json:"topics"`
	Token     string   `json:"token"`
}

type MJRResponse struct {
	Type      string `json:"type"`
	Nonce     string `json:"nonce"`
	ChannelId int    `json:"channel_id"`
	Topic     string `json:"topic"`
}

type ChannePointRedeemResponse struct {
	MJRResponse
	Message ChannePointRedeem `json:"message"`
}

type ChannePointRedeem struct {
	Redemption ChannePointRedemption `json:"redemption"`
}

type ChannePointRedemption struct {
	Reward    ChannePointReward `json:"reward"`
	UserInput string            `json:"user_input"`
}

type ChannePointReward struct {
	Id string `json:"id"`
}
