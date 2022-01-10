package repository

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_repository_URL(t *testing.T) {
	type fields struct {
		cache []string
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
				cache: []string{
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
				cache: []string{
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
				cache: tt.fields.cache,
			}

			got, err := r.URL(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("URL(%v)", tt.args.id)) {
				return
			}

			assert.Equalf(t, tt.want, got, "URL(%v)", tt.args.id)
		})
	}
}

func Test_repository_SaveURL(t *testing.T) {
	type fields struct {
		cache []string
	}

	type args struct {
		url string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Добавление ссылки",
			fields: fields{
				cache: []string{
					"https://ya.ru",
				},
			},
			args: args{
				url: "https://practicum.yandex.ru",
			},
			want:    1,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				cache: tt.fields.cache,
			}

			got, err := r.SaveURL(tt.args.url)
			if !tt.wantErr(t, err, fmt.Sprintf("SaveURL(%v)", tt.args.url)) {
				return
			}

			assert.Equalf(t, tt.want, got, "SaveURL(%v)", tt.args.url)
		})
	}
}
