package unit

import (
	"api_client/model"
	v1 "api_client/request"
	"api_client/service"
	"reflect"
	"testing"
)

func TestFoundPersonService_GetPersonList(t *testing.T) {
	type args struct {
		info *v1.PersonListParams
	}
	var x = v1.PersonListParams{
		CrPerson: model.CrPerson{Name: "ajiu"},
		Page:     1,
		PageSize: 10,
	}
	tests := []struct {
		name      string
		args      args
		wantList  interface{}
		wantTotal int64
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:      "axin",
			args:      args{info: &x},
			wantList:  x,
			wantTotal: 60,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &service.FoundPersonService{}
			gotList, gotTotal, err := c.GetPersonList(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("GetPersonList() gotList = %v, want %v", gotList, tt.wantList)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("GetPersonList() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
