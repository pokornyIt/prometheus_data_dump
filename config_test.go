package main

import "testing"

func TestConfig_overWriteFromLine(t *testing.T) {
	type args struct {
		server string
		path   string
	}
	type fields struct {
		Server string
		Path   string
		Days   int
		Jobs   []string
		Step   int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect args
	}{
		{"none", fields{"server.local", "/my", 1, nil, 10}, args{}, args{"server.local", "/my"}},
		{"server", fields{"server.local", "/my", 1, nil, 10}, args{"server.new", ""}, args{"server.new", "/my"}},
		{"path", fields{"server.local", "/my", 1, nil, 10}, args{"", "/newPath"}, args{"server.local", "/newPath"}},
		{"both", fields{"server.local", "/my", 1, nil, 10}, args{"server.new", "/newPath"}, args{"server.new", "/newPath"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configuration{
				Server: tt.fields.Server,
				Path:   tt.fields.Path,
				Days:   tt.fields.Days,
				Jobs:   tt.fields.Jobs,
				Step:   tt.fields.Step,
			}
			server = &tt.args.server
			directoryData = &tt.args.path
			c.overWriteFromLine()
			if c.Server != tt.expect.server {
				t.Errorf("overWriteFromLine() = server %s, want %s", c.Server, tt.expect.server)
			}
			if c.Path != tt.expect.path {
				t.Errorf("overWriteFromLine() = path %s, want %s", c.Path, tt.expect.path)
			}
		})
	}
}

func TestConfig_validate(t *testing.T) {
	type fields struct {
		Server string
		Path   string
		Days   int
		Jobs   []string
		Step   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"all valid", fields{"server.local", "./", 2, nil, 10}, false},
		{"no server", fields{"", "./", 2, nil, 10}, true},
		{"no path", fields{"", "", 2, nil, 10}, true},
		{"wrong day", fields{"", "./", -2, nil, 10}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &configuration{
				Server: tt.fields.Server,
				Path:   tt.fields.Path,
				Days:   tt.fields.Days,
				Jobs:   tt.fields.Jobs,
				Step:   tt.fields.Step,
			}
			if err := c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
