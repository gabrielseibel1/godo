package commands

import (
	"testing"

	"github.com/gabrielseibel1/godo/data"
	"github.com/gabrielseibel1/godo/types"
)

func TestDo_String(t *testing.T) {
	type fields struct {
		id   types.ID
		repo data.Repository
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no id",
			fields: fields{
				id:   "",
				repo: nil,
			},
			want: "command do",
		},
		{
			name: "some id",
			fields: fields{
				id:   "bofa",
				repo: nil,
			},
			want: "command do bofa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Do{
				id:   tt.fields.id,
				repo: tt.fields.repo,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("Do.String() = %v, want %v", tt.want, got)
			}
		})
	}
}
