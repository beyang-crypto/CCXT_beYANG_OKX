package main

import (
	"log"
	"time"

	//rest "github.com/TestingAccMar/CCXT_beYANG_OKX/okx/rest"
	ws "github.com/TestingAccMar/CCXT_beYANG_OKX/okx/ws"
)

func main() {
	cfg := &ws.Configuration{
		Addr:           ws.HostPublicWebSocket,
		ApiKey:         "xxx",
		SecretKey:      "xxx",
		APIKeyPassword: "xxx",
		DebugMode:      true,
	}
	b := ws.New(cfg)
	b.Start()

	time.Sleep(5 * time.Second)

	pair := b.GetPair("BTC", "USDT")
	b.Subscribe2(ws.ChannelTicker, pair)

	b.On(ws.ChannelTicker, handleBookTicker)
	b.On(ws.ChannelTicker, handleBestBidPrice)

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	balance := r.GetBalance(rest.RestURL)
	// 	for _, datas := range balance.Data {
	// 		for _, details := range datas.Details {
	// 			log.Printf("coin = %s, total = %s", details.Ccy, details.CashBal)
	// 		}
	// 	}
	// }()

	//	не дает прекратить работу программы
	forever := make(chan struct{})
	<-forever
}

func handleBookTicker(symbol string, data ws.Tickers) {
	log.Printf("OKX Ticker  %s: %v", symbol, data)
}

func handleBestBidPrice(symbol string, data ws.Tickers) {
	log.Printf("OKX BookTicker  %s: BestBidPrice : %s", symbol, data.Data[0].BidPx)
}

func handleWalletBalance(data ws.WalletBalance) {
	log.Printf("OKX WalletBalance: %v", data)
}

func handleWalletBalanceOfSymbol(data ws.WalletBalance) {

	for _, dd := range data.Data {
		for _, bd := range dd.BalData {
			log.Printf("%s:  %s", bd.Ccy, bd.CashBal)
		}
	}
}
