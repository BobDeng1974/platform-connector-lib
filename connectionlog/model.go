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

package connectionlog

import "time"

type ConnectorLog struct {
	Connected bool         `json:"connected"`
	Connector string       `json:"connector"`
	Gateways  []GatewayLog `json:"gateways"`
	Devices   []DeviceLog  `json:"devices"`
	Time      time.Time    `json:"time"`
}

type GatewayLog struct {
	Connected bool      `json:"connected"`
	Connector string    `json:"connector"`
	Gateway   string    `json:"gateway"`
	Time      time.Time `json:"time"`
}

type DeviceLog struct {
	Connected bool      `json:"connected"`
	Connector string    `json:"connector"`
	Device    string    `json:"device"`
	Time      time.Time `json:"time"`
}