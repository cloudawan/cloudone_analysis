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

package control

import (
	"github.com/cloudawan/kubernetes_management_utility/logger"
	"github.com/cloudawan/kubernetes_management_utility/restclient"
	"strconv"
)

func GetAllNamespaceName(kubeapiHost string, kubeapiPort int) (returnedNameSlice []string, returnedError error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("GetAllNamespaceName Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
			returnedNameSlice = nil
			returnedError = err.(error)
		}
	}()

	jsonMap, err := restclient.RequestGet("http://"+kubeapiHost+":"+strconv.Itoa(kubeapiPort)+"/api/v1/namespaces/", true)
	if err != nil {
		log.Error("Fail to get all namespace name with host: %s, port: %d, error: %s", kubeapiHost, kubeapiPort, err.Error())
		return nil, err
	}

	nameSlice := make([]string, 0)
	for _, data := range jsonMap.(map[string]interface{})["items"].([]interface{}) {
		name, ok := data.(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		if ok {
			nameSlice = append(nameSlice, name)
		}
	}

	return nameSlice, nil
}
