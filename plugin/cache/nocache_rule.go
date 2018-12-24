package cache

import (
	"context"
	"regexp"
	"time"

	"github.com/yeqown/gateway/config/rule"
)

// no cache rule settings, if the URI macthed any rule in rules
// then abort cache plugin processing
func (c *Cache) load(rules []rule.Nocacher) {
	c.regexps = make([]*regexp.Regexp, len(rules))
	for idx, r := range rules {
		c.regexps[idx], _ = regexp.Compile(r.Regular())
	}

	c.cntRegexp = len(rules)
}

// match NocacheRule, true means no cache
// fasle means cache
func (c *Cache) matchNoCacheRule(uri string) bool {
	if c.cntRegexp == 0 {
		return false
	}

	var (
		checkChan = make(chan bool, c.cntRegexp)
		counter   int
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer close(checkChan)

	for _, reg := range c.regexps {
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
		if counter >= c.cntRegexp {
			break
		}
	}
	return false
}
