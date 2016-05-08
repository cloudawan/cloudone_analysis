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
	"strconv"
	"time"
)

func registerWebServiceHistoricalReplicationControllerMetric() {
	ws := new(restful.WebService)
	ws.Path("/api/v1/historicalreplicationcontrollermetrics")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	restful.Add(ws)

	ws.Route(ws.GET("/{namespace}").Filter(authorize).Filter(auditLog).To(getAllHistoricalReplicationControllerMetric).
		Doc("Get all historical replication controllers in the namespace").
		Param(ws.PathParameter("namespace", "Kubernetes namespace").DataType("string")).
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("aggregationAmount", "Aggregation amount").DataType("int")).
		Do(returns200JsonMap, returns400, returns404, returns500))

	ws.Route(ws.GET("/{namespace}/{replicationcontroller}").Filter(authorize).Filter(auditLog).To(getHistoricalReplicationControllerMetric).
		Doc("Get the historical replication controller in the namespace").
		Param(ws.PathParameter("namespace", "Kubernetes namespace").DataType("string")).
		Param(ws.PathParameter("replicationcontroller", "Kubernetes replication controller name").DataType("string")).
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("aggregationAmount", "Aggregation amount").DataType("int")).
		Do(returns200JsonMap, returns400, returns404, returns500))
}

func getAllHistoricalReplicationControllerMetric(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	fromText := request.QueryParameter("from")
	toText := request.QueryParameter("to")
	aggregationAmountText := request.QueryParameter("aggregationAmount")

	from, err := time.Parse(time.RFC3339Nano, fromText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse fromText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["fromText"] = fromText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	to, err := time.Parse(time.RFC3339Nano, toText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse toText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["toText"] = toText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	aggregationAmount, err := strconv.Atoi(aggregationAmountText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse aggregationAmountText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["aggregationAmountText"] = aggregationAmountText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	jsonMap, err := monitor.GetAllHistoricalReplicationControllerMetrics(
		namespace, aggregationAmount, from, to)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Get historical replication controller metrics with the criteria failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["namespace"] = namespace
		jsonMap["from"] = from
		jsonMap["to"] = to
		jsonMap["aggregationAmount"] = aggregationAmount
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}

	response.WriteJson(jsonMap, "Json")
}

func getHistoricalReplicationControllerMetric(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	replicationControllerName := request.PathParameter("replicationcontroller")
	fromText := request.QueryParameter("from")
	toText := request.QueryParameter("to")
	aggregationAmountText := request.QueryParameter("aggregationAmount")

	from, err := time.Parse(time.RFC3339Nano, fromText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse fromText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["fromText"] = fromText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	to, err := time.Parse(time.RFC3339Nano, toText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse toText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["toText"] = toText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	aggregationAmount, err := strconv.Atoi(aggregationAmountText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse aggregationAmountText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["aggregationAmountText"] = aggregationAmountText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	jsonMap, err := monitor.GetHistoricalReplicationControllerMetrics(
		namespace, replicationControllerName, aggregationAmount, from, to)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Get historical replication controller metrics with the criteria failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["namespace"] = namespace
		jsonMap["replicationControllerName"] = replicationControllerName
		jsonMap["from"] = from
		jsonMap["to"] = to
		jsonMap["aggregationAmount"] = aggregationAmount
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}

	response.WriteJson(jsonMap, "Json")
}
