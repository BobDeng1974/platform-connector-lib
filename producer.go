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
)


func InitProducer(amqpUrl string) (producer *iot_broker_client_lib.Publisher,err error){
	return iot_broker_client_lib.NewPublisher(amqpUrl)
}

func (this *Connector) produce(topic string, message []byte, keys ...string)(err error) {
	return this.producer.Publish(topic, message, keys...)
}

func (this *Connector) CloseProducer(){
	this.producer.Close()
}