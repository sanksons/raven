package raven

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer(receiver *RavenReceiver) error {

	fmt.Println("Starting Server ...")
	//make sure we are running in release mode.
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	//r := gin.Default()

	receiverHolder := &ReceiverHolder{receiver, r}
	//Define routes
	receiverHolder.defineRoutes()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		//fmt.Println("got interrupt")
		StopServer(receiver)
		//	time.Sleep(time.Second * 10)
		os.Exit(0)
	}()
	if err := receiverHolder.startListening(); err != nil {
		StopServer(receiver)
		return err
	}
	return nil
}

func StopServer(receiver *RavenReceiver) error {
	fmt.Println("Stopping Server...")
	fmt.Println("######################################")
	receiver.Stop()
	fmt.Println("######################################")
	return nil
}

type ReceiverHolder struct {
	receiver *RavenReceiver
	engine   *gin.Engine
}

func (this *ReceiverHolder) defineRoutes() {

	this.engine.GET("/", this.ping)
	this.engine.GET("/ping", this.ping)
	this.engine.GET("/stats", this.stats)
	this.engine.GET("/showDeadBox", this.showDeadBox)

	//kill receiver/restart
	//show dead messages.
	this.engine.POST("/flushDead", this.flushDeadQ)
	this.engine.POST("/flushAll", this.flushAll)
}

func (this *ReceiverHolder) startListening() error {

	var port string = "0"
	if this.receiver.port != "" {
		port = this.receiver.port
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	this.receiver.port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	this.receiver.ShowMessage()

	s := http.Server{}
	s.Handler = this.engine
	return s.Serve(listener)
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

func (this *ReceiverHolder) flushAll(c *gin.Context) {
	responsedata := this.receiver.FlushAll()
	data := responsedata
	c.JSON(200, data)
}

func (this *ReceiverHolder) ping(c *gin.Context) {
	c.JSON(200, "OK")
}

func (this *ReceiverHolder) showDeadBox(c *gin.Context) {
	msgs, err := this.receiver.ShowDeadBox()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, msgs)
}
