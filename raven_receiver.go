package raven

import (
	"fmt"
	"time"

	"github.com/newrelic/go-agent"

	"github.com/sanksons/gowraps/util"
)

// Initiate a Raven Receiver.
func newRavenReceiver(id string, source Source) (*RavenReceiver, error) {
	rr := new(RavenReceiver)

	rr.setSource(source).setId(id)
	return rr, nil
}

//
// Raven Receiver / Message collector
//
type RavenReceiver struct {
	//Id assigned to this receiver.
	//Will be helpful to keep this unique across receivers.
	id string

	//Source where to look for ravens.
	source Source

	//Options define characteristics of a receiver.
	options struct {
		//Specifies if we want to use reliable Q or not
		//@todo: ordering is yet to be implemented.
		isReliable, ordering bool
	}

	//Q to store processing and dead messages.
	// used only when marked reliable.
	processingQ Q
	deadQ       Q

	// Farm to which reveiver belongs.
	farm *Farm

	startedAt time.Time
}

func (this RavenReceiver) String() string {
	return fmt.Sprintf("id: %s, source: %s , reliable: %v, processingQ: %s, deadQ: %s",
		this.id, this.source.GetName(), this.options.isReliable, this.processingQ.GetName(),
		this.deadQ.GetName(),
	)
}

func (this *RavenReceiver) getNewrelicTransaction() newrelic.Transaction {
	if this.farm.newrelicApp != nil {
		return this.farm.newrelicApp.StartTransaction(this.id, nil, nil)
	}
	return nil
}

func (this *RavenReceiver) endNewrelicTransaction(txn newrelic.Transaction) {
	if txn == nil {
		return
	}
	txn.End()
}

//Record heartbeat of consumer.
func (this *RavenReceiver) recordHeartBeat(inflightCount int) {

	if this.farm.newrelicApp == nil {
		return
	}
	//Record Heart Beat
	this.farm.newrelicApp.RecordCustomEvent(
		fmt.Sprintf("Heartbeat-%s", this.id), map[string]interface{}{
			"inflightcount": inflightCount,
			"checkedAt":     time.Now(),
			"queue":         this.source.GetName(),
		},
	)

	this.getLogger().Info(this.source.GetName(), this.id, "HeartBeat",
		fmt.Sprintf("In Flight Ravens: %d", inflightCount),
	)
}

//get the logger object.
func (this *RavenReceiver) getLogger() Logger {
	return this.farm.logger
}

// Mark the Q as reliable.
func (this *RavenReceiver) MarkReliable() *RavenReceiver {
	this.options.isReliable = true
	this.defineProcessingQ().defineDeadQ()
	return this
}

// Mark the Q as ordered.
func (this *RavenReceiver) MarkOrdered() *RavenReceiver {
	this.options.ordering = true
	return this
}

func (this *RavenReceiver) setSource(s Source) *RavenReceiver {
	this.source = s
	return this
}

func (this *RavenReceiver) setId(id string) *RavenReceiver {
	this.id = id
	return this
}

func (this *RavenReceiver) defineProcessingQ() *RavenReceiver {

	qname := fmt.Sprintf("%s_processing_%s", this.source.GetRawName(), this.id)
	this.processingQ = createQ(qname, this.source.GetBucket())
	return this
}

func (this *RavenReceiver) defineDeadQ() *RavenReceiver {

	qname := fmt.Sprintf("%s_dead", this.source.GetRawName())
	this.deadQ = createQ(qname, this.source.GetBucket())
	return this
}

//@todo: implement all the necessary validations required for a receiver.
func (this *RavenReceiver) validate() error {
	//Check if Id, Source and farm are defined.
	if this.id == "" {
		return fmt.Errorf("An Id needs to be assigned to Receiver. Make sure its unique within source")
	}
	if this.source.IsEmpty() {
		return fmt.Errorf("Receiver Source cannot be Empty")
	}
	if this.farm == nil {
		return fmt.Errorf("You need to define to which farm this receiver belongs.")
	}
	return nil
}

// Get Messages published but not picked for processing.
func (this *RavenReceiver) GetInFlightRavens() (int, error) {
	return this.farm.manager.InFlightMessages(*this)
}

