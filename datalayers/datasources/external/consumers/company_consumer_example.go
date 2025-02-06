// ==== we will get data from an external source via message queue or http/gRPC from other services === */
/*

FOR EXAMPLE, WE NEED INFORMATION OF COMPANY ENTITES FROM COMPANY SERVICE VIA MESSAGE QUEUE
BELOW IS AN EXAMPLE OF HOW TO GET COMPANY INFORMATION FROM COMPANY SERVICE VIA MESSAGE QUEUE
(not working yet)
*/

package company_kafka

import (

	"proposal-template/pkg/kafka"
)

type CompanyMessage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

// EventName implements the ConsumerMessage interface
func (c *CompanyMessage) EventName() string {
	return "CompanyMessage"
}

// // CompanyConsumer handles consuming messages from the "company-updates" topic
type CompanyConsumer struct {
	subscriber kafka.Subscriber
}

func NewCompanyConsumer(subscriber kafka.Subscriber)*CompanyConsumer {
	return &CompanyConsumer{}
}
func (c *CompanyConsumer) ConsumeMessages() {
	c.subscriber.ConsumeMessages()
}

