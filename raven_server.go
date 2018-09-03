package raven

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartServer(receiver *RavenReceiver) error {

	receiverHolder := &ReceiverHolder{receiver}
	r := gin.Default()
	//r.GET("/ping", receiverHolder.ping)
	r.GET("/stats", receiverHolder.stats)
	r.POST("/flushAll", receiverHolder.flushAll)
	r.POST("/flushDead", receiverHolder.flushDeadQ)

	r.Run(":5656")
	return nil
}

type ReceiverHolder struct {
	receiver *RavenReceiver
}

func (this *ReceiverHolder) ping(c *gin.Context) {
	cc, err := this.receiver.GetInFlightRavens()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"success": cc,
		})
	}
}

func (this *ReceiverHolder) stats(c *gin.Context) {

	cc, err := this.receiver.GetInFlightRavens()
	status := "OK"
	if err != nil {
		status = fmt.Sprintf("Not OK, Error: %s", err.Error())
	}

	var deadLen string
	deadLenI, err := this.receiver.farm.manager.ShowDeadQ(*this.receiver)
	if err != nil {
		deadLen = err.Error()
	} else {
		deadLen = fmt.Sprintf("%d", len(deadLenI))
	}
	data := gin.H{
		"Status":     status,
		"Receiver":   this.receiver.id,
		"Qname":      this.receiver.source.GetRawName(),
		"Bucket":     this.receiver.source.GetBucket(),
		"IsReliable": this.receiver.options.isReliable,
		"Inflight":   cc,
		"Uptime":     this.receiver.startedAt,
		"Dead":       deadLen,
	}
	c.JSON(200, data)
}

func (this *ReceiverHolder) flushDeadQ(c *gin.Context) {

	data := gin.H{
		"success": "OK",
	}
	c.JSON(200, data)
}

func (this *ReceiverHolder) flushAll(c *gin.Context) {

	data := gin.H{
		"success": "OK",
	}
	c.JSON(200, data)
}
