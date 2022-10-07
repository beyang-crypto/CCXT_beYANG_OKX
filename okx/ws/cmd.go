package ws

//	Необходим для удобного создания подписо
type Cmd struct {
	Op   string    `json:"op"`
	Args []ArgsCmd `json:"args"`
}
type ArgsCmd struct {
	Channel  string `json:"channel"`
	InstType string `json:"instType"`
	Uly      string `json:"uly"`
	InstID   string `json:"instId"`
	Ccy      string `json:"ccy"`
	AlgoId   string `json:"algoId"`
}

//	https://www.okx.com/docs-v5/en/#websocket-api-login
type Auth struct {
	Op   string     `json:"op"`
	Args []ArgsAuth `json:"args"`
}
type ArgsAuth struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}
