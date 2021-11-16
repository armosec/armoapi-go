package armotypes

import "fmt"

// context attributes based structure to get more flexible and searchable options
type ArmoContext struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	Source    string `json:"source"`
	// open question: should we suppport in "not exists" or "not equal". E.g. "apply this recommendation only for non GCP hosted K8s clusters"
}

func DesignatorToArmoContext(designator *PortalDesignator, designatorPrefix string) []ArmoContext {
	ctxSlice := make([]ArmoContext, 0, len(designator.Attributes))
	sourceStr := fmt.Sprintf("%s.attributes", designatorPrefix)
	if designatorPrefix == "" {
		sourceStr = "attributes"
	}
	for attr := range designator.Attributes {
		ctxSlice = append(ctxSlice, ArmoContext{Attribute: attr, Value: designator.Attributes[attr], Source: sourceStr})
	}
	if designator.WLID != "" {
		ctxSlice = append(ctxSlice, ArmoContext{Attribute: "wlid", Value: designator.WLID, Source: designatorPrefix})
	}
	if designator.WildWLID != "" {
		ctxSlice = append(ctxSlice, ArmoContext{Attribute: "wildwlid", Value: designator.WildWLID, Source: designatorPrefix})
	}
	if designator.SID != "" {
		ctxSlice = append(ctxSlice, ArmoContext{Attribute: "sid", Value: designator.SID, Source: designatorPrefix})
	}
	return ctxSlice
}

// checks if all the context values match in designators
func IsDesignatorsMatchContext(ctxSlice []ArmoContext, designator *PortalDesignator, designatorPrefix string) bool {
	sourceStr := fmt.Sprintf("%s.attributes", designatorPrefix)
	if designatorPrefix == "" {
		sourceStr = "attributes"
	}
	desiredMatches := len(ctxSlice)
	for ctxElemIdx := range ctxSlice {
		switch ctxSlice[ctxElemIdx].Source {
		case sourceStr:
			if val, ok := designator.Attributes[ctxSlice[ctxElemIdx].Attribute]; !ok || val != ctxSlice[ctxElemIdx].Value {
				return false
			}
			desiredMatches--
		case designatorPrefix:
			switch ctxSlice[ctxElemIdx].Attribute {
			case "wlid":
				if ctxSlice[ctxElemIdx].Value == designator.WLID {
					desiredMatches--
				}
			case "wildwlid":
				if ctxSlice[ctxElemIdx].Value == designator.WildWLID {
					desiredMatches--
				}
			case "sid":
				if ctxSlice[ctxElemIdx].Value == designator.SID {
					desiredMatches--
				}
			}
		default:
			// not a designator attribute
			desiredMatches--
		}

	}
	return desiredMatches == 0
}
