/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package platform_connector_lib

import (
	"github.com/SENERGY-Platform/iot-broker-client"
	"log"

	"github.com/Shopify/sarama"
)

type CommandHandler func(endpoint string, protocolParts map[string]string) (responseParts map[string]string, err error)

type Connector struct {
	Config         Config
	CommandHandler CommandHandler //must be able to handle concurrent calls
	kafkaproducer  sarama.AsyncProducer
	producer	   *iot_broker_client.Publisher
	consumer       *iot_broker_client.Consumer
	openid         *OpenidToken
}

func Init(config Config, commandHandler CommandHandler) (connector *Connector, err error) {
	connector = &Connector{Config: config, CommandHandler: commandHandler, kafkaproducer: initKafkaProducer(config.ZookeeperUrl)}
	connector.producer, err = InitProducer()
	if err != nil {
		return
	}
	connector.consumer, err = connector.InitConsumer()
	return
}

func (this *Connector) Stop() {
	this.kafkaproducer.Close()
	this.producer.Close()
	this.consumer.Close()
}

func (this *Connector) HandleEvent(username string, password string, endpoint string, protocolParts map[string]string) (total int, success int, ignore int, fail int, err error) {
	token, err := this.GetOpenidPasswordToken(username, password)
	if err != nil {
		log.Println("ERROR HandleEvent::GetOpenidPasswordToken()", err)
		return total, success, ignore, fail, err
	}
	return this.HandleEventWithAuthToken(token.JwtToken(), endpoint, protocolParts)
}

//is able to handle concurrent calls
func (this *Connector) HandleEventWithAuthToken(token JwtToken, endpoint string, protocolParts map[string]string) (total int, success int, ignore int, fail int, err error) {
	protocol := []ProtocolPart{}
	for key, val := range protocolParts {
		protocol = append(protocol, ProtocolPart{Name: key, Value: val})
	}
	return this.handleEvent(token, endpoint, protocol)
}
