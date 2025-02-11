package handlers

import (
	"net/http"
	"reflect"
	"testing"
)

func TestAddUrl(t *testing.T) {
	type args struct {
		db *storage.Store
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddUrl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUrl(t *testing.T) {
	type args struct {
		db *storage.Store
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUrl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedirectUrl(t *testing.T) {
	type args struct {
		db *storage.Store
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedirectUrl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedirectUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
