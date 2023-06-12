package unit

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	convey.Convey("start test new", t, func() {
		convey.So(42, convey.ShouldEqual, 42)
		fmt.Println(222222)
	})
}

func TestGetPersonDetail(t *testing.T) {
	convey.Convey("将两数相减", t, func() {
		fmt.Println(111)
	})

	type args struct {
		username string
	}

	tests := []struct {
		name    string
		args    args
		want    *PersonDetail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPersonDetail(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPersonDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPersonDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "email valid",
			args: args{
				email: "1234567@qq.com",
			},
			want: true,
		},
		{
			name: "email invalid",
			args: args{
				email: "test.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		got := checkEmail(tt.args.email)
		assert.Equal(t, tt.want, got)
	}
}

func TestToQiepIAN(t *testing.T) {
	var c map[string]interface{}
	c = make(map[string]interface{})
	c = map[string]interface{}{
		"img": []string{"xxx", "aaa"},
	}
	fmt.Println(c["img"].([]string))
}
