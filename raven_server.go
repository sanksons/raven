package raven

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer(receiver *RavenReceiver) error {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	receiverHolder := &ReceiverHolder{receiver, r}
	//Define routes
	receiverHolder.defineRoutes()

	return receiverHolder.startListening()
}

type ReceiverHolder struct {
	receiver *RavenReceiver
	engine   *gin.Engine
}

func (this *ReceiverHolder) defineRoutes() {

	this.engine.GET("/stats", this.stats)
	//r.POST("/flushAll", receiverHolder.flushAll)
	this.engine.POST("/flushDead", this.flushDeadQ)
}

func (this *ReceiverHolder) startListening() error {

	if this.receiver.port != "" {
		return this.engine.Run(fmt.Sprintf(":%s", this.receiver.port))
	}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	this.receiver.port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	this.receiver.ShowMessage()
	return http.Serve(listener, this.engine)
}

func (this *ReceiverHolder) stats(c *gin.Context) {

	flightData := this.receiver.GetInFlightRavens()

	deadBoxData := this.receiver.GetDeadBoxCount()
	boxes := make([]string, 0)
	for _, box := range this.receiver.msgReceivers {
		boxes = append(boxes, box.id)
	}

	data := gin.H{
		"Queue":      this.receiver.source.GetName(),
		"IsReliable": this.receiver.options.isReliable,
		"Boxes":      boxes,
		"Inflight":   flightData,
		"DeadBox":    deadBoxData,
	}
	c.JSON(200, data)
}

func (this *ReceiverHolder) flushDeadQ(c *gin.Context) {

	responsedata := this.receiver.FlushDeadBox()
	data := responsedata
	c.JSON(200, data)
}

// func (this *ReceiverHolder) flushAll(c *gin.Context) {

// 	data := gin.H{
// 		"success": "OK",
// 	}
// 	c.JSON(200, data)
// }
