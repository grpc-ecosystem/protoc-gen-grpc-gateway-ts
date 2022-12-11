package data

import "testing"

func TestGetModuleName(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		fileName    string
		want        string
	}{
		{"empty", "", "", ""},
		{"simple", "mypackage", "service.proto", "MypackageService"},
		{"with file path", "mypackage", "path/to/proto/file/service.proto", "MypackageService"},
		{"with underscore", "my_package", "cool_service.proto", "MyPackageCoolService"},
		{"with dash", "my-package", "cool-service.proto", "MyPackageCoolService"},
		{"with dash and underscore", "my-package", "cool_service.proto", "MyPackageCoolService"},
		{"with dots", "my.package", "cool.service.proto", "MyPackageCoolService"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetModuleName(tt.packageName, tt.fileName); got != tt.want {
				t.Errorf("GetModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}
