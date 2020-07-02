package main

import (
	"github.com/prometheus/common/model"
	"testing"
)

func TestPoint_UnmarshalJSON(t *testing.T) {
	type fields struct {
		T model.Time
		V float64
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Empty", fields{}, args{[]byte("")}, true},
		{"Bracket", fields{}, args{[]byte("[]")}, true},
		{"No number", fields{}, args{[]byte("[abc,\"abc\"]")}, true},
		{"Mixed number", fields{}, args{[]byte("[abc12 , \"a12bc\"]")}, true},
		{"Success format", fields{}, args{[]byte("[1593561600,\"10\"]")}, false},
		{"Format with space", fields{}, args{[]byte("[1593561600, \"10\"]")}, false},
		{"Format many spaces", fields{}, args{[]byte("  [  1593561600  , \"10\"   ]")}, false},
		{"Float T", fields{}, args{[]byte("  [  1593561600.360, \"10\"   ]")}, false},
		{"Float V", fields{}, args{[]byte("  [  1593561600, \"10.1\"   ]")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &point{
				T: tt.fields.T,
				V: tt.fields.V,
			}
			if err := p.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
