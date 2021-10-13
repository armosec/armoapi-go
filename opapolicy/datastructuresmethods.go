package opapolicy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/golang/glog"
	"github.com/open-policy-agent/opa/rego"
)

func (pn *PolicyNotification) ToJSONBytesBuffer() (*bytes.Buffer, error) {
	res, err := json.Marshal(pn)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(res), err
}

func (RuleResponse *RuleResponse) GetSingleResultStatus() string {
	if RuleResponse.Exception != nil {
		if RuleResponse.Exception.IsAlertOnly() {
			return "warning"
		}
		if RuleResponse.Exception.IsDisable() {
			return "ignore"
		}
	}
	return "failed"

}

func (ruleReport *RuleReport) GetRuleStatus() (string, []RuleResponse, []RuleResponse) {
	if len(ruleReport.RuleResponses) == 0 {
		return "success", nil, nil
	}
	exceptions := make([]RuleResponse, 0)
	failed := make([]RuleResponse, 0)

	for _, rule := range ruleReport.RuleResponses {
		if rule.ExceptionName != "" {
			exceptions = append(exceptions, rule)
		} else if rule.Exception != nil {
			exceptions = append(exceptions, rule)
		} else {
			failed = append(failed, rule)
		}
	}

	status := "failed"
	if len(failed) == 0 && len(exceptions) > 0 {
		status = "warning"
	}
	return status, failed, exceptions
}

func (controlReport *ControlReport) GetNumberOfFailedResources() int {
	sum := 0
	for i := range controlReport.RuleReports {
		sum += controlReport.RuleReports[i].GetNumberOfFailedResources()
	}
	return sum
}

// func (controlReport *ControlReport) GetNumberOfResources() int {
// 	sum := 0
// 	for i := range controlReport.RuleReports {
// 		sum += controlReport.RuleReports[i].GetNumberOfResources()
// 	}
// 	return sum
// }

func ParseRegoResult(regoResult *rego.ResultSet) ([]RuleResponse, error) {
	var errs error
	ruleResponses := []RuleResponse{}
	for _, result := range *regoResult {
		for desicionIdx := range result.Expressions {
			if resMap, ok := result.Expressions[desicionIdx].Value.(map[string]interface{}); ok {
				for objName := range resMap {
					jsonBytes, err := json.Marshal(resMap[objName])
					if err != nil {
						err = fmt.Errorf("in parseRegoResult, json.Marshal failed. name: %s, obj: %v, reason: %s", objName, resMap[objName], err)
						glog.Error(err)
						errs = fmt.Errorf("%s\n%s", errs, err)
						continue
					}
					desObj := make([]RuleResponse, 0)
					if err := json.Unmarshal(jsonBytes, &desObj); err != nil {
						err = fmt.Errorf("in parseRegoResult, json.Unmarshal failed. name: %s, obj: %v, reason: %s", objName, resMap[objName], err)
						glog.Error(err)
						errs = fmt.Errorf("%s\n%s", errs, err)
						continue
					}
					ruleResponses = append(ruleResponses, desObj...)
				}
			}
		}
	}
	return ruleResponses, errs
}

func (controlReport *ControlReport) GetNumberOfResources() int {
	sum := 0
	for i := range controlReport.RuleReports {
		if controlReport.RuleReports[i].ListInputResources != nil {
			sum += len(controlReport.RuleReports[i].ListInputResources)
		}
	}
	return sum
}

func (ruleReport *RuleReport) GetNumberOfWarningResources() int {
	sum := 0
	for i := range ruleReport.RuleResponses {
		if ruleReport.RuleResponses[i].GetSingleResultStatus() == "warning" {
			sum += len(ruleReport.RuleResponses[i].AlertObject.K8SApiObjects)
		}
	}
	return sum
}

