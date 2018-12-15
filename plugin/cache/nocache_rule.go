package cache

import (
	"context"
	"regexp"
	"time"
)

// if the URI match rules then will not enable cache plugin

var (
	regexps   []*regexp.Regexp
	cntRegexp int
)

// Rule ...
type Rule struct {
	Regular string `json:"regular"`
}

func initRules(rules []Rule) {
	regexps = make([]*regexp.Regexp, len(rules))
	for idx, rule := range rules {
		regexps[idx], _ = regexp.Compile(rule.Regular)
	}

	cntRegexp = len(rules)
}

// match NocacheRule
func matchNoCacheRule(uri string) bool {
	var (
		checkChan = make(chan bool, cntRegexp)
		counter   int
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer close(checkChan)

	for _, reg := range regexps {
		// fmt.Printf("reg: %s matched\n", reg.String())
		go func(ctx context.Context, reg *regexp.Regexp, c chan<- bool) {
			// to catch send on close channel
			go func() { recover() }()

			select {
			case <-ctx.Done():
				println("timeout matchNoCacheRule")
				break
			default:
				c <- reg.MatchString(uri)
			}
		}(ctx, reg, checkChan)
	}

	for checked := range checkChan {
		// fmt.Printf("%s, %d, %v\n", uri, cntRegexp, checked)
		if checked {
			return checked
		}
		// counter to mark all gorountine called finished
		counter++
		if counter >= cntRegexp {
			break
		}
	}
	return false
}
