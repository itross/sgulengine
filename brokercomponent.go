package sgulengine

import (
	"fmt"

	"github.com/itross/sgul"
)

type (
	// Event Broker configuration structs

	// OutboundEvent is an evt_to_producer configuration entry.
	OutboundEvent struct {
		Name      string
		Publisher string
	}

	// InboundEvent is a evt_to_consumer configuration entry.
	InboundEvent struct {
		Name       string
		Subscriber string
	}

	// Events is the event mapping configuration struct.
	Events struct {
		Outbound []OutboundEvent
		Inbound  []InboundEvent
	}

	// BrokerConfig .
	BrokerConfig struct {
		Events Events
	}

	outboudBroker struct {
		publishers map[string]*sgul.AMQPPublisher
	}

	inboundBroker struct {
		subscribers map[string]*sgul.AMQPSubscriber
	}

	// BrokerComponent is the AMQP communication manager.
	BrokerComponent struct {
		BaseComponent
		config     BrokerConfig
		connection *sgul.AMQPConnection
		outB       outboudBroker
		inB        inboundBroker
	}
)

// NewBroker returns new Broker component instance.
func NewBroker() *BrokerComponent {
	return &BrokerComponent{
		BaseComponent: NewBaseComponent("broker"),
		outB:          outboudBroker{publishers: make(map[string]*sgul.AMQPPublisher)},
		inB:           inboundBroker{subscribers: make(map[string]*sgul.AMQPSubscriber)},
	}
}

// Configure .
func (brk *BrokerComponent) Configure(conf interface{}) error {
	sgul.LoadConfiguration(brk.config)
	return nil
}

// Start will start the Broker component starting a connection to the AMQP server.
func (brk *BrokerComponent) Start(e *Engine) error {
	brk.connection = sgul.NewAMQPConnection()
	brk.logger.Debugf("Connecting to AMQP server at: %s", brk.connection.URI)

	if err := brk.connection.Connect(); err != nil {
		e.cErrs <- fmt.Errorf("error starting Event Broker component: %s", err)
		return nil
	}
	brk.logger.Info("AMQP connection esabilished")

	// initialize events -> publishers map
	for _, evt := range brk.config.Events.Outbound {
		brk.outB.publishers[evt.Name] = brk.connection.Publishers[evt.Publisher]
		brk.logger.Infof("%s publisher registered for %s event", evt.Publisher, evt.Name)
	}

	return nil
}

// Shutdown will stop AMQP channel and connection.
func (brk *BrokerComponent) Shutdown() {
	if brk.connection != nil {
		if err := brk.connection.Close(); err != nil {
			brk.logger.Errorf("error shutting down Broker Component", "error", err)
		}
	}
}

// PublisherFor returns the AMQP Publisher registered for the event name.
func (brk *BrokerComponent) PublisherFor(evtName string) *sgul.AMQPPublisher {
	return brk.outB.publishers[evtName]
}

// func (brk *BrokerComponent) AddEventPublisher(eventName, routingKey string) error {
// 	if err := brk.connection.Publisher()
// }
