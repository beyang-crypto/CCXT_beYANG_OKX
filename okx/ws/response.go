package ws

//	https://www.okx.com/docs-v5/en/#websocket-api-public-channel-instruments-channel
type Tickers struct {
	Arg  argTickers    `json:"arg"`
	Data []dataTickers `json:"data"`
}
type argTickers struct {
	Channel string `json:"channel"`
	InstID  string `json:"instId"`
}
type dataTickers struct {
	InstType  string `json:"instType"`
	InstID    string `json:"instId"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24H   string `json:"open24h"`
	High24H   string `json:"high24h"`
	Low24H    string `json:"low24h"`
	VolCcy24H string `json:"volCcy24h"`
	Vol24H    string `json:"vol24h"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
	Ts        string `json:"ts"`
}

//	https://www.okx.com/docs-v5/en/#websocket-api-private-channel-balance-and-position-channel
type WalletBalance struct {
	Arg  argWalletBalance    `json:"arg"`
	Data []dataWalletBalance `json:"data"`
}
type argWalletBalance struct {
	Channel string `json:"channel"`
	UID     string `json:"uid"`
}
type balDataWalletBalance struct {
	Ccy     string `json:"ccy"`
	CashBal string `json:"cashBal"`
	UTime   string `json:"uTime"`
}
type posDataWalletBalance struct {
	PosID    string `json:"posId"`
	TradeID  string `json:"tradeId"`
	InstID   string `json:"instId"`
	InstType string `json:"instType"`
	MgnMode  string `json:"mgnMode"`
	PosSide  string `json:"posSide"`
	Pos      string `json:"pos"`
	Ccy      string `json:"ccy"`
	PosCcy   string `json:"posCcy"`
	AvgPx    string `json:"avgPx"`
	UTIme    string `json:"uTIme"`
}
type dataWalletBalance struct {
	PTime     string                 `json:"pTime"`
	EventType string                 `json:"eventType"`
	BalData   []balDataWalletBalance `json:"balData"`
	PosData   []posDataWalletBalance `json:"posData"`
}
