package utils

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

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

func TestEncodeFormToString(t *testing.T) {
	getReq, err := http.NewRequest("GET",
		"http://localhost:8080/api/url?q=params", nil)
	if err != nil {
		t.Errorf("could not generate a new request: %v", err)
	}
	formEncoded := EncodeFormToString(getReq)
	t.Logf("url %s, encoded: %s",
		"http://localhost:8080/api/url?q=params", formEncoded)

	// Test encde post form
	body := url.Values{}
	body.Add("q", "value")
	body.Add("q1", "value1")

	postReq, err := http.NewRequest("POST",
		"http://localhost:8080/api/url", strings.NewReader(body.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Errorf("could not generate a new request: %v", err)
	}
	formEncoded = EncodeFormToString(postReq)
	t.Logf("url %s, encoded: %s", "http://localhost:8080/api/url?q=params", formEncoded)
}

func BenchmarkEncodeFormToString(b *testing.B) {
	b.StopTimer()
	getReq, err := http.NewRequest("GET",
		"http://localhost:8080/api/url?q=params", nil)
	if err != nil {
		b.Errorf("could not generate a new request: %v", err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		formEncoded := EncodeFormToString(getReq)
		_ = formEncoded
	}
}

func TestCopyRequest(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CopyRequest(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CopyRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRequestForm(t *testing.T) {
	// generate request var
	getReq, _ := http.NewRequest("GET",
		"http://localhost:8080/api/url?q=value", nil)

	body := url.Values{}
	body.Add("q", "value")
	body.Add("q1", "value1")
	postReq, _ := http.NewRequest("POST",
		"http://localhost:8080/api/url", strings.NewReader(body.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// test case ...
	type args struct {
		cpyReq *http.Request
	}
	tests := []struct {
		name string
		args args
		want url.Values
	}{
		{
			name: "case-1-GET",
			args: args{
				cpyReq: getReq,
			},
			want: url.Values{
				"q": []string{"value"},
			},
		},
		{
			name: "case-2-POST",
			args: args{
				cpyReq: postReq,
			},
			want: url.Values{
				"q":  []string{"value"},
				"q1": []string{"value1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRequestForm(tt.args.cpyReq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRequestForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkParseRequestForm(b *testing.B) {
	b.StopTimer()
	getReq, _ := http.NewRequest("GET",
		"http://localhost:8080/api/url?q=value&q1=value1&q2=value2", nil)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		form := ParseRequestForm(getReq)
		_ = form
	}
}
