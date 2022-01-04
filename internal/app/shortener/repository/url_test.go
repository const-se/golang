package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_repository_URL(t *testing.T) {
	type fields struct {
		storage []string
	}

	type args struct {
		id int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Получение ссылки",
			fields: fields{
				storage: []string{
					"https://ya.ru",
					"https://practicum.yandex.ru",
				},
			},
			args: args{
				id: 1,
			},
			want:    "https://practicum.yandex.ru",
			wantErr: assert.NoError,
		},
		{
			name: "Получение несуществующей ссылки",
			fields: fields{
				storage: []string{
					"https://ya.ru",
					"https://practicum.yandex.ru",
				},
			},
			args: args{
				id: 2,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				storage: tt.fields.storage,
			}

			got, err := r.URL(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("URL(%v)", tt.args.id)) {
				return
			}

			assert.Equalf(t, tt.want, got, "URL(%v)", tt.args.id)
		})
	}
}
