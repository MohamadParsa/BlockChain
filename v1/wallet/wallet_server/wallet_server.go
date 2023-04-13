package wallet_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"text/template"

	"github.com/MohamadParsa/BlockChain/v1/wallet"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	TEMPLATE_FOLDER                = "./v1/wallet/wallet_server/pages/"
	INDEX_TEMPLATE                 = "index.html"
	ADD_TRANSACTION_METHOD_ADDRESS = "/v1/AddTransaction"
	AMOUNT_METHOD_ADDRESS          = "/v1/amount"
)

type Server struct {
	wallet                   *wallet.Wallet
	blockChainServersAddress map[string]string
}

func New(wallet *wallet.Wallet) *Server {
	blockChainServersAddress := make(map[string]string)
	blockChainServersAddress["main"] = "http://localhost:8080"
	return &Server{wallet: wallet, blockChainServersAddress: blockChainServersAddress}
}
func (server Server) Serve(port string) {
	router := gin.New()

	setHealthMethod(router)

	routerV1 := router.Group("/v1")
	server.setAPIMethodsV1(routerV1)

	log.Printf("server listening at %v \n", port)

	log.Fatal(http.ListenAndServe(port, router))
}
func (server Server) setAPIMethodsV1(router *gin.RouterGroup) {
	router.GET("/index", server.index)
	router.GET("/balance", server.balance)
	router.POST("/sendCrypto", server.sendCrypto)

}

func (server Server) index(c *gin.Context) {
	html, err := createContentByWalletTemplate(server.wallet, INDEX_TEMPLATE)
	if err != nil {
		log.Error("error parsing template", err)
		c.JSON(500, gin.H{"result": "internal error"})
	} else {
		c.Data(200, "text/html", []byte(html))
	}
}
func (server Server) balance(c *gin.Context) {
	responseLog := make(map[string]float64, len(server.blockChainServersAddress))
	waitGroup := sync.WaitGroup{}

	for _, serverAddress := range server.blockChainServersAddress {
		waitGroup.Add(1)
		responseLog[serverAddress] = 0
		go func() {
			response, err := http.Get(serverAddress + AMOUNT_METHOD_ADDRESS + "/" + server.wallet.Address())
			if err == nil {
				res, _ := ioutil.ReadAll(response.Body)
				log.Println("res", string(res))
				amount, err := strconv.ParseFloat(string(res), 64)
				if err == nil {
					responseLog[serverAddress] = amount
				}
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	mostRepeatedResponse := 0.0
	responseMap := make(map[float64]int, len(server.blockChainServersAddress))

	for _, response := range responseLog {
		responseMap[response] = +1
		if mostRepeatedResponse == 0.0 ||
			responseMap[response] > responseMap[mostRepeatedResponse] {
			mostRepeatedResponse = response
		}

	}
	c.Data(200, "text/html", []byte(strconv.FormatFloat(mostRepeatedResponse, 'f', -1, 64)))

}
func (server Server) sendCrypto(c *gin.Context) {
	responseLog := make(map[string]int, len(server.blockChainServersAddress))
	waitGroup := sync.WaitGroup{}
	if sendCryptoRequest, ok := extractSendCryptoData(c.Request.Body, c.Request.Header); ok {
		transactionDTO, err := server.wallet.SendCrypto(sendCryptoRequest.RecipientAddress, stringTofloat64(sendCryptoRequest.Amount))
		if err == nil {
			jsonByte, err := json.Marshal(transactionDTO)
			jsonBuf := bytes.NewBuffer(jsonByte)
			if err == nil {

				for _, serverAddress := range server.blockChainServersAddress {
					waitGroup.Add(1)
					responseLog[serverAddress] = 0
					go func() {
						response, err := http.Post(serverAddress+ADD_TRANSACTION_METHOD_ADDRESS, "application/json", jsonBuf)
						if err == nil {
							res, _ := ioutil.ReadAll(response.Body)
							log.Println("res", string(res))
							responseLog[serverAddress] = response.StatusCode
						}
						waitGroup.Done()
					}()
				}
				waitGroup.Wait()
			}

		}
		successResponseCount := 0
		for _, response := range responseLog {
			if response == 200 {
				successResponseCount++
			}
		}

		responseStatusCode := 500
		if successResponseCount >= (len(server.blockChainServersAddress)/2)+1 {
			responseStatusCode = 200
		}
		writeDefaultResponse(c, responseStatusCode)
	} else {

		c.JSON(500, gin.H{"result": "error in parameters value"})

	}
}
func writeDefaultResponse(c *gin.Context, statusCode int) {
	fmt.Println(statusCode)
	message := "success"
	if statusCode != 200 {
		message = "failed"
	}
	c.Data(statusCode, "application/json", []byte("\"result\":\""+message+"\""))
}
func setHealthMethod(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}

func createContentByWalletTemplate(wallet *wallet.Wallet, templateName string) (string, error) {
	contentTemplate := template.New("tmp")
	contentTemplate, err := contentTemplate.ParseFiles(TEMPLATE_FOLDER + templateName)

	if err != nil {
		log.Error("error parsing template", err)
		return "", err
	}
	var byteBuffer bytes.Buffer
	err = contentTemplate.ExecuteTemplate(&byteBuffer, templateName, wallet)
	if err != nil {
		log.Error("error execute template", err)
		return "", err
	}
	return byteBuffer.String(), nil
}
func extractSendCryptoData(b io.Reader, h http.Header) (*SendCryptoRequestDTO, bool) {
	data, _ := io.ReadAll(b)

	var sendCryptoRequest SendCryptoRequestDTO
	err := json.Unmarshal(data, &sendCryptoRequest)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to extract sendCrypto request information from body"))
		return nil, false
	}
	if stringTofloat64(sendCryptoRequest.Amount) <= 0 || sendCryptoRequest.RecipientAddress == "" {
		log.Error("failed to extract sendCrypto request required information from body")
		return nil, false
	}
	return &sendCryptoRequest, true
}
func stringTofloat64(num string) float64 {
	val, _ := strconv.ParseFloat(num, 64)
	return val
}
