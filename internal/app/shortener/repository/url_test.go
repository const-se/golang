package repository

import (
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
		wantErr bool
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
			want: "https://practicum.yandex.ru",
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				storage: tt.fields.storage,
			}

			got, err := r.URL(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("URL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
