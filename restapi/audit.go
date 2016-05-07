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
	"github.com/cloudawan/cloudone_analysis/audit"
	utilityaudit "github.com/cloudawan/cloudone_utility/audit"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"time"
)

func registerWebServiceAuditLog() {
	ws := new(restful.WebService)
	ws.Path("/api/v1/auditlogs")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	restful.Add(ws)

	ws.Route(ws.GET("/").Filter(authorize).Filter(auditLog).To(getAllAuditLog).
		Doc("Get audit logs in the time range").
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("size", "The amount of data to return").DataType("int")).
		Param(ws.QueryParameter("offset", "The offset from the result").DataType("int")).
		Do(returns200AuditLogSlice, returns400, returns404, returns500))

	// Don't audit itself to prevent loop. Also, this is used only by system
	ws.Route(ws.POST("/").Filter(authorize).To(postAuditLog).
		Doc("Create the audit log").
		Do(returns200, returns404, returns500).
		Reads(utilityaudit.AuditLog{}))

	ws.Route(ws.GET("/{user}").Filter(authorize).Filter(auditLog).To(getAuditLog).
		Doc("Get the audit logs belonging to the user").
		Param(ws.PathParameter("user", "User name").DataType("string")).
		Param(ws.QueryParameter("from", "Time start from in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("to", "Time end to in RFC3339Nano formt").DataType("string")).
		Param(ws.QueryParameter("size", "The amount of data to return").DataType("int")).
		Param(ws.QueryParameter("offset", "The offset from the result").DataType("int")).
		Do(returns200AuditLogSlice, returns400, returns404, returns500))
}

func getAllAuditLog(request *restful.Request, response *restful.Response) {
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

	auditLogSlice, err := audit.SearchAuditLog("*", from, to, size, offset)
	if err != nil {
		errorText := fmt.Sprintf("Fail to get all audit logs with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}

	response.WriteJson(auditLogSlice, "[]AuditLog")
}

func getAuditLog(request *restful.Request, response *restful.Response) {
	user := request.PathParameter("user")
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

	auditLogSlice, err := audit.SearchAuditLog(user, from, to, size, offset)
	if err != nil {
		errorText := fmt.Sprintf("Fail to get all audit logs belonging to user %s logs with error %s", user, err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}

	response.WriteJson(auditLogSlice, "[]AuditLog")
}

func postAuditLog(request *restful.Request, response *restful.Response) {
	auditLog := &utilityaudit.AuditLog{}
	err := request.ReadEntity(&auditLog)

	if err != nil {
		errorText := fmt.Sprintf("POST Audit Log with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(400, `{"Error": "`+errorText+`"}`)
		return
	}

	err = audit.SaveAudit(auditLog, false)
	if err != nil {
		errorText := fmt.Sprintf("Fail to create audit log with error %s", err)
		log.Error(errorText)
		response.WriteErrorString(404, `{"Error": "`+errorText+`"}`)
		return
	}
}

func returns200AuditLogSlice(b *restful.RouteBuilder) {
	b.Returns(http.StatusOK, "OK", []utilityaudit.AuditLog{})
}