// Start HeartBeat of Receiver.
func (this *RavenReceiver) StartHeartBeat() error {

	for {
		func() {
			// Incase of panic, restart for loop.
			defer util.PanicHandler("HeartBeat")

			// Pulse rate
			time.Sleep(10 * time.Second)

			cc, err := this.GetInFlightRavens()
			if err != nil {
				this.getLogger().Error(this.source.GetName(), this.id, "HeartBeat",
					fmt.Sprintf("Error: %s", err.Error()),
				)
				return
			}

			//Check if we can record health.
			this.recordHeartBeat(cc)

		}()
	}
}

func (this *RavenReceiver) Start(f func(m *Message, txn newrelic.Transaction) error) error {

	//Start HeartBeat
	go this.StartHeartBeat()

	return this.start(f)

}

func (this *RavenReceiver) StartServer() error {
	return StartServer(this)
}

func (this *RavenReceiver) start(f func(m *Message, txn newrelic.Transaction) error) error {

	this.startedAt = time.Now()

	this.getLogger().Info(this.source.GetName(), this.id,
		fmt.Sprintf("Starting Raven receiver with config, %s", this),
	)
	if verr := this.validate(); verr != nil {
		return verr
	}
	receiver := *this
	//startup con
	if err := this.farm.manager.PreStartup(receiver); err != nil {
		return err
	}

	// this blocks
	for {
		//this blocks, so no need for wait on empty Q.
		msg, err := this.farm.manager.Receive(receiver)
		if err != nil && err == ErrEmptyQueue {
			//Q is empty, Simply recheck.
			this.getLogger().Info(this.source.GetName(), this.id, "Queue is empty recheck")
			continue
		}
		// Something went wrong.
		if err != nil {
			//add a wait here.
			//log error
			this.getLogger().Error(this.source.GetName(), this.id, fmt.Sprintf("Got Error while receiving. Error:%s",
				err.Error()),
			)

			this.getLogger().Info(this.source.GetName(), this.id, "Waiting for 5 seconds before retrying.")
			time.Sleep(5 * time.Second)
			continue
		}

		this.getLogger().Info(this.source.GetName(), this.id, fmt.Sprintf("Received Message: %s",
			msg),
		)

		//
		// Send Message for processing.
		//
		var execerr error
		var txn newrelic.Transaction
		func() {
			// handle any panics occuring from client code.
			defer func() {
				if r := recover(); r != nil {
					emsg := fmt.Sprintf("Panic Occurred !!! Handled Gracefully \n Message: %s", msg)
					execerr = fmt.Errorf(emsg)
				}
				// Check if transaction is started and needs to be wrapped up.
				this.endNewrelicTransaction(txn)
			}()

			// Send Message for processing.
			// Note: pass newrelic transaction alongside so that client can
			// make use of it and record segments.
			txn = this.getNewrelicTransaction()
			execerr = f(msg, txn)
		}()

		if execerr == nil {
			//free up message from processing Q
			err := this.farm.manager.MarkProcessed(msg, receiver)
			if err != nil {
				this.getLogger().Error(
					this.source.GetName(), this.id,
					fmt.Sprintf("Could Not mark message as processed. Message : %s", msg),
				)
			}
		} else if execerr == ErrTmpFailure {
			this.getLogger().Error(
				this.source.GetName(), this.id,
				fmt.Sprintf("Got temporary error while processing message [%s], requeing it", msg),
			)
			err := this.farm.manager.RequeMessage(*msg, receiver)
			if err != nil {
				this.getLogger().Error(
					this.source.GetName(), this.id,
					fmt.Sprintf("Could Not Reque message. Message : %s", msg),
				)
			}
			//sleep till 5 seconds, before repulling message.
			time.Sleep(5 * time.Second)

		} else {
			// Found a permanent error while processing message.
			this.getLogger().Error(
				this.source.GetName(), this.id,
				fmt.Sprintf(
					"Got permanent error while processing message [%s], Discarding it, Error: %s", msg, execerr.Error(),
				),
			)

			//store in DeadQ
			err := this.farm.manager.MarkFailed(msg, receiver)
			if err != nil {
				this.getLogger().Error(
					this.source.GetName(), this.id,
					fmt.Sprintf("Could Not mark message as dead. Message : %s", msg),
				)
			}
		}

	}
}
