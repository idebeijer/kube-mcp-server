package resource

import "testing"

func TestExtractNamespaceFromURI(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		args    args
		want    string
		wantErr bool
	}{
		{args{uri: "k8s://default/pods"}, "default", false},
		{args{uri: "k8s://test/services"}, "test", false},
		{args{uri: "k8s://default"}, "default", false},
		{args{uri: "k8s://"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.args.uri, func(t *testing.T) {
			got, err := ExtractNamespaceFromURI(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractNamespaceFromURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractNamespaceFromURI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
