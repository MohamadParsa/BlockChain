package miner_server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/MohamadParsa/BlockChain/v1/transaction/transaction_request"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ERR_INPUT_INVALID = "input invalid"
)

type RestFull struct {
	miner *miner.Miner
}

func New(miner *miner.Miner) *RestFull {
	return &RestFull{miner: miner}
}
func (restFull RestFull) Serve(port string) {
	router := gin.New()

	setHealthMethod(router)

	routerV1 := router.Group("/v1")
	restFull.setAPIMethodsV1(routerV1)

	log.Printf("server listening at %v \n", port)

	log.Fatal(http.ListenAndServe(port, router))
}
func (restFull RestFull) setAPIMethodsV1(router *gin.RouterGroup) {
	router.GET("/blockChain", restFull.blockChain)
	router.POST("/Transaction", restFull.transaction)
}

func (restFull RestFull) blockChain(c *gin.Context) {
	jsonByteResult, err := json.Marshal(restFull.miner.BlockChain())
	writeResponse(c, jsonByteResult, err)
}
func (restFull RestFull) transaction(c *gin.Context) {
	transactionRequest, ok := extractTransactionRequest(c.Request.Body)
	if ok {
		ok, err := restFull.miner.AddTransaction(transactionRequest)
		jsonByteResult, _ := json.Marshal(ok)
		writeResponse(c, jsonByteResult, err)
	} else {
		writeResponse(c, nil, errors.New(ERR_INPUT_INVALID))
	}
}

func writeResponse(c *gin.Context, jsonByteResult []byte, err error) {
	if err != nil {
		c.JSON(500, gin.H{"result": "internal error"})
		return
	}
	c.Data(200, "application/json", jsonByteResult)
}
func setHealthMethod(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}
func extractTransactionRequest(b io.Reader) (*transaction_request.TransactionRequest, bool) {
	decoder := json.NewDecoder(b)
	var transactionRequest transaction_request.TransactionRequest
	err := decoder.Decode(&transactionRequest)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to extract description request information from body"))
		return nil, false
	}

	return &transactionRequest, isTransactionRequestValid(&transactionRequest)
}
func isTransactionRequestValid(transactionRequest *transaction_request.TransactionRequest) bool {
	result := true
	if transactionRequest.RecipientAddress() == "" || transactionRequest.SenderAddress() == "" || transactionRequest.Value() < 0 {
		result = false
	}

	return result
}
