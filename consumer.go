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
	"github.com/SENERGY-Platform/iot-broker-client-lib"
	"log"

	"encoding/json"
)

func (this *Connector) InitConsumer() (consumer *iot_broker_client_lib.Consumer, err error) {
	consumer, err = iot_broker_client_lib.NewConsumer(this.Config.AmqpUrl, "queue_"+this.Config.Protocol, this.Config.Protocol, false, int(this.Config.AmqpPrefetchCount), func(msg []byte) error {
		return this.handleMessage(string(msg))
	})
	if err != nil {
		log.Println("ERROR: unable to create amqp consumer", err)
		return
	}
	err = consumer.BindAll()
	if err != nil {
		log.Println("ERROR: unable to bind consumer to all devices", err)
		return
	}
	return
}

func (this *Connector) handleMessage(msg string) (err error) {
	log.Println("consume kafka msg: ", msg)
	envelope := Envelope{}
	err = json.Unmarshal([]byte(msg), &envelope)
	if err != nil {
		log.Println("ERROR: ", err)
		return nil //ignore marshaling errors --> no repeat; errors would definitely reoccur
	}
	payload, err := json.Marshal(envelope.Value)
	if err != nil {
		log.Println("ERROR: ", err)
		return nil //ignore marshaling errors --> no repeat; errors would definitely reoccur
	}
	return this.handleCommand(string(payload))
}
