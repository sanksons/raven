package raven

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/sanksons/gowraps/util"

	"github.com/sanksons/raven/childlock"
)

//
// Initiate a Raven Receiver.
//
func newRavenReceiver(id string, source Source) (*RavenReceiver, error) {
	rr := new(RavenReceiver)
	//Define source and Id for receiver.
	rr.setSource(source).setId("")

	// Generate Message Receivers, for each message box
	msgreceivers := make([]*MsgReceiver, 0, len(source.MsgBoxes))
	for _, box := range source.MsgBoxes {
		m := &MsgReceiver{
			msgbox:  box,
			parent:  rr,
			stopped: make(chan bool),
		}
		// Set Id for msgReceiver.
		m.setId(box.GetName())
		msgreceivers = append(msgreceivers, m)
	}
	rr.msgReceivers = msgreceivers

	return rr, nil
}

type RavenReceiver struct {

	//A Unique Id that distinguishes this Receiver from other receivers.
	id string

	// Source from which receiver will fetch Ravens.
	source Source

	// If a PORT is specified raven will be locked to this port
	// else an ephemeral port is picked.
	port string

	// Receiving options.
	options struct {
		//Specifies if we want to use reliable Q or not
		//@todo: ordering is yet to be implemented.
		isReliable, ordering bool
	}

	//All the child receivers.
	msgReceivers []*MsgReceiver

	// Access to Raven farm and underlying adapters.
	farm *Farm

	//A lock which ensures singleton receiver.
	lock *childlock.Lock
}

func (this *RavenReceiver) defineAccessPort(port string) {
	this.port = port
}

func (this *RavenReceiver) setSource(s Source) *RavenReceiver {
	this.source = s
	return this
}

func (this *RavenReceiver) setId(id string) *RavenReceiver {
	// Make sure we are setting source before Id.
	// Since a Source can have only one receiver at a time, it makes perfect sense to allot
	// source name as ID.
	this.id = this.source.GetName()
	return this
}

func (this *RavenReceiver) GetId() string {
	return this.id
}

func (this *RavenReceiver) GetPort() string {
	return this.port
}

//
// Not mandatory, but if specified receiver will use the specified port for
// communications.
//
func (this *RavenReceiver) SetPort(p string) {
	this.port = p
}

//
// Markall the allotted message receivers as reliable.
//
func (this *RavenReceiver) MarkReliable() *RavenReceiver {
	this.options.isReliable = true

	for _, msgReceiver := range this.msgReceivers {
		msgReceiver.MarkReliable()
	}
	return this
}

//@todo: implement all the necessary validations required for a receiver.
func (this *RavenReceiver) validate() error {

	//Check if Id, Source and farm are defined.
	// check if atleast one receiver is assigned.

	if this.id == "" {
		return fmt.Errorf("An Id needs to be assigned to Receiver. Make sure its unique within source")
	}
	if this.source.GetName() == "" {
		return fmt.Errorf("Receiver Source cannot be Empty")
	}
	if this.farm == nil {
		return fmt.Errorf("You need to define to which farm this receiver belongs.")
	}
	if len(this.msgReceivers) <= 0 {
		return fmt.Errorf("Atleast one msg Receiver needs to be assigned")
	}
	return nil
}

func (this *RavenReceiver) Lock() error {
	if this.lock == nil {
		return nil
	}
	//	fmt.Println("lock")
	r := time.Now().Format(time.RFC3339)
	if err := this.lock.Acquire(r); err != nil {
		return err
	}
	return nil
}

func (this *RavenReceiver) Unlock() error {
	if this.lock == nil {
		return nil
	}
	//fmt.Println("unlock")
	if err := this.lock.Release(); err != nil {
		return err
	}
	//	fmt.Println("unlock done")
	return nil
}

func (this *RavenReceiver) RefreshLock() error {
	if this.lock == nil {
		return nil
	}

	go func() {
		for {
			time.Sleep(CHILD_LOCK_REFRESH_INTERVAL * time.Second)
			func() {
				//fmt.Println("referesh")
				defer util.PanicHandler("Lock Refresh failed")
				if err := this.lock.Refresh(); err != nil {
					fmt.Printf("Lock refresh failed, Error: %s", err.Error())
				}

			}()
		}

	}()

	return nil
}

