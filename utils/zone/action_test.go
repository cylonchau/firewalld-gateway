package zone

import (
	"github.com/godbus/dbus/v5"
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	tests := []struct {
		name string
		want *Rish
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInstance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRish_AddRich(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		rule Rich
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if err := AddRich(tt.args.rule); (err != nil) != tt.wantErr {
				t.Errorf("AddRich() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRish_DelRich(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		rule Rich
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if err := DelRich(tt.args.rule); (err != nil) != tt.wantErr {
				t.Errorf("DelRich() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRish_QueryAllRish(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	tests := []struct {
		name      string
		fields    fields
		wantRichs []string
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			gotRichs, err := QueryAllRish()
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryAllRish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRichs, tt.wantRichs) {
				t.Errorf("QueryAllRish() gotRichs = %v, want %v", gotRichs, tt.wantRichs)
			}
		})
	}
}

func TestRish_QueryRich(t *testing.T) {

	var (
		queryParam = QueryRich{
			Address:  "10.0.0.19",
			Port:     8088,
			Protocol: "tcp",
			Type:     "",
			Expire:   0,
			Zone:     "",
		}
	)

	type args struct {
		rule QueryRich
	}
	tests := []struct {
		name         string
		args         args
		wantIsExists bool
	}{
		{
			name: "query one from firewalld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if gotIsExists := QueryRich(tt.args.rule); gotIsExists != tt.wantIsExists {
				t.Errorf("QueryRich() = %v, want %v", gotIsExists, tt.wantIsExists)
			}
		})
	}
}

func TestRish_listRichParse(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		lists []string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantJsonList []string
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			gotJsonList, err := listRichParse(tt.args.lists)
			if (err != nil) != tt.wantErr {
				t.Errorf("listRichParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJsonList, tt.wantJsonList) {
				t.Errorf("listRichParse() gotJsonList = %v, want %v", gotJsonList, tt.wantJsonList)
			}
		})
	}
}

func TestRish_oneRichParse(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		richs []string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantJsonrich []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if gotJsonrich := oneRichParse(tt.args.richs); !reflect.DeepEqual(gotJsonrich, tt.wantJsonrich) {
				t.Errorf("oneRichParse() = %v, want %v", gotJsonrich, tt.wantJsonrich)
			}
		})
	}
}

func TestRish_richExists(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		rule QueryRich
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if got := richExists(tt.args.rule); got != tt.want {
				t.Errorf("richExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRish_structToQueryRich(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		r QueryRich
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if got := structToQueryRich(tt.args.r); got != tt.want {
				t.Errorf("structToQueryRich() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRish_structToRich(t *testing.T) {
	type fields struct {
		Conn *dbus.Conn
	}
	type args struct {
		r Rich
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &Rish{
				Conn: tt.fields.Conn,
			}
			if got := structToRich(tt.args.r); got != tt.want {
				t.Errorf("structToRich() = %v, want %v", got, tt.want)
			}
		})
	}
}
