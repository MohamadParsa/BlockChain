package wallet_server

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/MohamadParsa/BlockChain/v1/wallet"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	TEMPLATE_FOLDER = "./v1/wallet/wallet_server/pages/"
	INDEX_TEMPLATE  = "index.html"
)

type RestFull struct {
	wallet *wallet.Wallet
}

func New(wallet *wallet.Wallet) *RestFull {
	return &RestFull{wallet: wallet}
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
	router.GET("/index", restFull.index)
}

func (restFull RestFull) index(c *gin.Context) {
	html, err := createContentByWalletTemplate(restFull.wallet, INDEX_TEMPLATE)
	if err != nil {
		log.Error("error parsing template", err)
		c.JSON(500, gin.H{"result": "internal error"})
	} else {
		c.Data(200, "text/html", []byte(html))
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
