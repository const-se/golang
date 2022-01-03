package repository

import (
	"testing"
)

func Test_repository_SaveURL(t *testing.T) {
	type fields struct {
		storage []string
	}

	type args struct {
		url string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Добавление ссылки",
			fields: fields{
				storage: []string{
					"https://ya.ru",
				},
			},
			args: args{
				url: "https://practicum.yandex.ru",
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				storage: tt.fields.storage,
			}

			got, err := r.SaveURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("SaveURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
