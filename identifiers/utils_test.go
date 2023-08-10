package identifiers

import "testing"

func TestCalcHashFNV(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestCalcHashFNV",
			args: args{
				id: "123",
			},
			want: "5003431119771845851",
		},
		{
			name: "TestCalcHashFNV-1",
			args: args{
				id: "1234",
			},
			want: "2282126479029740061",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcHashFNV(tt.args.id); got != tt.want {
				t.Errorf("CalcHashFNV() = %v, want %v", got, tt.want)
			}
		})
	}
}
