package main

import (
	"testing"
)

//func TestNewStorage(t *testing.T) {
//	type args struct {
//		path    string
//		sources Sources
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Storage
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewStorage(tt.args.path, tt.args.sources); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStorage_prepareDirectory(t *testing.T) {
//	type fields struct {
//		MainPath   string
//		TimePath   string
//		SourcePath string
//		Prepared   bool
//		Accessible bool
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Storage{
//				MainPath:   tt.fields.MainPath,
//				TimePath:   tt.fields.TimePath,
//				SourcePath: tt.fields.SourcePath,
//				Prepared:   tt.fields.Prepared,
//				Accessible: tt.fields.Accessible,
//			}
//			if err := s.prepareDirectory(); (err != nil) != tt.wantErr {
//				t.Errorf("prepareDirectory() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestStorage_saveAllData(t *testing.T) {
//	type fields struct {
//		MainPath   string
//		TimePath   string
//		SourcePath string
//		Prepared   bool
//		Accessible bool
//	}
//	type args struct {
//		saveAllData []SaveAllData
//		fileName    string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Storage{
//				MainPath:   tt.fields.MainPath,
//				TimePath:   tt.fields.TimePath,
//				SourcePath: tt.fields.SourcePath,
//				Prepared:   tt.fields.Prepared,
//				Accessible: tt.fields.Accessible,
//			}
//		})
//	}
//}
//
//func TestStorage_saveJson(t *testing.T) {
//	type fields struct {
//		MainPath   string
//		TimePath   string
//		SourcePath string
//		Prepared   bool
//		Accessible bool
//	}
//	type args struct {
//		data         []byte
//		fullFileName string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Storage{
//				MainPath:   tt.fields.MainPath,
//				TimePath:   tt.fields.TimePath,
//				SourcePath: tt.fields.SourcePath,
//				Prepared:   tt.fields.Prepared,
//				Accessible: tt.fields.Accessible,
//			}
//		})
//	}
//}
//
//func TestStorage_saveOrganized(t *testing.T) {
//	type fields struct {
//		MainPath   string
//		TimePath   string
//		SourcePath string
//		Prepared   bool
//		Accessible bool
//	}
//	type args struct {
//		services OrganizedServices
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Storage{
//				MainPath:   tt.fields.MainPath,
//				TimePath:   tt.fields.TimePath,
//				SourcePath: tt.fields.SourcePath,
//				Prepared:   tt.fields.Prepared,
//				Accessible: tt.fields.Accessible,
//			}
//		})
//	}
//}

func Test_cleanFilePathName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"success", "Ab GtHnJzt._456", "ab gthnjzt._456"},
		{"spaces", "  AbGtHnJzt._456  ", "abgthnjzt._456"},
		{"dot end", "  AbGtHnJzt._456..", "abgthnjzt._456"},
		{"invalid chars", "abcdef/ge*h+ij\\kl|m", "abcdefgehijklm"},
		{"invalid chars", "abcdef/ge*h+ij\\kl|m", "abcdefgehijklm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanFilePathName(tt.path); got != tt.want {
				t.Errorf("cleanFilePathName() = %v, want %v", got, tt.want)
			}
		})
	}
}
