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

package model

import "strings"

type TaskInfo struct {
	WorkerId           string `json:"worker_id"`
	TaskId             string `json:"task_id"`
	CompletionStrategy string `json:"completion_strategy"`
	Time               string `json:"time"`
}

type ProtocolRequest struct {
	Input map[string]string `json:"input"`
}

type ProtocolResponse struct {
	Output map[string]string `json:"output"`
}

type Metadata struct {
	Device               Device   `json:"device"`
	Service              Service  `json:"service"`
	Protocol             Protocol `json:"protocol"`
	InputCharacteristic  string         `json:"input_characteristic,omitempty"`
	OutputCharacteristic string         `json:"output_characteristic,omitempty"`
}

type ProtocolMsg struct {
	Request  ProtocolRequest  `json:"request"`
	Response ProtocolResponse `json:"response"`
	TaskInfo TaskInfo         `json:"task_info"`
	Metadata Metadata         `json:"metadata"`
}

const Optimistic = "optimistic"

type Envelope struct {
	DeviceId  string      `json:"device_id,omitempty"`
	ServiceId string      `json:"service_id,omitempty"`
	Value     interface{} `json:"value"`
}

func ServiceIdToTopic(id string) string {
	id = strings.ReplaceAll(id, "#", "_")
	id = strings.ReplaceAll(id, ":", "_")
	return id
}
