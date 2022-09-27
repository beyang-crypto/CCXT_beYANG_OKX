package ws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/buger/jsonparser"      //  Для вытаскивания одного значения из файла json
	"github.com/chuckpreslar/emission" // Эмитер необходим для удобного выполнения функции в какой-то момент
	"github.com/goccy/go-json"         // для создания собственных json файлов и преобразования json в структуру
	"github.com/gorilla/websocket"
)

const (
	HostPublicWebSocket     = "wss://ws.okx.com:8443/ws/v5/public"
	HostPrivateWebSocket    = "wss://ws.okx.com:8443/ws/v5/private"
	HostPublicWebSocketAWS  = "wss://wsaws.okx.com:8443/ws/v5/public"
	HostPrivateWebSocketAWS = "wss://wsaws.okx.com:8443/ws/v5/private"
)

const (
	ChannelBalanceAndPosition = "balance_and_position"
	ChannelRfqs               = "rfqs"
	ChannelQuotes             = "quotes"
	ChannelTicker             = "tickers"
	ChannelOptSummary         = "opt-summary"
	ChannelAccount            = "account"
	ChannelAccountGreeks      = "account-greeks"
	ChannelGridPositions      = "grid-positions"
	ChannelGridSubOrders      = "grid-sub-orders"
	ChannelLiquidationWarning = "liquidation-warning"
	ChannelGridOrdersSpot     = "grid-orders-spot"
	ChannelGridOrdersContract = "grid-orders-contract"
	ChannelGridOrdersMoon     = "grid-orders-moon"
	ChannelAlgoAdvance        = "algo-advance"
	ChannelEstimatedPrice     = "estimated-price"
	ChannelPositions          = "positions"
	ChannelOrders             = "orders"
	ChannelOrdersAlgo         = "orders-algo"
)

const (
	InstTypeSpot    = "SPOT"
	InstTypeMargin  = "MARGIN"
	InstTypeSwap    = "SWAP"
	InstTypeFutures = "FUTURES"
	InstTypeOption  = "OPTION"
	InstTypeAny     = "ANY"
)

type Configuration struct {
	Addr           string `json:"addr"`
	ApiKey         string `json:"api_key"`
	SecretKey      string `json:"secret_key"`
	APIKeyPassword string `json:"passphrase"`
	DebugMode      bool   `json:"debug_mode"`
}

type OKXWS struct {
	cfg  *Configuration
	conn *websocket.Conn

	mu            sync.RWMutex
	subscribeCmds []Cmd //	сохраняем все подписки у данной биржи, чтоб при переподключении можно было к ним повторно подключиться

	emitter *emission.Emitter
}

func (b *OKXWS) GetPair(coin1 string, coin2 string) string {
	return coin1 + "-" + coin2
}

func New(config *Configuration) *OKXWS {

	// 	потом тут добавятся различные другие настройки
	b := &OKXWS{
		cfg:     config,
		emitter: emission.NewEmitter(),
	}
	return b
}

// 	подписка только с channel
func (b *OKXWS) Subscribe1(channel string) {

	var args []ArgsCmd

	args = append(args, ArgsCmd{
		Channel: channel,
	})

	b.subscribe((args))
}

// 	подписка c 2 аргументами
func (b *OKXWS) Subscribe2(channel string, secondArg string) {
	var args []ArgsCmd
	switch channel {
	case "opt-summary":
		args = append(args, ArgsCmd{
			Channel: channel,
			Uly:     secondArg,
		})
	case "account", "account-greeks":
		args = append(args, ArgsCmd{
			Channel: channel,
			Ccy:     secondArg,
		})
	case "grid-positions", "grid-sub-orders":
		args = append(args, ArgsCmd{
			Channel: channel,
			AlgoId:  secondArg,
		})
	case "liquidation-warning", "grid-orders-spot", "grid-orders-contract", "grid-orders-moon":
		args = append(args, ArgsCmd{
			Channel:  channel,
			InstType: secondArg,
		})
	default:
		args = append(args, ArgsCmd{
			Channel: channel,
			InstID:  secondArg,
		})
	}

	b.subscribe((args))
}

