package presistence

import (
	"reflect"
	"testing"
)

func Test_Store(t *testing.T) {

}

func Test_PluginCode(t *testing.T) {
	if PlgCodeProxyPath != PluginCode(0x01) {
		t.Errorf("want: %v, got: %v", PluginCode(0x01), PlgCodeProxyPath)
	}
	if PlgCodeProxyServer != PluginCode(0x02) {
		t.Errorf("want: %v, got: %v", PluginCode(0x02), PlgCodeProxyServer)
	}
	if PlgCodeProxyReverseSrv != PluginCode(0x04) {
		t.Errorf("want: %v, got: %v", PluginCode(0x04), PlgCodeProxyReverseSrv)
	}
	if PlgCodeCache != PluginCode(0x08) {
		t.Errorf("want: %v, got: %v", PluginCode(0x08), PlgCodeCache)
	}
	if PlgCodeRatelimit != PluginCode(0x10) {
		t.Errorf("want: %v, got: %v", PluginCode(0x10), PlgCodeRatelimit)
	}
}

func TestListPlgByCode(t *testing.T) {
	type args struct {
		code PluginCode
	}
	tests := []struct {
		name string
		args args
		want []PluginCode
	}{
		{
			name: "case 1",
			args: args{
				code: PluginCode(0x0f),
			},
			want: []PluginCode{
				PlgCodeProxyPath,
				PlgCodeProxyServer,
				PlgCodeProxyReverseSrv,
				PlgCodeCache,
			},
		},
		{
			name: "case 2",
			args: args{
				code: PluginCode(0x11),
			},
			want: []PluginCode{
				PlgCodeProxyPath,
				PlgCodeRatelimit,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListPlgByCode(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListPlgByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
