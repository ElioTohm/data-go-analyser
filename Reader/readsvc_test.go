package Reader

import (
	"strings"
	"testing"
)

func TestAnalyseData(t *testing.T) {
	type args struct {
		customers []*Customer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    t.Name(),
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		dummyfile := strings.NewReader("id,first_name,last_name,customer_reference,status\n1,Alex,Richard,ade3-11ed,Active\nid,customer_reference,order_status,order_reference,order_timestamp\n1,ade3-11ed,Delivered,dc0aa69c,1676539508\n2,ade3-11ed,Delivered,dc0ab1aa,1676539608\nid,order_reference,item_name,quantity,total_price\n1,dc0aa69c,XYZ,2,20\n2,dc0aa69c,ABC,1,55\n3,dc0ab1aa,BBB,1,12.5")
		constructedData, _ := ConstructData(ReadFile(dummyfile))
		t.Run(tt.name, func(t *testing.T) {
			if err := AnalyseData(constructedData); (err != nil) != tt.wantErr {
				t.Errorf("AnalyseData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
