package crypto

import (
	"testing"
)

func TestGenerateRSA(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Returns error on failing to generate keys",
			args: args{
				size: 8,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateRSA(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRSA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
