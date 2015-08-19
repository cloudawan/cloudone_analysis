// Copyright 2015 CloudAwan LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event

/*
import (
	"fmt"
	"testing"
	"time"
)


func TestConsistency(t *testing.T) {
	id := "_api_v1_namespaces_default_events_nginx.13f8ec6987d1e416"
	acknowledge := false
	fmt.Println(Acknowledge("default", id, acknowledge))
	jsonMap, _ := GetEvent(indexKubernetesEventIndex, "default", id)
	fmt.Println(jsonMap["searchMetaData"].(map[string]interface{})["acknowledge"])

	found := false
	jsonSlice, _ := SearchHistoricalEvent("*", nil, nil, acknowledge)
	for _, json := range jsonSlice {
		if json.(map[string]interface{})["_id"] == id {
			//fmt.Println("=1=", json.(map[string]interface{})["_source"].(map[string]interface{})["searchMetaData"])
			found = true
		}
	}
	fmt.Println("acknowledge: ", acknowledge, ", 1th found: ", found)

	found = false
	jsonSlice, _ = SearchHistoricalEvent("*", nil, nil, !acknowledge)
	for _, json := range jsonSlice {
		if json.(map[string]interface{})["_id"] == id {
			//fmt.Println("=1=", json.(map[string]interface{})["_source"].(map[string]interface{})["searchMetaData"])
			found = true
		}
	}
	fmt.Println("acknowledge: ", !acknowledge, ", 1th found: ", found)

	time.Sleep(1000 * time.Millisecond)

	found = false
	jsonSlice, _ = SearchHistoricalEvent("*", nil, nil, acknowledge)
	for _, json := range jsonSlice {
		if json.(map[string]interface{})["_id"] == id {
			//fmt.Println("=2=", json.(map[string]interface{})["_source"].(map[string]interface{})["searchMetaData"])
			found = true
		}
	}
	fmt.Println("acknowledge: ", acknowledge, ", 2nd found: ", found)

	found = false
	jsonSlice, _ = SearchHistoricalEvent("*", nil, nil, !acknowledge)
	for _, json := range jsonSlice {
		if json.(map[string]interface{})["_id"] == id {
			//fmt.Println("=2=", json.(map[string]interface{})["_source"].(map[string]interface{})["searchMetaData"])
			found = true
		}
	}
	fmt.Println("acknowledge: ", !acknowledge, ", 2nd found: ", found)
}

func TestAnowledge(t *testing.T) {
	fmt.Println(Acknowledge("default", "_api_v1_namespaces_default_events_nginx.13f8ec6987d1e416", true))
}


func TestGetAllEvent(t *testing.T) {
	fmt.Println(RecordHistoricalEvent("172.16.0.113", 8080))
}

func TestGetAllEvent(t *testing.T) {
	from := time.Now().AddDate(-1, 0, 0)
	to := time.Now()
	jsonSlice, err := SearchHistoricalEvent("*", &from, &to, false, 1024*1024, 0)
	fmt.Println(jsonSlice, err)
}
*/
