package urlshortener

import "testing"

func TestMakeUrlShort(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeUrlShort(); got != tt.want {
				t.Errorf("MakeUrlShort() = %v, want %v", got, tt.want)
			}
		})
	}
}
