package add_pf_metadata

import (
	"github.com/elastic/beats/libbeat/common"
	"fmt"
	"os/exec"
	"bytes"
	"log"
	"regexp"
)

var ruleset_exp = regexp.MustCompile(`(?m)^@(?P<rule_number>\d+)(?:\((?P<tracker_number>\d+)\))? (?P<effective_rule>.+?)(?: label \"(?P<label>.*)\")?$`)

func parseRules(ruleset string) common.MapStr {
	rules := common.MapStr{}

	matches := ruleset_exp.FindAllStringSubmatch(ruleset, -1)
	
	if matches == nil { 
		return rules
	}

	for _, match := range matches {
		var rule = match[1]
		rules.Put(rule+".pf.rule_info.tracker_number", match[2])
		rules.Put(rule+".pf.rule_info.effective_rule", match[3])
		rules.Put(rule+".pf.rule_info.label", match[4])
	}
	return rules
}

func getRuleset() common.MapStr {
	cmd := exec.Command("pfctl", "-gsr")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(stderr.String())
	}
	return parseRules(stdout.String())
}

func main() {
	event := common.MapStr{}
	event.Put("pf.rule", "0")
	ruleset := getRuleset()

	rule_number, err := event.GetValue("pf.rule")
	if err != nil {
		log.Fatal("Rule number was not set in event:", err)
	}

	rule_info, err := ruleset.GetValue(rule_number.(string))
	if err != nil {
		log.Fatal("Requested rule not found in ruleset:", err)
	}
	
	event.Update(rule_info.(common.MapStr).Clone())

	fmt.Println(event.StringToPrint())
}