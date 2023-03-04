package wallet_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/MohamadParsa/BlockChain/v1/wallet"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	TEMPLATE_FOLDER = "./v1/wallet/wallet_server/pages/"
	INDEX_TEMPLATE  = "index.html"
)

type Server struct {
	wallet *wallet.Wallet
}

func New(wallet *wallet.Wallet) *Server {
	return &Server{wallet: wallet}
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
func (server Server) sendCrypto(c *gin.Context) {
	if sendCryptoRequest, ok := extractSendCryptoData(c.Request.Body, c.Request.Header); ok {
		tran, err := server.wallet.SendCrypto(sendCryptoRequest.RecipientAddress, stringTofloat64(sendCryptoRequest.Amount))
		if err != nil {
			fmt.Println(tran)
		}
		message := "success"
		if err != nil {
			message = "failed"
		}

		writeResponse(c, []byte(`{"result":"`+message+`"}`), err)
	} else {

		c.JSON(500, gin.H{"result": "error in parameters value"})

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
func extractSendCryptoData(b io.Reader, h http.Header) (*SendCryptoRequest, bool) {
	decoder := json.NewDecoder(b)

	fmt.Println(b)
	var sendCryptoRequest SendCryptoRequest
	err := decoder.Decode(&sendCryptoRequest)
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