// 	подписка с тремя аргументами
func (b *OKXWS) Subscribe3(channel string, arg1 string, arg2 string) {

	var args []ArgsCmd
	switch channel {
	case "algo-advance":
		args = append(args, ArgsCmd{
			Channel:  channel,
			InstType: arg1,
			InstID:   arg2,
		})
	case "estimated-price":
		args = append(args, ArgsCmd{
			Channel:  channel,
			InstType: arg1,
			Uly:      arg2,
		})
	}

	b.subscribe((args))
}

// 	подписка с четырьмя аргументами
func (b *OKXWS) Subscribe4(channel string, arg1 string, arg2 string, arg3 string) {

	var args []ArgsCmd
	switch channel {
	case "positions", "orders", "orders-algo":
		args = append(args, ArgsCmd{
			Channel:  channel,
			InstType: arg1,
			Uly:      arg2,
			InstID:   arg3,
		})
	}

	b.subscribe((args))
}

func (b *OKXWS) subscribe(args []ArgsCmd) {
	cmd := Cmd{
		Op:   "subscribe",
		Args: args,
	}
	b.subscribeCmds = append(b.subscribeCmds, cmd)
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 1")
	}
	b.SendCmd(cmd)
}

//	отправка команды на сервер в отдельной функции для того, чтобы при переподключении быстро подписаться на все предыдущие каналы
func (b *OKXWS) SendCmd(cmd Cmd) {
	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_OKX/okx",
				"File": "client.go",
				"Functions" : "(b *OKXWS) sendCmd(cmd Cmd)",
				"Function where err" : "json.Marshal",
				"Exchange" : "OKX",
				"Data" : [%s],
				"Error" : %s
			}`, cmd, err)
		log.Fatal()
	}
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку part 2")
	}
	b.Send(string(data))
}

func (b *OKXWS) Send(msg string) (err error) {
	defer func() {
		// recover необходим для корректной обработки паники
		if r := recover(); r != nil {
			if err != nil {
				log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_OKX/okx",
						"File": "client.go",
						"Functions" : "(b *OKXWS) Send(msg string) (err error)",
						"Function where err" : "b.conn.WriteMessage",
						"Exchange" : "OKX",
						"Data" : [websocket.TextMessage, %s],
						"Error" : %s,
						"Recover" : %v
					}`, msg, err, r)
				log.Fatal()
			}
			err = errors.New(fmt.Sprintf("OKXWs send error: %v", r))
		}
	}()
	if b.cfg.DebugMode {
		log.Printf("Отправка сообщения на сервер. текст сообщения:%s", msg)
	}

	err = b.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	return
}

// подключение к серверу и постоянное чтение приходящих ответов
func (b *OKXWS) Start() error {
	if b.cfg.DebugMode {
		log.Printf("Начало подключения к серверу")
	}
	b.connect()

	cancel := make(chan struct{})

	go func() {
		t := time.NewTicker(time.Second * 15)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				b.ping()
			case <-cancel:
				return
			}
		}
	}()

	go func() {
		defer close(cancel)

		for {
			_, data, err := b.conn.ReadMessage()
			if err != nil {

				if websocket.IsCloseError(err, 1006) {
					b.closeAndReconnect()
					//Необходим вызв SubscribeToTicker в отдельной горутине, рекурсия, думаю, тут неуместна
					log.Printf("Status: INFO	ошибка 1006 начинается переподключение к серверу")

				} else {
					b.conn.Close()
					log.Printf(`
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_OKX/okx",
						"File": "client.go",
						"Functions" : "(b *OKXWS) Start() error",
						"Function where err" : "b.conn.ReadMessage",
						"Exchange" : "OKX",
						"Error" : %s
					}`, err)
					log.Fatal()
				}

			} else {
				b.messageHandler(data)
			}
		}
	}()

	return nil
}

