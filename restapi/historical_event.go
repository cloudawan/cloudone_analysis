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
	"fmt"
	"github.com/cloudawan/kubernetes_management_analysis/event"
	"github.com/emicklei/go-restful"
	"strconv"
	"time"
)

func registerWebServiceHistoricalEvent() {
	ws := new(restful.WebService)
	ws.Path("/api/v1/historicalevents")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	restful.Add(ws)

	ws.Route(ws.GET("/").To(getAllHistoricalEvent).
		Doc("Get all historical events").
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("acknowledge", "Already acknowledged or not").DataType("boolean")).
		Param(ws.QueryParameter("size", "The amount of data to return").DataType("int")).
		Param(ws.QueryParameter("offset", "The offset from the result").DataType("int")).
		Do(returns200JsonMap, returns400, returns404, returns500))

	ws.Route(ws.GET("/{namespace}").To(getHistoricalEvent).
		Doc("Get the historical events in the namespace").
		Param(ws.PathParameter("namespace", "Kubernetes namespace").DataType("string")).
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("acknowledge", "Already acknowledged or not").DataType("boolean")).
		Param(ws.QueryParameter("size", "The amount of data to return").DataType("int")).
		Param(ws.QueryParameter("offset", "The offset from the result").DataType("int")).
		Do(returns200JsonMap, returns400, returns404, returns500))

	ws.Route(ws.PUT("/{namespace}/{id}").To(acknowledgeHistoricalEvent).
		Doc("Acknowledge the historical events in the namespace").
		Param(ws.PathParameter("namespace", "Kubernetes namespace").DataType("string")).
		Param(ws.PathParameter("id", "Kubernetes event id").DataType("string")).
		Param(ws.QueryParameter("acknowledge", "acknowledge or unacknowledge").DataType("boolean")).
		Do(returns200, returns400, returns404, returns500))
}

func getAllHistoricalEvent(request *restful.Request, response *restful.Response) {
	fromText := request.QueryParameter("from")
	toText := request.QueryParameter("to")
	acknowledgeText := request.QueryParameter("acknowledge")
	sizeText := request.QueryParameter("size")
	offsetText := request.QueryParameter("offset")

	var from *time.Time
	if fromText == "" {
		from = nil
	} else {
		fromValue, err := time.Parse(time.RFC3339Nano, fromText)
		if err != nil {
			errorText := fmt.Sprintf("Parse from %s with error %s", fromText, err)
			log.Error(errorText)
			response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
			return
		} else {
			from = &fromValue
		}
	}

	var to *time.Time
	if toText == "" {
		to = nil
	} else {
		toValue, err := time.Parse(time.RFC3339Nano, toText)
		if err != nil {
			errorText := fmt.Sprintf("Parse from %s with error %s", toText, err)
			log.Error(errorText)
			response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
			return
		} else {
			to = &toValue
		}
	}

	acknowledge, err := strconv.ParseBool(acknowledgeText)
	if err != nil {
		errorText := fmt.Sprintf("Parse from %s with error %s", acknowledgeText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	size, err := strconv.Atoi(sizeText)
	if err != nil {
		errorText := fmt.Sprintf("Parse from %s with error %s", sizeText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	offset, err := strconv.Atoi(offsetText)
	if err != nil {
		errorText := fmt.Sprintf("Parse from %s with error %s", offsetText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	jsonMap, err := event.SearchHistoricalEvent("*", from, to, acknowledge, size, offset)
	if err != nil {
		errorText := fmt.Sprintf("Fail to get all historical events with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}

	response.WriteJson(jsonMap, "Json")
}

func getHistoricalEvent(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	fromText := request.QueryParameter("from")
	toText := request.QueryParameter("to")
	acknowledgeText := request.QueryParameter("acknowledge")
	sizeText := request.QueryParameter("size")
	offsetText := request.QueryParameter("offset")

	var from *time.Time
	if fromText == "" {
		from = nil
	} else {
		fromValue, err := time.Parse(time.RFC3339Nano, fromText)
		if err != nil {
			errorText := fmt.Sprintf("Parse from %s with error %s", fromText, err)
			log.Error(errorText)
			response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
			return
		} else {
			from = &fromValue
		}
	}

	var to *time.Time
	if toText == "" {
		to = nil
	} else {
		toValue, err := time.Parse(time.RFC3339Nano, toText)
		if err != nil {
			errorText := fmt.Sprintf("Parse from %s with error %s", toText, err)
			log.Error(errorText)
			response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
			return
		} else {
			to = &toValue
		}
	}

	acknowledge, err := strconv.ParseBool(acknowledgeText)
	if err != nil {
		errorText := fmt.Sprintf("Parse from %s with error %s", acknowledgeText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	size, err := strconv.Atoi(sizeText)
	if err != nil {
		errorText := fmt.Sprintf("Parse from %s with error %s", sizeText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	offset, err := strconv.Atoi(offsetText)
	if err != nil {
		errorText := fmt.Sprintf("Parse from %s with error %s", offsetText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	jsonSlice, err := event.SearchHistoricalEvent(namespace, from, to, acknowledge, size, offset)
	if err != nil {
		errorText := fmt.Sprintf("Fail to get all historical events with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}

	response.WriteJson(jsonSlice, "Json")
}

func acknowledgeHistoricalEvent(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	id := request.PathParameter("id")
	acknowledgeText := request.QueryParameter("acknowledge")

	acknowledge, err := strconv.ParseBool(acknowledgeText)
	if err != nil {
		errorText := fmt.Sprintf("Parse to %s with error %s", acknowledgeText, err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	err = event.Acknowledge(namespace, id, acknowledge)
	if err != nil {
		errorText := fmt.Sprintf("Fail to acknowledge historical events with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}
}
