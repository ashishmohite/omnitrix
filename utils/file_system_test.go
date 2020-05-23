package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Should return true when file is present at given path",
			args:    struct{ path string }{path: "./file_system.go"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Should return false when file is not present at given path",
			args:    struct{ path string }{path: "../file_system.go"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Should return error when path is not a file",
			args:    struct{ path string }{path: ".."},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Should return true when directory is present at given path",
			args:    struct{ path string }{path: "."},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Should return false when file is not present at given path",
			args:    struct{ path string }{path: "../random_dir"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Should return error when path is not a directory",
			args:    struct{ path string }{path: "./file_system.go"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DirectoryExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirectoryExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DirectoryExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryCleanup(t *testing.T) {
	type args struct {
		path      string
		fileMode  os.FileMode
		createDir bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should delete all the files from given directory",
			args:    args{path: "dir/dir_1", fileMode: os.FileMode(0777), createDir: true},
			wantErr: false,
		},
		{
			name:    "Should fail if error while deleting directory content",
			args:    args{path: "dir/dir_1", fileMode: os.FileMode(0777), createDir: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.createDir {
				_ = os.MkdirAll(filepath.Join(tt.args.path, "test"), os.FileMode(tt.args.fileMode))
			}
			if err := DirectoryCleanup(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("DirectoryCleanup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.createDir {
				path := filepath.Join(tt.args.path, "test")
				if _, err := IsDirectoryEmpty(path); err == nil {
					t.Errorf("DirectoryCleanup() error = path %s exists", path)
				}

				_ = os.RemoveAll(tt.args.path)
			}
		})
	}
}
