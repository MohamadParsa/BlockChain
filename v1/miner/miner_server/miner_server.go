package miner_server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MohamadParsa/BlockChain/v1/miner"
	transaction_request "github.com/MohamadParsa/BlockChain/v1/transaction/transactionDTO"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ERR_INPUT_INVALID = "input invalid"
	MINING_PERIOD     = time.Minute / 4
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
	restFull.runMiningCore()

	log.Printf("server listening at %v \n", port)

	log.Fatal(http.ListenAndServe(port, router))
}
func (restFull RestFull) runMiningCore() {
	go func() {
		for {
			err := restFull.miner.Mining()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(MINING_PERIOD)
		}
	}()
}
func (restFull RestFull) setAPIMethodsV1(router *gin.RouterGroup) {
	router.GET("/blockChain", restFull.blockChain)
	router.GET("/amount/:walletAddress", restFull.amount)
	router.POST("/AddTransaction", restFull.transaction)
}

func (restFull RestFull) blockChain(c *gin.Context) {
	jsonByteResult, err := json.Marshal(restFull.miner.BlockChain())
	writeResponse(c, jsonByteResult, err)
}
func (restFull RestFull) amount(c *gin.Context) {
	walletAddress := c.Param("walletAddress")
	jsonByteResult, err := json.Marshal(restFull.miner.BlockChain().CalculateTotalAmount(walletAddress))
	writeResponse(c, jsonByteResult, err)
}
func (restFull RestFull) transaction(c *gin.Context) {
	transactionDTO, ok := extractTransactionDTO(c.Request.Body)

	if ok {
		ok, err := restFull.miner.AddTransaction(transactionDTO)
		fmt.Println(ok, err)

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
func extractTransactionDTO(b io.Reader) (*transaction_request.TransactionDTO, bool) {
	data, _ := io.ReadAll(b)
	b.Read(data)
	var transactionDTO transaction_request.TransactionDTO
	err := json.Unmarshal(data, &transactionDTO)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to extract description request information from body"))
		return nil, false
	}
	return &transactionDTO, isTransactionDTOValid(&transactionDTO)
}
func isTransactionDTOValid(transactionDTO *transaction_request.TransactionDTO) bool {
	result := true
	if transactionDTO.RecipientAddress() == "" || transactionDTO.SenderAddress() == "" || transactionDTO.Value() < 0 {
		result = false
	}

	return result
}