func (controlReport *ControlReport) GetNumberOfWarningResources() int {
	sum := 0
	for i := range controlReport.RuleReports {
		sum += controlReport.RuleReports[i].GetNumberOfWarningResources()
	}
	return sum
}
func (controlReport *ControlReport) ListControlsInputKinds() []string {
	listControlsInputKinds := []string{}
	for i := range controlReport.RuleReports {
		listControlsInputKinds = append(listControlsInputKinds, controlReport.RuleReports[i].ListInputKinds...)
	}
	return listControlsInputKinds
}

func (controlReport *ControlReport) Passed() bool {
	for i := range controlReport.RuleReports {
		if len(controlReport.RuleReports[i].RuleResponses) != 0 {
			return false
		}
	}
	return true
}

func (controlReport *ControlReport) Warning() bool {
	if controlReport.Passed() || controlReport.Failed() {
		return false
	}
	for i := range controlReport.RuleReports {
		if status, _, _ := controlReport.RuleReports[i].GetRuleStatus(); status == "warning" {
			return true
		}
	}
	return false
}

func (controlReport *ControlReport) Failed() bool {
	if controlReport.Passed() {
		return false
	}
	for i := range controlReport.RuleReports {
		if status, _, _ := controlReport.RuleReports[i].GetRuleStatus(); status == "failed" {
			return true
		}
	}
	return false
}

func (ruleReport *RuleReport) GetNumberOfResources() int {
	return len(ruleReport.ListInputResources)
}

func (ruleReport *RuleReport) GetNumberOfFailedResources() int {
	sum := 0
	for i := len(ruleReport.RuleResponses) - 1; i >= 0; i-- {
		if ruleReport.RuleResponses[i].GetSingleResultStatus() == "failed" {
			sum += len(ruleReport.RuleResponses[i].AlertObject.K8SApiObjects)
		}
	}
	return sum
}
func (ctrl *ControlReport) GetID() string {
	h := fnv.New32a()
	h.Write([]byte(ctrl.Name))
	s := fmt.Sprintf("%d", h.Sum32())

	return "C-" + s
}

func RemoveResponse(slice []RuleResponse, index int) []RuleResponse {
	return append(slice[:index], slice[index+1:]...)
}

// TODO - receive list full json paths
func (postureReport *PostureReport) RemoveData(keepFields, keepMetadataFields []string) {
	for i := range postureReport.FrameworkReports {
		postureReport.FrameworkReports[i].RemoveData(keepFields, keepMetadataFields)
	}
}
func (frameworkReport *FrameworkReport) RemoveData(keepFields, keepMetadataFields []string) {
	for i := range frameworkReport.ControlReports {
		frameworkReport.ControlReports[i].RemoveData(keepFields, keepMetadataFields)
	}
}
func (controlReport *ControlReport) RemoveData(keepFields, keepMetadataFields []string) {
	for i := range controlReport.RuleReports {
		controlReport.RuleReports[i].RemoveData(keepFields, keepMetadataFields)
	}
}

func (ruleReport *RuleReport) RemoveData(keepFields, keepMetadataFields []string) {
	for i := range ruleReport.RuleResponses {
		ruleReport.RuleResponses[i].RemoveData(keepFields, keepMetadataFields)
	}
}

func (r *RuleResponse) RemoveData(keepFields, keepMetadataFields []string) {
	r.AlertObject.ExternalObjects = nil

	for i := range r.AlertObject.K8SApiObjects {
		deleteFromMap(r.AlertObject.K8SApiObjects[i], keepFields)
		for k := range r.AlertObject.K8SApiObjects[i] {
			if k == "metadata" {
				if b, ok := r.AlertObject.K8SApiObjects[i][k].(map[string]interface{}); ok {
					deleteFromMap(b, keepMetadataFields)
					r.AlertObject.K8SApiObjects[i][k] = b
				}
			}
		}
	}
}

func deleteFromMap(m map[string]interface{}, keepFields []string) {
	for k := range m {
		if StringInSlice(keepFields, k) {
			continue
		}
		delete(m, k)
	}
}

func StringInSlice(strSlice []string, str string) bool {
	for i := range strSlice {
		if strSlice[i] == str {
			return true
		}
	}
	return false
}
