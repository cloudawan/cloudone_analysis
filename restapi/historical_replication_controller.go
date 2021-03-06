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

package restapi

import (
	"encoding/json"
	"github.com/cloudawan/cloudone_analysis/monitor"
	"github.com/emicklei/go-restful"
	"net/http"
)

func registerWebServiceHistoricalReplicationController() {
	ws := new(restful.WebService)
	ws.Path("/api/v1/historicalreplicationcontrollers")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	restful.Add(ws)

	ws.Route(ws.GET("/names/{namespace}").Filter(authorize).Filter(auditLog).To(getAllHistoricalReplicationControllerName).
		Doc("Get all historical replication controller names in the namespace").
		Param(ws.PathParameter("namespace", "Kubernetes namespace").DataType("string")).
		Do(returns200StringSlice, returns404, returns500))
}

func getAllHistoricalReplicationControllerName(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")

	nameSlice, err := monitor.GetAllReplicationControllerNameInNameSpace(namespace)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Get all hittorical replication controller name in the namespace failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["namespace"] = namespace
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}

	response.WriteJson(nameSlice, "[]String")
}

func returns200StringSlice(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", make([]string, 0))
}
