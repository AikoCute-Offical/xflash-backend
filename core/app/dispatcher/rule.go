package dispatcher

import (
	"github.com/AikoCute-Offical/xflash-backend/api/panel"
	"reflect"
	"sync"
)

type Rule struct {
	Rule *sync.Map // Key: Tag, Value: *panel.DetectRule
}

func NewRule() *Rule {
	return &Rule{
		Rule: new(sync.Map),
	}
}

func (r *Rule) UpdateRule(tag string, newRuleList *panel.DetectRule) error {
	if value, ok := r.Rule.LoadOrStore(tag, newRuleList); ok {
		oldRuleList := value.([]panel.DestinationRule)
		if !reflect.DeepEqual(oldRuleList, newRuleList) {
			r.Rule.Store(tag, newRuleList)
		}
	}
	return nil
}

func (r *Rule) Detect(tag string, destination string, protocol string) (reject bool) {
	reject = false
	// If we have some rule for this inbound
	if value, ok := r.Rule.Load(tag); ok {
		ruleList := value.(*panel.DetectRule)
		for i, _ := range ruleList.DestinationRule {
			if ruleList.DestinationRule[i].Pattern.Match([]byte(destination)) {
				reject = true
				break
			}
		}
		if !reject {
			for _, v := range ruleList.ProtocolRule {
				if v == protocol {
					return true
				}
			}
		}
	}
	return reject
}
