package miner_server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MohamadParsa/BlockChain/v1/miner"
	"github.com/gin-gonic/gin"
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
}

func (restFull RestFull) blockChain(c *gin.Context) {
	jsonByteResult, err := json.Marshal(restFull.miner.BlockChain())
	writeResponse(c, jsonByteResult, err)
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
