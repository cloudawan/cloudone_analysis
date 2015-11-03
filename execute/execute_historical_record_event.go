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

package execute

import (
	"github.com/cloudawan/cloudone_analysis/event"
	"github.com/cloudawan/cloudone_analysis/utility/configuration"
	"github.com/cloudawan/cloudone_utility/logger"
	"time"
)

func loopHistoricalRecordEvent(ticker *time.Ticker, checkingInterval time.Duration) {
	for {
		select {
		case <-ticker.C:
			// Historical record
			periodicalRunHistoricalRecordEvent()
		case <-quitChannel:
			ticker.Stop()
			log.Info("Loop historical record event quit")
			return
		}
	}
}

func periodicalRunHistoricalRecordEvent() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("periodicalRunHistoricalRecordEvent Error: %s", err)
			log.Error(logger.GetStackTrace(4096, false))
		}
	}()

	kubeapiHost, ok := configuration.LocalConfiguration.GetString("kubeapiHost")
	if ok == false {
		log.Error("Fail to get configuration kubeapiHost")
		return
	}
	kubeapiPort, ok := configuration.LocalConfiguration.GetInt("kubeapiPort")
	if ok == false {
		log.Error("Fail to get configuration kubeapiPort")
		return
	}

	if err := event.RecordHistoricalEvent(kubeapiHost, kubeapiPort); err != nil {
		log.Error(err)
		return
	}
}
