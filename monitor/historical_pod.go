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
	"github.com/cloudawan/kubernetes_management_utility/logger"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strconv"
	"strings"
	"time"
)

func RecordHistoricalPod(kubeapiHost string, kubeapiPort int, namespace string, replicationControllerName string, podName string) (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("RecordHistoricalPod Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	result, err := restclient.RequestGet("http://"+kubeapiHost+":"+strconv.Itoa(kubeapiPort)+"/api/v1/namespaces/"+namespace+"/pods/"+podName+"/", true)
	if err != nil {
		log.Error("Fail to get pod inofrmation with host %s, port: %d, namespace: %s, pod name: %s, error %s", kubeapiHost, kubeapiPort, namespace, podName, err.Error())
		return err
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
		result, err := restclient.RequestGet(url, true)
		containerJsonMap, _ := result.(map[string]interface{})
		if err != nil {
			errorHappened = true
			log.Error("Request to url %s error %s", url, err)
			errorBuffer.WriteString("Request to url " + url + " error " + err.Error())
		} else {
			// Historical data
			err := splitHistoricalDataIntoSecondBased(namespace, replicationControllerName,
				podName, containerName, containerJsonMap)
			if err != nil {
				errorHappened = true
				log.Error("Save container record %s error %s", containerJsonMap, err)
				errorBuffer.WriteString("Save container record  " + containerName + " error " + err.Error())
			}
		}
	}

	if errorHappened {
		log.Error("Fail to get all container inofrmation with host %s, port: %d, namespace: %s, pod name: %s, error %s", kubeapiHost, kubeapiPort, namespace, podName, errorBuffer.String())
		return errors.New(errorBuffer.String())
	} else {
		return nil
	}
}

func splitHistoricalDataIntoSecondBased(namespace string, replicationControllerName string,
	podName string, containerName string, jsonMap map[string]interface{}) error {
	statsSlice, ok := jsonMap["stats"].([]interface{})
	if ok {
		for _, stats := range statsSlice {
			statsJsonMap, ok := stats.(map[string]interface{})
			if ok {
				timestampField, ok := statsJsonMap["timestamp"].(string)
				if ok {
					timestamp, err := time.Parse(time.RFC3339Nano, timestampField)
					if err != nil {
						log.Error("Parse timestamp error %s", stats)
						return errors.New("Parse timestamp error")
					} else {
						index := getDocumentIndex(namespace)
						documentType := getDocumentType(replicationControllerName)
						id := getDocumentID(podName, containerName, timestamp)
						jsonMap["searchMetaData"] = make(map[string]interface{})
						jsonMap["searchMetaData"].(map[string]interface{})["namespace"] = namespace
						jsonMap["searchMetaData"].(map[string]interface{})["replicationControllerName"] = replicationControllerName
						jsonMap["searchMetaData"].(map[string]interface{})["podName"] = podName
						jsonMap["searchMetaData"].(map[string]interface{})["containerName"] = containerName
						jsonMap["stats"] = stats

						err := saveContainerRecord(index, documentType, id, jsonMap)

						if err != nil {
							log.Error("Save error %s", err)
							return err
						}
					}

				} else {
					log.Error("Container doesn't have timestamp field in stats %s", stats)
					return errors.New("Container doesn't have timestamp field in stats")
				}
			} else {
				log.Error("Container doesn't have stats field as map %s", stats)
				return errors.New("Container doesn't have stats field as map")
			}
		}

		return nil
	} else {
		log.Error("Container doesn't have stats field")
		return errors.New("Container doesn't have stats field")
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
