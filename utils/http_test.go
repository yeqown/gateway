package utils

import "testing"

func TestParseURIPrefix(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				uri: "/prefix/uri/hhh/",
			},
			want: "/prefix",
		},
		{
			name: "case 2",
			args: args{
				uri: "prefix/uri",
			},
			want: "/prefix",
		},
		{
			name: "case 3",
			args: args{
				uri: "/prefix",
			},
			want: "/prefix",
		},
		{
			name: "case 4",
			args: args{
				uri: "prefix",
			},
			want: "/nonePrefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseURIPrefix(tt.args.uri); got != tt.want {
				t.Errorf("ParseURIPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
