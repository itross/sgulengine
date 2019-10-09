package sgulengine

import (
	"github.com/itross/sgul"
)

type outboudBroker struct {
	publishers map[string]*sgul.AMQPPublisher
}

type inboundBroker struct {
	subscribers map[string]*sgul.AMQPSubscriber
}

// BrokerComponent is the AMQP communication manager.
type BrokerComponent struct {
	BaseComponent
	config     sgul.AMQP
	connection *sgul.AMQPConnection
	outB       outboudBroker
	inB        inboundBroker
}

// NewBroker returns new Broker component instance.
func NewBroker() *BrokerComponent {
	return &BrokerComponent{
		BaseComponent: NewBaseComponent("broker"),
	}
}

// Configure .
func (brk *BrokerComponent) Configure(conf interface{}) error {
	brk.config = conf.(sgul.AMQP)
	return nil
}

// Start will start the Broker component starting a connection to the AMQP server.
func (brk *BrokerComponent) Start(e *Engine) error {
	brk.connection = sgul.NewAMQPConnection()
	brk.logger.Debugf("Connecting to AMQP server at: %s", brk.connection.URI)

	if err := brk.connection.Connect(); err != nil {
		return err
	}
	brk.logger.Debug("AMQP connection esabilished")
	return nil
}

// Shutdown will stop AMQP channel and connection.
func (brk *BrokerComponent) Shutdown() {
	if err := brk.connection.Close(); err != nil {
		brk.logger.Errorf("error shutting down Broker Component", "error", err)
	}
}