func (this *RavenReceiver) Stop() {

	defer func() {
		this.Unlock()
		fmt.Println("Lock released")
	}()
	chanx := make(chan bool)
	for _, receiver := range this.msgReceivers {
		fmt.Printf("\nStopping MsgReceiver: %s", receiver.id)

		go func(receiver *MsgReceiver) {
			receiver.stop()
			chanx <- true
		}(receiver)
	}

	for _ = range this.msgReceivers {
		<-chanx
	}

}

func (this *RavenReceiver) Start(f MessageHandler) error {

	if err := this.validate(); err != nil {
		return err
	}

	//Take lock, this ensures only one receiver is receiving from Q.
	if err := this.Lock(); err != nil {
		return err
	}
	defer this.Unlock()

	//Refresh lock
	this.RefreshLock()

	// execute prestart hook of all receivers.
	// once all prestart hooks are successfull start receivers.
	for _, msgreceiver := range this.msgReceivers {
		if err := msgreceiver.preStart(); err != nil {
			return err
		}
	}

	//@todo: Start receivers.
	// Since the start functions of receivers block, we need to start
	// receivers as seperate goroutines.
	// @todo: need to control these receivers from channels.
	for _, msgreceiver := range this.msgReceivers {
		go msgreceiver.StartHeartBeat()
		go msgreceiver.start(f)
	}

	//Once all the receivers are up boot up the server.
	if err := StartServer(this); err != nil {
		return err
	}
	return nil
}

//
// Get all the ravens still wandering around.
//
func (this *RavenReceiver) GetInFlightRavens() map[string]string {
	holder := make(map[string]string, len(this.msgReceivers))
	for _, r := range this.msgReceivers {
		var val string
		cc, err := r.GetInFlightRavens()
		if err != nil {
			val = err.Error()
		} else {
			val = strconv.Itoa(cc)
		}
		holder[r.id] = val
	}
	return holder
}

func (this *RavenReceiver) GetDeadBoxCount() map[string]string {
	holder := make(map[string]string, 0)
	for _, r := range this.msgReceivers {
		var val string
		cc, err := r.GetDeadBoxCount()
		if err != nil {
			val = err.Error()
		} else {
			val = strconv.Itoa(cc)
		}
		holder[r.id] = val
	}
	return holder
}

func (this *RavenReceiver) FlushDeadBox() map[string]string {
	holder := make(map[string]string, 0)
	for _, r := range this.msgReceivers {
		var val string = "OK"
		if err := r.flushDeadBox(); err != nil {
			val = err.Error()
		}
		holder[r.id] = val
	}
	return holder
}

func (this *RavenReceiver) FlushAll() map[string]string {
	holder := make(map[string]string, 0)
	for _, r := range this.msgReceivers {
		var val string = "OK"
		if err := r.flushAll(); err != nil {
			val = err.Error()
		}
		holder[r.id] = val
	}
	return holder
}

func (this *RavenReceiver) ShowDeadBox() ([]*Message, error) {
	m := make([]*Message, 0)
	for _, r := range this.msgReceivers {

		msgs, err := r.showDeadBox()
		if err != nil {
			return nil, err
		}
		m = append(m, msgs...)
	}
	return m, nil
}

func (this *RavenReceiver) ShowMessage() {
	fmt.Println("\n\n--------------------------------------------")
	fmt.Printf("MessageReceivers Started:\n")
	w := tabwriter.NewWriter(os.Stdout, 10, 10, 1, ' ', tabwriter.AlignRight)

	fmt.Fprintln(w)
	fmt.Fprintf(w, "ReceiverId \t IsReliable \t ProcBox \t DeadBox \t ")
	fmt.Fprintln(w)
	for _, r := range this.msgReceivers {
		fmt.Fprintf(w, "%s \t %t \t %s \t %s \t", r.id, r.options.isReliable, r.procBox.GetName(), r.deadBox.GetName())
		fmt.Fprintln(w)
	}
	w.Flush()
	fmt.Println("--------------------------------------------")
	fmt.Printf("\nConsumer Communication Port: %s\n", this.GetPort())
	fmt.Println("--------------------------------------------")
	return
}
