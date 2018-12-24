package cache

import (
	"testing"

	"github.com/yeqown/gateway/config/rule"
)

type nocacheCfg struct {
	Regexp string
}

func (c nocacheCfg) Enabled() bool   { return true }
func (c nocacheCfg) ID() string      { return "id" }
func (c nocacheCfg) SetID(string)    { return }
func (c nocacheCfg) Regular() string { return c.Regexp }

func Test_initRules(t *testing.T) {
	c := &Cache{}

	type args struct {
		rules []rule.Nocacher
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 1",
			args: args{
				rules: []rule.Nocacher{
					nocacheCfg{Regexp: "^/api$"},
					nocacheCfg{Regexp: "/d{1,2}*"},
				},
			},
		},
		{
			name: "case 2",
			args: args{
				rules: []rule.Nocacher{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.load(tt.args.rules)
			if want := len(tt.args.rules); c.cntRegexp != want {
				t.Errorf("could not initRules, not equal length: %d, want %d",
					c.cntRegexp, want)
			}
		})
	}
}

func Test_matchNoCacheRule(t *testing.T) {
	c := &Cache{}
	c.load([]rule.Nocacher{
		nocacheCfg{Regexp: "^/api/url$"},
		nocacheCfg{Regexp: "^/api/test$"},
		nocacheCfg{Regexp: "^/api/hire$"},
	})

	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case 1",
			args: args{
				uri: "/api/url",
			},
			want: true,
		},
		{
			name: "case 1",
			args: args{
				uri: "/api/hhhh/ashdak",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.matchNoCacheRule(tt.args.uri); got != tt.want {
				t.Errorf("matchNoCacheRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
