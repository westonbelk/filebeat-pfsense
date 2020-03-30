package add_pf_metadata

import (
	"fmt"
	"sync"
	"time"
	"regexp"
	"log"

	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

type pfMetadata struct {
	config  Config
	rules common.MapStrPointer
	lastUpdate struct {
		time.Time
		sync.Mutex
	}
}

const (
	processorName = "add_pf_metadata"
)

var rule_number_re = regexp.MustCompile(`^(\d+),`)

func init() {
	processors.RegisterPlugin(processorName, New)
}


// New constructs a new add_host_metadata processor.
func New(cfg *common.Config) (processors.Processor, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return nil, errors.Wrapf(err, "fail to unpack the %v configuration", processorName)
	}

	p := &pfMetadata{
		config: config,
		rules: common.NewMapStrPointer(nil),
	}
	p.rules.Set(getRuleset())

	return p, nil
}

// Run enriches the given event with the host meta data
func (p *pfMetadata) Run(event *beat.Event) (*beat.Event, error) {
	fl_message, err := event.Fields.GetValue("message")
	if err != nil {
		log.Println("Message was not set in event:", err)
		log.Println(event.Fields.String())
		// log.Fatal("Rule number was not set in event:", err)
	}

	// Regex to pull out the rule_number
	matches := rule_number_re.FindStringSubmatch(fl_message.(string))
	if matches == nil {
		log.Fatal("Rule number was not foudn in event")
	}
	rule_number := matches[1]

	data, err := p.loadRuleInfo(rule_number)
	if err != nil {
		return nil, err
	}

	// Update the fields
	event.Fields.Update(data)

	return event, nil
}

func (p *pfMetadata) expired() bool {
	if p.config.CacheTTL <= 0 {
		return true
	}

	p.lastUpdate.Lock()
	defer p.lastUpdate.Unlock()

	if p.lastUpdate.Add(p.config.CacheTTL).After(time.Now()) {
		return false
	}
	p.lastUpdate.Time = time.Now()
	return true
}

func (p *pfMetadata) loadRuleInfo(rule_number string) (common.MapStr, error) {
	if p.expired() {
		p.rules.Set(getRuleset())
		log.Println("Updated ruleset")
	}

	rule_info, err := p.rules.Get().GetValue(rule_number)
	if err != nil {
		return common.MapStr{}, err
	}
	return rule_info.(common.MapStr).Clone(), nil
}

func (p *pfMetadata) String() string {
	return fmt.Sprintf("%v=[cache.ttl=[%v],rules=[N/I]]",
		processorName, p.config.CacheTTL)
}