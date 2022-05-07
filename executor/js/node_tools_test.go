package js

import (
	"reflect"
	"testing"
)

func Test_formatArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No args",
			args: struct{ args []string }{args: []string{}},
			want: "",
		},
		{
			name: "Single arg",
			args: struct{ args []string }{args: []string{"firstArg"}},
			want: "'firstArg'",
		},
		{
			name: "Two arg",
			args: struct{ args []string }{args: []string{"firstArg", "secondArg"}},
			want: "'firstArg','secondArg'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatArgs(tt.args.args); got != tt.want {
				t.Errorf("formatArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepareNodeCommand(t *testing.T) {
	type args struct {
		path    string
		command string
		args    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No ags command",
			args: struct {
				path    string
				command string
				args    []string
			}{path: "/test/command.js", command: "exec", args: []string{}},
			want: "console.log(require('/test/command.js').exec())",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareNodeCommand(tt.args.path, tt.args.command, tt.args.args); got != tt.want {
				t.Errorf("prepareNodeCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseNodeArrayResult(t *testing.T) {
	type args struct {
		result string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Empty array",
			args:    struct{ result string }{result: "[]"},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "Single string value array",
			args:    struct{ result string }{result: "['firstArg']"},
			want:    []string{"firstArg"},
			wantErr: false,
		},
		{
			name:    "Two string value array",
			args:    struct{ result string }{result: "['firstArg',   'secondArg'  ]"},
			want:    []string{"firstArg", "secondArg"},
			wantErr: false,
		},
		{
			name:    "Single number value array",
			args:    struct{ result string }{result: "[ 12.3 ]"},
			want:    []string{"12.3"},
			wantErr: false,
		},
		{
			name:    "Two number value array",
			args:    struct{ result string }{result: "[32.0,   11]"},
			want:    []string{"32.0", "11"},
			wantErr: false,
		},
		{
			name:    "Two mixed value array",
			args:    struct{ result string }{result: "['test',   11]"},
			want:    []string{"test", "11"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNodeArrayResult(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNodeArrayResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNodeArrayResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}
