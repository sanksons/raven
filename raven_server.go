package raven

import (
	"github.com/gin-gonic/gin"
)

func StartServer(receiver *RavenReceiver) error {

	receiverHolder := &ReceiverHolder{receiver}
	r := gin.Default()
	r.GET("/stats", receiverHolder.stats)
	//r.POST("/flushAll", receiverHolder.flushAll)
	r.POST("/flushDead", receiverHolder.flushDeadQ)

	r.Run(":5656")
	return nil
}

type ReceiverHolder struct {
	receiver *RavenReceiver
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
