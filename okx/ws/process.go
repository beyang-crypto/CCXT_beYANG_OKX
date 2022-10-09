package ws

func (b *OKXWS) processTickers(name string, symbol string, data Tickers) {
	b.Emit(ChannelTicker, name, symbol, data)
}

func (b *OKXWS) processWalletBalance(name string, data WalletBalance) {
	b.Emit(ChannelBalanceAndPosition, name, data)
}
