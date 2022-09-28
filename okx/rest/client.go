package rest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

// для создания собственных json файлов и преобразования json в структуру

const (
	RestURL    = "https://www.okx.com"
	RestURLAWS = "https://aws.okx.com"
)

type Configuration struct {
	Addr           string `json:"addr"`
	ApiKey         string `json:"api_key"`
	SecretKey      string `json:"secret_key"`
	APIKeyPassword string `json:"passphrase"`
	DebugMode      bool   `json:"debug_mode"`
}

type OKXWS struct {
	cfg *Configuration
}

func New(config *Configuration) *OKXWS {

	// 	потом тут добавятся различные другие настройки
	b := &OKXWS{
		cfg: config,
	}
	return b
}

func (ex *OKXWS) GetBalance() interface{} {
	//	https://docs.ftx.com/#get-balances
	//	получение времяни
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	ts := time.Now().UTC().Format(time.RFC3339)

	apiKey := ex.cfg.ApiKey
	secretKey := ex.cfg.SecretKey
	passphrase := ex.cfg.APIKeyPassword

	url := ex.cfg.Addr + "/api/v5/account/balance"
	signature_payload := fmt.Sprintf("%sGET/api/v5/account/balance", ts)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(signature_payload))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	//	реализация метода GET
	req, err := http.NewRequest("GET", url, nil)

	log.Printf(signature_payload)
	log.Printf(url)
	req.Header.Set("OK-ACCESS-KEY", apiKey)
	req.Header.Set("OK-ACCESS-SIGN", signature)
	req.Header.Set("OK-ACCESS-TIMESTAMP", ts)
	req.Header.Set("OK-ACCESS-PASSPHRASE", passphrase)
	//	код для вывода полученных данных
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Do(req)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if ex.cfg.DebugMode {
		log.Printf("OKXlletBalance %v", string(data))
	}

	var walletBalance WalletBalance
	err = json.Unmarshal(data, &walletBalance)
	if err != nil {
		log.Printf(`
			{
				"Status" : "Error",
				"Path to file" : "CCXT_BEYANG_OKX/okx/rest",
				"File": "client.go",
				"Functions" : "(ex *OKXWS) GetBalance() (WalletBalance)",
				"Function where err" : "json.Unmarshal",
				"Exchange" : "OKX",
				"Comment" : %s to WalletBalance struct,
				"Error" : %s
			}`, string(data), err)
		log.Fatal()
	}

	return walletBalance

}
