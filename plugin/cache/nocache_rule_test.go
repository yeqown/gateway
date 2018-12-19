package cache

import (
	"testing"

	"github.com/yeqown/gateway/config/presistence/filestore"
	"github.com/yeqown/gateway/config/rule"
)

func Test_initRules(t *testing.T) {
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
					&filestore.NocacheCfg{Regexp: "^/api$"},
					&filestore.NocacheCfg{Regexp: "/d{1,2}*"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initRules(tt.args.rules)
			if want := len(tt.args.rules); cntRegexp != want {
				t.Errorf("could not initRules, not equal length: %d, want %d",
					cntRegexp, want)
			}
		})
	}
}

func Test_matchNoCacheRule(t *testing.T) {
	initRules([]rule.Nocacher{
		&filestore.NocacheCfg{Regexp: "^/api/url$"},
		&filestore.NocacheCfg{Regexp: "^/api/test$"},
		&filestore.NocacheCfg{Regexp: "^/api/hire$"},
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
			if got := matchNoCacheRule(tt.args.uri); got != tt.want {
				t.Errorf("matchNoCacheRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
