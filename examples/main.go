package main

import (
	"log"
	"time"

	ws "github.com/TestingAccMar/CCXT_beYANG_OKX/okx/ws"
)

func main() {
	cfg := &ws.Configuration{
		Addr:           ws.HostPublicWebSocket,
		ApiKey:         "",
		SecretKey:      "",
		APIKeyPassword: "",
		DebugMode:      true,
	}
	b := ws.New(cfg)
	b.Start()

	time.Sleep(5 * time.Second)

	pair1 := b.GetPair("btc", "usdt")
	pair2 := b.GetPair("eth", "usdt")
	b.Subscribe(ws.ChannelTicker, []string{pair1})
	b.Subscribe(ws.ChannelTicker, []string{pair2})

	//b.On(ws.ChannelTicker, handleBookTicker)
	b.On(ws.ChannelTicker, handleBestBidPrice)

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	balance := rest.OKXToWalletBalance(r.GetBalance())
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

func handleBookTicker(name string, symbol string, data ws.Tickers) {
	log.Printf("%s Ticker  %s: %v", name, symbol, data)
}

func handleBestBidPrice(name string, symbol string, data ws.Tickers) {
	log.Printf("%s BookTicker  %s: BestBidPrice : %s", name, symbol, data.Data[0].BidPx)
}

func handleWalletBalance(name string, data ws.WalletBalance) {
	log.Printf("%s WalletBalance: %v", name, data)
}

func handleWalletBalanceOfSymbol(name string, data ws.WalletBalance) {

	for _, dd := range data.Data {
		for _, bd := range dd.BalData {
			log.Printf("%s:  %s", bd.Ccy, bd.CashBal)
		}
	}
}
