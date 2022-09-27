package main

import (
	rest "ccxt_beyang_okx/okx/rest"
	ws "ccxt_beyang_okx/okx/ws"
	"log"
	"time"
)

func main() {
	cfgpriv := &ws.Configuration{
		Addr:           ws.HostPrivateWebSocket,
		ApiKey:         "xxx",
		SecretKey:      "xxx",
		APIKeyPassword: "xxx",
		DebugMode:      true,
	}
	cfg := &ws.Configuration{
		Addr:           ws.HostPublicWebSocket,
		ApiKey:         "xxx",
		SecretKey:      "xxx",
		APIKeyPassword: "xxx",
		DebugMode:      true,
	}
	cfgrest := &ws.Configuration{
		Addr:           ws.HostPrivateWebSocket,
		ApiKey:         "xxx",
		SecretKey:      "xxx",
		APIKeyPassword: "xxx",
		DebugMode:      true,
	}
	r := rest.New((*rest.Configuration)(cfgrest))
	p := ws.New(cfgpriv)
	b := ws.New(cfg)
	p.Start()
	b.Start()

	p.Auth()

	time.Sleep(5 * time.Second)
	p.Subscribe1(ws.ChannelBalanceAndPosition)

	b.On(ws.ChannelBalanceAndPosition, handleWalletBalanceOfSymbol)
	pair := b.GetPair("BTC", "USDT")
	b.Subscribe2(ws.ChannelTicker, pair)

	b.On(ws.ChannelTicker, handleBookTicker)
	b.On(ws.ChannelTicker, handleBestBidPrice)

	go func() {
		time.Sleep(5 * time.Second)
		balance := r.GetBalance(rest.RestURL)
		for _, datas := range balance.Data {
			for _, details := range datas.Details {
				log.Printf("coin = %s, total = %s", details.Ccy, details.CashBal)
			}
		}
	}()

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
