package realip

import (
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{r: &http.Request{}},
			want: "",
		},
		{
			name: "ok_real_ipv4_with_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"10.255.30.50:8878"}}}},
			want: "10.255.30.50:8878",
		},
		{
			name: "ok_real_ipv4_without_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"10.255.30.50"}}}},
			want: "10.255.30.50",
		},
		{
			name: "ok_real_ipv6_with_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"[2001:db8:3333:4444:5555:6666:7777:8888]:8878"}}}},
			want: "[2001:db8:3333:4444:5555:6666:7777:8888]:8878",
		},
		{
			name: "ok_real_ipv6_without_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"2001:db8:3333:4444:5555:6666:7777:8888"}}}},
			want: "2001:db8:3333:4444:5555:6666:7777:8888",
		},
		{
			name: "error_real_ipv4",
			args: args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"10.255.30.500"}}}},
			want: "",
		},
		{
			name: "error_real_text",
			args: args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"text"}}}},
			want: "",
		},
		{
			name: "ok_forwarded_ipv4_with_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"151.255.30.50:8878", "27.10.10.11:33254"}}}},
			want: "151.255.30.50:8878",
		},
		{
			name: "ok_forwarded_ipv4_without_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"151.255.30.50", "27.10.10.11"}}}},
			want: "151.255.30.50",
		},
		{
			name: "ok_forwarded_ipv4_without_port_with_private_ip",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"10.255.30.50, 151.255.30.50, 27.10.10.11"}}}},
			want: "151.255.30.50",
		},
		{
			name: "ok_forwarded_ipv6_with_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"[2001:db8:3333:4444:5555:6666:7777:8888]:8878", "[2001:db8:1234:ffff:ffff:ffff:ffff:ffff]:5521"}}}},
			want: "[2001:db8:3333:4444:5555:6666:7777:8888]:8878",
		},
		{
			name: "ok_forwarded_ipv6_without_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"2001:db8:3333:4444:5555:6666:7777:8888", "[2001:db8:1234:ffff:ffff:ffff:ffff:ffff]:5521"}}}},
			want: "2001:db8:3333:4444:5555:6666:7777:8888",
		},
		{
			name: "error_forwarded_ipv4",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"10.255.30.500"}}}},
			want: "",
		},
		{
			name: "error_forwarded_ipv6",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"2001:ffdb8:3333:4444:5555:6666:7777:8888"}}}},
			want: "",
		},
		{
			name: "error_forwarded_text",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"text"}}}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.args.r)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
