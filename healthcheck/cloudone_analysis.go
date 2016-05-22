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

package healthcheck

import (
	"errors"
	"github.com/cloudawan/cloudone_analysis/utility/configuration"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strconv"
	"time"
)

func CreateCloudoneAnalysisControl() (*CloudoneAnalysisControl, error) {
	restapiPort, ok := configuration.LocalConfiguration.GetInt("restapiPort")
	if ok == false {
		log.Error("Can't find restapiPort")
		return nil, errors.New("Can't find restapiPort")
	}
	cloudoneAnalysisControl := &CloudoneAnalysisControl{
		restapiPort,
	}
	return cloudoneAnalysisControl, nil
}

const (
	restApiTimeout = time.Millisecond * 300
)

type CloudoneAnalysisControl struct {
	RestapiPort int
}

func (cloudoneAnalysisControl *CloudoneAnalysisControl) testRestAPI() bool {
	result, _ := restclient.HealthCheck(
		"https://127.0.0.1:"+strconv.Itoa(cloudoneAnalysisControl.RestapiPort)+"/apidocs.json",
		nil,
		restApiTimeout)
	return result
}

func (cloudoneAnalysisControl *CloudoneAnalysisControl) testStorageElasticSearch() bool {
	jsonMap := make(map[string]interface{})
	jsonMap["updatedTime"] = time.Now().Format(time.RFC3339Nano)
	if err := saveTest("test", "test", "test", jsonMap); err != nil {
		log.Error(err)
		return false
	} else {
		return true
	}
}

func (cloudoneAnalysisControl *CloudoneAnalysisControl) GetStatus() map[string]interface{} {
	jsonMap := make(map[string]interface{})
	jsonMap["restapi"] = cloudoneAnalysisControl.testRestAPI()
	jsonMap["elasticsearch"] = cloudoneAnalysisControl.testStorageElasticSearch()
	return jsonMap
}
