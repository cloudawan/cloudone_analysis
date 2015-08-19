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

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cloudawan/kubernetes_management_analysis/control"
	"github.com/cloudawan/kubernetes_management_utility/logger"
	"strconv"
	"strings"
	"time"
)

func RecordHistoricalEvent(kubeapiHost string, kubeapiPort int) (returnedError error) {
	jsonMapSlice, err := control.GetAllEvent(kubeapiHost, kubeapiPort)
	if err != nil {
		log.Error(err)
		return err
	}
	hasError := false
	erroerMessageBuffer := bytes.Buffer{}
	for _, jsonMap := range jsonMapSlice {
		namespace, _ := jsonMap["metadata"].(map[string]interface{})["namespace"].(string)
		selfLink, _ := jsonMap["metadata"].(map[string]interface{})["selfLink"].(string)

		jsonMap["searchMetaData"] = make(map[string]interface{})
		jsonMap["searchMetaData"].(map[string]interface{})["acknowledge"] = false

		if err := saveKubernetesEvent(indexKubernetesEventIndex, namespace, getEventID(selfLink), jsonMap, false); err != nil {
			log.Error(err)
			erroerMessageBuffer.WriteString(err.Error())
			hasError = true
		} else {
			// Remove after saving in Elastic Search
			if err := control.DeleteEvent(kubeapiHost, kubeapiPort, selfLink); err != nil {
				log.Error(err)
				erroerMessageBuffer.WriteString(err.Error())
				hasError = true
			}
		}
	}
	if hasError {
		return errors.New(erroerMessageBuffer.String())
	} else {
		return nil
	}
}

func SearchHistoricalEvent(namespace string, from *time.Time,
	to *time.Time, acknowledge bool, size int, offset int) (returnedJsonSlice []interface{}, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("SearchHistoricalEvent Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedJsonSlice = nil
			returnedError = err.(error)
		}
	}()

	if from != nil && to != nil && from.After(*to) {
		return nil, errors.New("From " + from.String() + " can't be after to " + to.String())
	}

	var acknowledgeText string
	if acknowledge {
		acknowledgeText = "true"
	} else {
		acknowledgeText = "false"
	}

	var queryField string
	if from == nil && to != nil {
		lte := to.UTC().Format(time.RFC3339Nano)
		queryField = `"query": {		
			"range" : {
				"lastTimestamp" : {
					"lte": "` + lte + `",
					"time_zone": "+0:00"
				}
			}
	    },`
	} else if from != nil && to == nil {
		gte := from.UTC().Format(time.RFC3339Nano)
		queryField = `"query": {		
			"range" : {
				"lastTimestamp" : {
					"gte": "` + gte + `",
					"time_zone": "+0:00"
				}
			}
	    },`
	} else if from != nil && to != nil {
		lte := to.UTC().Format(time.RFC3339Nano)
		gte := from.UTC().Format(time.RFC3339Nano)
		queryField = `"query": {		
			"range" : {
				"lastTimestamp" : {
					"lte": "` + lte + `",
					"gte": "` + gte + `",
					"time_zone": "+0:00"
				}
			}
	    },`
	} else {
		queryField = ``
	}

	query := `
	{
		"query": {
			"filtered": {
				` + queryField + `
				"filter": {
					"term": { 
						"searchMetaData.acknowledge": ` + acknowledgeText + `
					}
				}
			}
		},
		"sort" : [
	 		{ 
				"lastTimestamp" : "desc"
			}
    	],
		"size": ` + strconv.Itoa(size) + `,
		"from": ` + strconv.Itoa(offset) + `
	}
	`

	byteSlice, err := searchKubernetesEventRawJson(indexKubernetesEventIndex, "*", query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	jsonMap := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(byteSlice))
	decoder.UseNumber()
	if err := decoder.Decode(&jsonMap); err != nil {
		log.Error(err)
		return nil, err
	}

	jsonSlice, ok := jsonMap["hits"].(map[string]interface{})["hits"].([]interface{})
	if ok {
		return jsonSlice, nil
	} else {
		return nil, errors.New("Fail to get with byteSlice " + string(byteSlice))
	}
}

func getEventID(selfLink string) string {
	return strings.Replace(selfLink, "/", "_", -1)
}

func Acknowledge(namespace string, id string, acknowledge bool) error {
	jsonMap, err := GetEvent(indexKubernetesEventIndex, namespace, id)
	if err != nil {
		log.Error(err)
		return err
	} else {
		jsonMap["searchMetaData"].(map[string]interface{})["acknowledge"] = acknowledge
		if err := saveKubernetesEvent(indexKubernetesEventIndex, namespace, id, jsonMap, true); err != nil {
			log.Error(err)
			return err
		} else {
			return nil
		}
	}
}
