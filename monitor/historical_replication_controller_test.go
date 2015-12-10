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

package monitor

import (
	"fmt"
	"testing"
	"time"
)

func TestSearchHistoricalReplicationControllerMetrics(t *testing.T) {
	current := time.Now()
	nodeAmount := 10
	from := current.Add(-1 * time.Minute)
	to := current.Add(-0 * time.Minute)
	byteSlice, err := searchHistoricalReplicationControllerMetrics("default", "cloudone-all", nodeAmount, from, to)
	fmt.Println(string(byteSlice), err)
}

/*
func TestGetAllHistoricalReplicationControllerMetrics(t *testing.T) {
	current := time.Now()
	nodeAmount := 10
	from := current.Add(-11 * time.Minute)
	to := current.Add(-1 * time.Minute)
	fmt.Println(GetAllHistoricalReplicationControllerMetrics("default", nodeAmount, from, to))
}


func TestGetHistoricalReplicationControllerMetrics(t *testing.T) {
	current := time.Now()
	nodeAmount := 10
	from := current.Add(-11 * time.Minute)
	to := current.Add(-1 * time.Minute)
	fmt.Println(GetHistoricalReplicationControllerMetrics("default", "private-repository", nodeAmount, from, to))
}


func TestRecordHistoricalReplicationController(t *testing.T) {
	fmt.Println(RecordHistoricalReplicationController("172.16.0.113", 8080, "default", "cassandra"))
}
*/
