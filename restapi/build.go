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
	"github.com/cloudawan/cloudone_analysis/build"
	utilitybuild "github.com/cloudawan/cloudone_utility/build"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"time"
)

func registerWebServiceBuildLog() {
	ws := new(restful.WebService)
	ws.Path("/api/v1/buildlogs")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	restful.Add(ws)

	ws.Route(ws.POST("/").Filter(authorize).Filter(auditLog).To(postBuildLog).
		Doc("Create the build log").
		Do(returns200, returns400, returns422, returns500).
		Reads(utilitybuild.BuildLog{}))

	ws.Route(ws.GET("/{imageinformation}").Filter(authorize).Filter(auditLog).To(getBuildLogBelongingToImageInformation).
		Doc("Get the build logs belonging to the image information").
		Param(ws.PathParameter("imageinformation", "Image information").DataType("string")).
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("size", "The amount of data to return").DataType("int")).
		Param(ws.QueryParameter("offset", "The offset from the result").DataType("int")).
		Do(returns200BuildLogSlice, returns400, returns404, returns500))

	ws.Route(ws.DELETE("/{imageinformation}").Filter(authorize).Filter(auditLog).To(deleteBuildLogBelongingToImageInformation).
		Doc("Delete the build logs belonging to the image information").
		Param(ws.PathParameter("imageinformation", "Image information").DataType("string")).
		Do(returns200, returns400, returns404, returns500))

	ws.Route(ws.GET("/{imageinformation}/{version}").Filter(authorize).Filter(auditLog).To(getBuildLog).
		Doc("Get the build logs belonging to the image information with the version").
		Param(ws.PathParameter("imageinformation", "Image information").DataType("string")).
		Param(ws.PathParameter("version", "Version").DataType("string")).
		Do(returns200BuildLog, returns404, returns500))

	ws.Route(ws.DELETE("/{imageinformation}/{version}").Filter(authorize).Filter(auditLog).To(deleteBuildLog).
		Doc("Delete the build logs belonging to the image information with the version").
		Param(ws.PathParameter("imageinformation", "Image information").DataType("string")).
		Param(ws.PathParameter("version", "Version").DataType("string")).
		Do(returns200, returns400, returns500))
}

func postBuildLog(request *restful.Request, response *restful.Response) {
	buildLog := &utilitybuild.BuildLog{}
	err := request.ReadEntity(&buildLog)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Read body failure"
		jsonMap["ErrorMessage"] = err.Error()
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	err = build.SaveBuildLog(buildLog, false)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Create build log failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["buildLog"] = buildLog
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(422, string(errorMessageByteSlice))
		return
	}
}

func getBuildLogBelongingToImageInformation(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")
	fromText := request.QueryParameter("from")
	toText := request.QueryParameter("to")
	sizeText := request.QueryParameter("size")
	offsetText := request.QueryParameter("offset")

	var from *time.Time
	if fromText == "" {
		from = nil
	} else {
		fromValue, err := time.Parse(time.RFC3339Nano, fromText)
		if err != nil {
			jsonMap := make(map[string]interface{})
			jsonMap["Error"] = "Could not parse fromText"
			jsonMap["ErrorMessage"] = err.Error()
			jsonMap["fromText"] = fromText
			errorMessageByteSlice, _ := json.Marshal(jsonMap)
			log.Error(jsonMap)
			response.WriteErrorString(400, string(errorMessageByteSlice))
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
			jsonMap := make(map[string]interface{})
			jsonMap["Error"] = "Could not parse toText"
			jsonMap["ErrorMessage"] = err.Error()
			jsonMap["toText"] = toText
			errorMessageByteSlice, _ := json.Marshal(jsonMap)
			log.Error(jsonMap)
			response.WriteErrorString(400, string(errorMessageByteSlice))
			return
		} else {
			to = &toValue
		}
	}

	size, err := strconv.Atoi(sizeText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse sizeText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["sizeText"] = sizeText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	offset, err := strconv.Atoi(offsetText)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Could not parse offsetText"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["offsetText"] = offsetText
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(400, string(errorMessageByteSlice))
		return
	}

	buildLogSlice, err := build.SearchBuildLog(imageInformation, from, to, size, offset)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Get build log with the criteria failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["imageInformation"] = imageInformation
		jsonMap["from"] = from
		jsonMap["to"] = to
		jsonMap["size"] = size
		jsonMap["offset"] = offset
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}

	response.WriteJson(buildLogSlice, "[]BuildLog")
}

func deleteBuildLogBelongingToImageInformation(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")

	err := build.DeleteBuildLogBelongingToImageInformation(imageInformation)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Delete build log belonging to image information failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["imageInformation"] = imageInformation
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}
}

func getBuildLog(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")
	version := request.PathParameter("version")

	buildLog, err := build.GetBuildLog(imageInformation, version)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Get build log failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["imageInformation"] = imageInformation
		jsonMap["version"] = version
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}

	response.WriteJson(buildLog, "BuildLog")
}

func deleteBuildLog(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")
	version := request.PathParameter("version")

	err := build.DeleteBuildLog(imageInformation, version)
	if err != nil {
		jsonMap := make(map[string]interface{})
		jsonMap["Error"] = "Delete build log failure"
		jsonMap["ErrorMessage"] = err.Error()
		jsonMap["imageInformation"] = imageInformation
		jsonMap["version"] = version
		errorMessageByteSlice, _ := json.Marshal(jsonMap)
		log.Error(jsonMap)
		response.WriteErrorString(404, string(errorMessageByteSlice))
		return
	}
}

func returns200BuildLogSlice(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", []utilitybuild.BuildLog{})
}

func returns200BuildLog(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", utilitybuild.BuildLog{})
}
