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

package cluster

import (
	"fmt"
	//"strconv"
	"testing"
	"time"
)

func TestSaveClusterSingletonLockIndex(t *testing.T) {
	jsonMap := make(map[string]interface{})
	jsonMap["id"] = "ip1"
	jsonMap["firstTimeStamp"] = time.Now().Format(time.RFC3339Nano)
	jsonMap["lastTimeStamp"] = time.Now().Format(time.RFC3339Nano)

	fmt.Println(saveClusterSingletonLock(
		indexClusterSingletonLock,
		typeCloudoneAnalysis,
		"test",
		jsonMap,
		false))
}

func TestLoadClusterSingletonLockIndex(t *testing.T) {
	jsonMap, err := loadClusterSingletonLock(indexClusterSingletonLock,
		typeCloudoneAnalysis,
		"test")
	fmt.Println(jsonMap, err)
	if err != nil && err.Error() == "record not found" {
		fmt.Println("record not found")
	}
	if jsonMap != nil {
		fmt.Println(time.Parse(time.RFC3339Nano, jsonMap["firstTimeStamp"].(string)))
	}

}
