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
		Do(returns200, returns404, returns500).
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
		Do(returns200BuildLog, returns400, returns404, returns500))

	ws.Route(ws.DELETE("/{imageinformation}/{version}").Filter(authorize).Filter(auditLog).To(deleteBuildLog).
		Doc("Delete the build logs belonging to the image information with the version").
		Param(ws.PathParameter("imageinformation", "Image information").DataType("string")).
		Param(ws.PathParameter("version", "Version").DataType("string")).
		Do(returns200, returns400, returns404, returns500))
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

	buildLogSlice, err := build.SearchBuildLog(imageInformation, from, to, size, offset)
	if err != nil {
		errorText := fmt.Sprintf("Fail to get all build logs belonging to image information %s logs with error %s", imageInformation, err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}

	response.WriteJson(buildLogSlice, "[]BuildLog")
}

func postBuildLog(request *restful.Request, response *restful.Response) {
	buildLog := &utilitybuild.BuildLog{}
	err := request.ReadEntity(&buildLog)

	if err != nil {
		errorText := fmt.Sprintf("POST Build Log with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	err = build.SaveBuildLog(buildLog, false)
	if err != nil {
		errorText := fmt.Sprintf("Fail to create build log with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}
}

func deleteBuildLogBelongingToImageInformation(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")

	err := build.DeleteBuildLogBelongingToImageInformation(imageInformation)
	if err != nil {
		errorText := fmt.Sprintf("Fail to delete build logs belonging to image information %s with error %s", imageInformation, err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}
}

func getBuildLog(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")
	version := request.PathParameter("version")

	buildLog, err := build.GetBuildLog(imageInformation, version)
	if err != nil {
		errorText := fmt.Sprintf("Fail to delete build logs belonging to image information %s with version %s with error %s", imageInformation, version, err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}

	response.WriteJson(buildLog, "BuildLog")
}

func deleteBuildLog(request *restful.Request, response *restful.Response) {
	imageInformation := request.PathParameter("imageinformation")
	version := request.PathParameter("version")

	err := build.DeleteBuildLog(imageInformation, version)
	if err != nil {
		errorText := fmt.Sprintf("Fail to delete build logs belonging to image information %s with version %s with error %s", imageInformation, version, err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}
}

func returns200BuildLogSlice(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", []utilitybuild.BuildLog{})
}

func returns200BuildLog(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", utilitybuild.BuildLog{})
}
