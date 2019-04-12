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

package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"log"
	"time"
)

type Producer struct {
	broker []string
	logger *log.Logger
}

func PrepareProducer(zk string) (*Producer, error) {
	broker, err := GetBroker(zk)
	if err == nil && len(broker) == 0 {
		err = errors.New("missing kafka broker")
	}
	return &Producer{broker: broker}, err
}

func (this *Producer) Log(logger *log.Logger) {
	this.logger = logger
}

func (this *Producer) getProducer(topic string) (writer *kafka.Writer, err error) {
	var logger *log.Logger
	if this.logger != nil {
		logger = this.logger
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:     this.broker,
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		MaxAttempts: 25,
		Logger:      logger,
	})
	return writer, err
}

func (this *Producer) Produce(topic string, message string) (err error) {
	if this.logger != nil {
		this.logger.Println("DEBUG: produce ", topic, message)
	}
	writer, err := this.getProducer(topic)
	if err != nil {
		if this.logger != nil {
			this.logger.Println("ERROR: while getting kafka producer:", err)
		}
		return err
	}
	defer writer.Close()
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(message),
			Time:  time.Now(),
		},
	)
}
