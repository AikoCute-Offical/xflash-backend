package core

import (
	"github.com/AikoCute-Offical/xflash-backend/api/xflash"
)

func (p *Core) UpdateRule(tag string, newRuleList []xflash.DetectRule) error {
	return p.dispatcher.RuleManager.UpdateRule(tag, newRuleList)
}

func (p *Core) UpdateProtocolRule(tag string, newRuleList []string) error {

	return p.dispatcher.RuleManager.UpdateProtocolRule(tag, newRuleList)
}

func (p *Core) GetDetectResult(tag string) ([]xflash.DetectResult, error) {
	return p.dispatcher.RuleManager.GetDetectResult(tag)
}
