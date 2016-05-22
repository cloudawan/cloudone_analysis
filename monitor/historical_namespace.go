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
	"github.com/cloudawan/cloudone_analysis/control"
	"github.com/cloudawan/cloudone_utility/logger"
)

func RecordHistoricalAllNamespace(kubeApiServerEndPoint string, kubeApiServerToken string) (returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("RecordHistoricalAllNamespace Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedError = err.(error)
		}
	}()

	namespaceNameSlice, err := control.GetAllNamespaceName(kubeApiServerEndPoint, kubeApiServerToken)
	if err != nil {
		log.Error(err)
		return err
	}

	allNamespaceContainerRecordSlice := make([]map[string]interface{}, 0)
	for _, namespaceName := range namespaceNameSlice {
		replicationControllerNameSlice, err := control.GetAllReplicationControllerName(kubeApiServerEndPoint, kubeApiServerToken, namespaceName)
		if err != nil {
			log.Error(err)
		} else {
			for _, replicationControllerName := range replicationControllerNameSlice {
				replicationControllerContainerRecordSlice, err := RecordHistoricalReplicationController(kubeApiServerEndPoint, kubeApiServerToken, namespaceName, replicationControllerName)
				if err != nil {
					log.Error(err)
				} else {
					allNamespaceContainerRecordSlice = append(allNamespaceContainerRecordSlice, replicationControllerContainerRecordSlice...)
				}
			}
		}
	}

	for _, containerRecord := range allNamespaceContainerRecordSlice {
		index, _ := containerRecord["searchMetaData"].(map[string]interface{})["index"].(string)
		documentType, _ := containerRecord["searchMetaData"].(map[string]interface{})["documentType"].(string)
		id, _ := containerRecord["searchMetaData"].(map[string]interface{})["id"].(string)
		if err := saveContainerRecord(index, documentType, id, containerRecord); err != nil {
			log.Error("Save error %s", err)
		}
	}

	return nil
}
