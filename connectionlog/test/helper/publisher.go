/*
 * Copyright 2019 InfAI (CC SES)
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

package helper

import (
	"errors"
	"github.com/SENERGY-Platform/platform-connector-lib/connectionlog/test/config"
	"github.com/segmentio/kafka-go"
	"github.com/wvanbergen/kazoo-go"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
)

type Publisher struct {
	config          config.Config
	devicetypes     *kafka.Writer
	protocols       *kafka.Writer
	devices         *kafka.Writer
	hubs            *kafka.Writer
	concepts        *kafka.Writer
	characteristics *kafka.Writer
}

func GetProducer(broker []string, topic string, debug bool) (writer *kafka.Writer, err error) {
	var logger *log.Logger
	if debug {
		logger = log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0)
	} else {
		logger = log.New(ioutil.Discard, "", 0)
	}
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:     broker,
		Topic:       topic,
		MaxAttempts: 10,
		Logger:      logger,
	})
	return writer, err
}

func GetBroker(zk string) (brokers []string, err error) {
	return getBroker(zk)
}

func getBroker(zkUrl string) (brokers []string, err error) {
	zookeeper := kazoo.NewConfig()
	zookeeper.Logger = log.New(ioutil.Discard, "", 0)
	zk, chroot := kazoo.ParseConnectionString(zkUrl)
	zookeeper.Chroot = chroot
	if kz, err := kazoo.NewKazoo(zk, zookeeper); err != nil {
		return brokers, err
	} else {
		return kz.BrokerList()
	}
}

func GetKafkaController(zkUrl string) (controller string, err error) {
	zookeeper := kazoo.NewConfig()
	zookeeper.Logger = log.New(ioutil.Discard, "", 0)
	zk, chroot := kazoo.ParseConnectionString(zkUrl)
	zookeeper.Chroot = chroot
	kz, err := kazoo.NewKazoo(zk, zookeeper)
	if err != nil {
		return controller, err
	}
	controllerId, err := kz.Controller()
	if err != nil {
		return controller, err
	}
	brokers, err := kz.Brokers()
	if err != nil {
		return controller, err
	}
	return brokers[controllerId], err
}

func InitTopic(zkUrl string, topics ...string) (err error) {
	return InitTopicWithConfig(zkUrl, 1, 1, topics...)
}

func InitTopicWithConfig(zkUrl string, numPartitions int, replicationFactor int, topics ...string) (err error) {
	controller, err := GetKafkaController(zkUrl)
	if err != nil {
		log.Println("ERROR: unable to find controller", err)
		return err
	}
	if controller == "" {
		log.Println("ERROR: unable to find controller")
		return errors.New("unable to find controller")
	}
	initConn, err := kafka.Dial("tcp", controller)
	if err != nil {
		log.Println("ERROR: while init topic connection ", err)
		return err
	}
	defer initConn.Close()
	for _, topic := range topics {
		err = initConn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
			ConfigEntries: []kafka.ConfigEntry{
				{ConfigName: "retention.ms", ConfigValue: "-1"},
				{ConfigName: "retention.bytes", ConfigValue: "-1"},
				{ConfigName: "cleanup.policy", ConfigValue: "compact"},
				{ConfigName: "delete.retention.ms", ConfigValue: "100"},
				{ConfigName: "segment.ms", ConfigValue: "100"},
				{ConfigName: "min.cleanable.dirty.ratio", ConfigValue: "0.01"},
			},
		})
		if err != nil {
			log.Println("ERROR:", err)
			debug.PrintStack()
			return err
		}
	}
	return nil
}
