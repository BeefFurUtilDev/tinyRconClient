package main

import (
	"os"
	"testing"
)

func Test_execCommand(t *testing.T) {
	cmdTest1 := "spark tps"
	cmdTest2 := "list"
	type args struct {
		clientSetup *client
		cmd         *string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				clientSetup: &client{
					addr:     os.Getenv("IP"),
					port:     25575,
					password: os.Getenv("PASSWORD"),
				},
				cmd: &cmdTest1,
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "test2",
			args: args{
				clientSetup: &client{
					addr:     os.Getenv("IP"),
					port:     25575,
					password: os.Getenv("PASSWORD"),
				},
				cmd: &cmdTest2,
			},
			wantErr:    false,
			wantResult: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := execCommand(tt.args.clientSetup, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("execCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("execCommand() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
