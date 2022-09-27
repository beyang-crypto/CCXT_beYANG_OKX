package ws

func (b *OKXWS) processTickers(symbol string, data Tickers) {
	b.Emit(ChannelTicker, symbol, data)
}

func (b *OKXWS) processWalletBalance(data WalletBalance) {
	b.Emit(ChannelBalanceAndPosition, data)
}
