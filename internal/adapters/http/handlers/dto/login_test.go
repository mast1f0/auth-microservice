package dto

import "testing"

func TestLoginReqValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     LoginReq
		wantErr bool
	}{
		{name: "valid", req: LoginReq{Login: "user", Password: "password123"}},
		{name: "missing login", req: LoginReq{Password: "password123"}, wantErr: true},
		{name: "missing password", req: LoginReq{Login: "user"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}