//	Необходим для приватных каналов
func (b *OKXWS) Auth() {
	ts := time.Now().Unix()
	sign := fmt.Sprintf("%d", ts) + "GET" + "/users/self/verify"
	sig := hmac.New(sha256.New, []byte(b.cfg.SecretKey))
	sig.Write([]byte(sign))
	signature := base64.StdEncoding.EncodeToString(sig.Sum(nil))

	log.Printf(sign)
	log.Printf(signature)

	var args []ArgsAuth
	args = append(args, ArgsAuth{
		APIKey:     b.cfg.ApiKey,
		Passphrase: b.cfg.APIKeyPassword,
		Timestamp:  fmt.Sprintf("%d", ts),
		Sign:       signature,
	})
	auth := Auth{
		Op:   "login",
		Args: args,
	}
	data, err := json.Marshal(auth)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_OKX/okx",
				"File": "client.go",
				"Functions" : "(b *OKXWS) Auth()",
				"Function where err" : "json.Marshal",
				"Exchange" : "OKX",
				"Data" : [%v],
				"Error" : %s
			}`, auth, err)
		log.Fatal()
	}
	if b.cfg.DebugMode {
		log.Printf("Создание json сообщения на подписку")
	}
	b.Send(string(data))
}

func (b *OKXWS) connect() {

	c, _, err := websocket.DefaultDialer.Dial(b.cfg.Addr, nil)
	if err != nil {
		log.Printf(`{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_OKX/okx",
						"File": "client.go",
						"Functions" : "(b *OKXWS) connect()",
						"Function where err" : "websocket.DefaultDialer.Dial",
						"Exchange" : "OKX",
						"Data" : [%s, nil],
						"Error" : %s
					}`, b.cfg.Addr, err)
		log.Fatal()
	}
	b.conn = c
	for _, cmd := range b.subscribeCmds {
		b.SendCmd(cmd)
	}
}

func (b *OKXWS) closeAndReconnect() {
	b.conn.Close()
	b.connect()
}

func (b *OKXWS) ping() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("OKXWs ping error: %v", r)
		}
	}()

	//	https://www.okx.com/docs-v5/en/#websocket-api-connect
	err := b.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
	if err != nil {
		log.Printf("OKXWs ping error: %v", err)
	}
}

func (b *OKXWS) messageHandler(data []byte) {

	if b.cfg.DebugMode {
		log.Printf("OKXWs %v", string(data))
	}

	//	в ошибке нет необходимости, т.к. она выходит каждый раз, когда не найдет элемент
	event, _ := jsonparser.GetString(data, "event")
	channel, _ := jsonparser.GetString(data, "arg", "channel")

	switch event {
	case "login":
		msg, _ := jsonparser.GetString(data, "msg")
		if len(msg) != 0 {
			log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_OKX/okx",
				"File": "client.go",
				"Functions" : "(b *OKXWS) messageHandler(data []byte)",
				"Exchange" : "OKX",
				"Message" : %s
			}`, string(data))
			log.Fatal()
		}
	case "error":
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_OKX/okx",
				"File": "client.go",
				"Functions" : "(b *OKXWS) messageHandler(data []byte)",
				"Exchange" : "OKX",
				"Message" : %s
			}`, string(data))
		log.Fatal()
	case "subscribe":
		//	Ну что сказать, хорошо, что подписались успешно
	case "unsubscribe":
		//
	default:
		switch channel {
		case "tickers":
			instId, _ := jsonparser.GetString(data, "instId")
			var ticker Tickers
			err := json.Unmarshal(data, &ticker)
			if err != nil {
				log.Printf(`Huobi
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_OKX/okx",
						"File": "client.go",
						"Functions" : "(b *OKXWS) messageHandler(data []byte)",
						"Function where err" : "json.Unmarshal",
						"Exchange" : "OKX",
						"Comment" : %s to BookTicker struct,
						"Error" : %s
					}`, string(data), err)
				log.Fatal()
			}
			b.processTickers(instId, ticker)
		case "balance_and_position":
			var walletBalance WalletBalance
			err := json.Unmarshal(data, &walletBalance)
			if err != nil {
				log.Printf(`Huobi
					{
						"Status" : "Error",
						"Path to file" : "CCXT_BEYANG_OKX/okx",
						"File": "client.go",
						"Functions" : "(b *OKXWS) messageHandler(data []byte)",
						"Function where err" : "json.Unmarshal",
						"Exchange" : "OKX",
						"Comment" : %s to BookTicker struct,
						"Error" : %s
					}`, string(data), err)
				log.Fatal()
			}
			b.processWalletBalance(walletBalance)
		default:
			if string(data) == "pong" {

			} else {
				log.Printf(`
				{
					"Status" : "INFO",
					"Path to file" : "CCXT_BEYANG_OKX/okx",
					"File": "client.go",
					"Functions" : "(b *OKXWS) messageHandler(data []byte)",
					"Exchange" : "OKX",
					"Comment" : "не известный ответ от сервера"
					"Message" : %s
				}`, string(data))
				log.Fatal()
			}
		}
	}
}
