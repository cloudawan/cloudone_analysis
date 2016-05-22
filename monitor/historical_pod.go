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
	"bytes"
	"errors"
	"github.com/cloudawan/cloudone_utility/logger"
	"github.com/cloudawan/cloudone_utility/restclient"
	"strings"
	"time"
)

func RecordHistoricalPod(kubeApiServerEndPoint string, kubeApiServerToken string, namespace string, replicationControllerName string, podName string) (returnedPodContainerRecordSlice []map[string]interface{}, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("RecordHistoricalPod Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedPodContainerRecordSlice = nil
			returnedError = err.(error)
		}
	}()

	podContainerRecordSlice := make([]map[string]interface{}, 0)

	headerMap := make(map[string]string)
	headerMap["Authorization"] = kubeApiServerToken

	result, err := restclient.RequestGet(kubeApiServerEndPoint+"/api/v1/namespaces/"+namespace+"/pods/"+podName+"/", headerMap, true)
	if err != nil {
		log.Error("Fail to get pod inofrmation with endpoint %s, token: %s, namespace: %s, pod name: %s, error %s", kubeApiServerEndPoint, kubeApiServerToken, namespace, podName, err.Error())
		return nil, err
	}
	jsonMap, _ := result.(map[string]interface{})

	errorBuffer := bytes.Buffer{}
	errorBuffer.WriteString("The following container has error: ")
	errorHappened := false

	kubeletHost, _ := jsonMap["status"].(map[string]interface{})["hostIP"].(string)
	uid, _ := jsonMap["metadata"].(map[string]interface{})["uid"].(string)
	containerSlice, _ := jsonMap["spec"].(map[string]interface{})["containers"].([]interface{})
	for _, container := range containerSlice {
		containerName, _ := container.(map[string]interface{})["name"].(string)
		url := "https://" + kubeletHost + ":10250/stats/" + namespace + "/" + podName + "/" + uid + "/" + containerName
		result, err := restclient.RequestGet(url, nil, true)
		containerJsonMap, _ := result.(map[string]interface{})
		if err != nil {
			errorHappened = true
			log.Error("Request to url %s error %s", url, err)
			errorBuffer.WriteString("Request to url " + url + " error " + err.Error())
		} else {
			// ElasticSearch doesn't allow to use character '.' in the field name so it should be replaced with '_'
			if containerJsonMap["spec"].(map[string]interface{})["labels"] != nil {
				for key, value := range containerJsonMap["spec"].(map[string]interface{})["labels"].(map[string]interface{}) {
					if strings.Contains(key, ".") {
						newKey := strings.Replace(key, ".", "_", -1)
						containerJsonMap["spec"].(map[string]interface{})["labels"].(map[string]interface{})[newKey] = value
						delete(containerJsonMap["spec"].(map[string]interface{})["labels"].(map[string]interface{}), key)
					}
				}
			}

			// Historical data
			containerRecordSlice, err := splitHistoricalDataIntoSecondBased(namespace, replicationControllerName,
				podName, containerName, containerJsonMap)
			if err != nil {
				errorHappened = true
				log.Error("Save container record %s error %s", containerJsonMap, err)
				errorBuffer.WriteString("Save container record  " + containerName + " error " + err.Error())
			} else {
				podContainerRecordSlice = append(podContainerRecordSlice, containerRecordSlice...)
			}
		}
	}

	if errorHappened {
		log.Error("Fail to get all container inofrmation with endpoint %s, token: %s, namespace: %s, pod name: %s, error %s", kubeApiServerEndPoint, kubeApiServerToken, namespace, podName, errorBuffer.String())
		return nil, errors.New(errorBuffer.String())
	} else {
		return podContainerRecordSlice, nil
	}
}

func splitHistoricalDataIntoSecondBased(namespace string, replicationControllerName string,
	podName string, containerName string, jsonMap map[string]interface{}) ([]map[string]interface{}, error) {
	statsSlice, ok := jsonMap["stats"].([]interface{})
	if ok {
		containerRecordSlice := make([]map[string]interface{}, 0)
		for _, stats := range statsSlice {
			statsJsonMap, ok := stats.(map[string]interface{})
			if ok {
				timestampField, ok := statsJsonMap["timestamp"].(string)
				if ok {
					timestamp, err := time.Parse(time.RFC3339Nano, timestampField)
					if err != nil {
						log.Error("Parse timestamp error %s", stats)
						return nil, errors.New("Parse timestamp error")
					} else {
						containerRecord := make(map[string]interface{})
						for key, value := range jsonMap {
							containerRecord[key] = value
						}

						index := getDocumentIndex(namespace)
						documentType := getDocumentType(replicationControllerName)
						id := getDocumentID(podName, containerName, timestamp)
						containerRecord["searchMetaData"] = make(map[string]interface{})
						containerRecord["searchMetaData"].(map[string]interface{})["namespace"] = namespace
						containerRecord["searchMetaData"].(map[string]interface{})["replicationControllerName"] = replicationControllerName
						containerRecord["searchMetaData"].(map[string]interface{})["podName"] = podName
						containerRecord["searchMetaData"].(map[string]interface{})["containerName"] = containerName
						containerRecord["searchMetaData"].(map[string]interface{})["index"] = index
						containerRecord["searchMetaData"].(map[string]interface{})["documentType"] = documentType
						containerRecord["searchMetaData"].(map[string]interface{})["id"] = id
						containerRecord["stats"] = stats

						containerRecordSlice = append(containerRecordSlice, containerRecord)
						/*
							//err := saveContainerRecord(index, documentType, id, jsonMap)

							if err != nil {
								log.Error("Save error %s", err)
								return nil, err
							}
						*/
					}

				} else {
					log.Error("Container doesn't have timestamp field in stats %s", stats)
					return nil, errors.New("Container doesn't have timestamp field in stats")
				}
			} else {
				log.Error("Container doesn't have stats field as map %s", stats)
				return nil, errors.New("Container doesn't have stats field as map")
			}
		}

		return containerRecordSlice, nil
	} else {
		log.Error("Container doesn't have stats field")
		return nil, errors.New("Container doesn't have stats field")
	}
}

func getDocumentIndex(namespace string) string {
	return indexContainerMetricsIndexPrefix + strings.ToLower(namespace)
}

func getDocumentType(replicationControllerName string) string {
	return indexContainerMetricsTypePrefix + strings.ToLower(replicationControllerName)
}

func getDocumentID(podName string, containerName string, timestamp time.Time) string {
	return podName + "_" + containerName + "_" + timestamp.UTC().Format("2006-01-02T15-04-05")
}

func getReplicationControllerNameFromDocumentType(documentType string) string {
	return documentType[len(indexContainerMetricsTypePrefix):len(documentType)]
}
