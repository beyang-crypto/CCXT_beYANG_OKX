package rest

//https://www.okx.com/docs-v5/en/#rest-api-account-get-balance
type WalletBalance struct {
	Code string              `json:"code"`
	Data []dataWalletBalance `json:"data"`
	Msg  string              `json:"msg"`
}
type detailsWalletBalance struct {
	AvailBal      string `json:"availBal"`
	AvailEq       string `json:"availEq"`
	CashBal       string `json:"cashBal"`
	Ccy           string `json:"ccy"`
	CrossLiab     string `json:"crossLiab"`
	DisEq         string `json:"disEq"`
	Eq            string `json:"eq"`
	EqUsd         string `json:"eqUsd"`
	FrozenBal     string `json:"frozenBal"`
	Interest      string `json:"interest"`
	IsoEq         string `json:"isoEq"`
	IsoLiab       string `json:"isoLiab"`
	IsoUpl        string `json:"isoUpl"`
	Liab          string `json:"liab"`
	MaxLoan       string `json:"maxLoan"`
	MgnRatio      string `json:"mgnRatio"`
	NotionalLever string `json:"notionalLever"`
	OrdFrozen     string `json:"ordFrozen"`
	Twap          string `json:"twap"`
	UTime         string `json:"uTime"`
	Upl           string `json:"upl"`
	UplLiab       string `json:"uplLiab"`
	StgyEq        string `json:"stgyEq"`
	SpotInUseAmt  string `json:"spotInUseAmt"`
}
type dataWalletBalance struct {
	AdjEq       string                 `json:"adjEq"`
	Details     []detailsWalletBalance `json:"details"`
	Imr         string                 `json:"imr"`
	IsoEq       string                 `json:"isoEq"`
	MgnRatio    string                 `json:"mgnRatio"`
	Mmr         string                 `json:"mmr"`
	NotionalUsd string                 `json:"notionalUsd"`
	OrdFroz     string                 `json:"ordFroz"`
	TotalEq     string                 `json:"totalEq"`
	UTime       string                 `json:"uTime"`
}

func OKXToWalletBalance(data interface{}) WalletBalance {
	bt, _ := data.(WalletBalance)
	return bt
}
