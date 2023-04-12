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
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{r: &http.Request{}},
			wantErr: true,
		},
		{
			name: "ok_remote_addr_ipv4_with_port",
			args: args{r: &http.Request{RemoteAddr: "10.255.30.50:8878"}},
			want: "10.255.30.50:8878",
		},
		{
			name: "ok_remote_addr_ipv4_without_port",
			args: args{r: &http.Request{RemoteAddr: "10.255.30.50"}},
			want: "10.255.30.50",
		},
		{
			name:    "error_remote_addr_ipv4",
			args:    args{r: &http.Request{RemoteAddr: "10.255.30.500"}},
			wantErr: true,
		},
		{
			name:    "error_remote_addr_text",
			args:    args{r: &http.Request{RemoteAddr: "text"}},
			wantErr: true,
		},
		{
			name: "ok_remote_addr_ipv6_with_port",
			args: args{r: &http.Request{RemoteAddr: "[2001:db8:3333:4444:5555:6666:7777:8888]:8878"}},
			want: "[2001:db8:3333:4444:5555:6666:7777:8888]:8878",
		},
		{
			name: "ok_remote_addr_ipv6_without_port",
			args: args{r: &http.Request{RemoteAddr: "2001:db8:3333:4444:5555:6666:7777:8888"}},
			want: "2001:db8:3333:4444:5555:6666:7777:8888",
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
			name:    "error_real_ipv4",
			args:    args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"10.255.30.500"}}}},
			wantErr: true,
		},
		{
			name:    "error_real_text",
			args:    args{r: &http.Request{Header: map[string][]string{"X-Real-Ip": {"text"}}}},
			wantErr: true,
		},

		{
			name: "ok_forwarded_ipv4_with_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"10.255.30.50:8878", "10.10.10.11:33254"}}}},
			want: "10.255.30.50:8878",
		},
		{
			name: "ok_forwarded_ipv4_without_port",
			args: args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"10.255.30.50", "10.10.10.11"}}}},
			want: "10.255.30.50",
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
			name:    "error_forwarded_ipv4",
			args:    args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"10.255.30.500"}}}},
			wantErr: true,
		},
		{
			name:    "error_forwarded_ipv6",
			args:    args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"2001:ffdb8:3333:4444:5555:6666:7777:8888"}}}},
			wantErr: true,
		},
		{
			name:    "error_forwarded_text",
			args:    args{r: &http.Request{Header: map[string][]string{"X-Forwarded-For": {"text"}}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
