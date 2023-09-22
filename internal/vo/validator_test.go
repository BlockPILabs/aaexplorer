package vo

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"GetUserOpsRequest", args{s: GetUserOpsRequest{
				PaginationRequest: NewDefaultPaginationRequest(),
				Network:           "polygon",
				LatestBlockNumber: 0,
				BlockNumber:       0,
				StartTime:         0,
				EndTime:           0,
				TxHash:            "",
				Bundler:           "",
				Paymaster:         "",
				Factory:           "",
				Account:           "",
				HashTerm:          "",
			}}, false,
		},
		{
			"GetUserOpsRequest", args{s: GetUserOpsRequest{
				PaginationRequest: NewDefaultPaginationRequest(),
				Network:           "polygon",
				LatestBlockNumber: 0,
				BlockNumber:       0,
				StartTime:         0,
				EndTime:           0,
				TxHash:            "",
				Bundler:           "0x12",
				Paymaster:         "",
				Factory:           "",
				Account:           "",
				HashTerm:          "",
			}}, true,
		},
		{
			"GetUserOpsRequest", args{s: GetUserOpsRequest{
				PaginationRequest: NewDefaultPaginationRequest(),
				Network:           "polygon",
				LatestBlockNumber: 0,
				BlockNumber:       0,
				StartTime:         0,
				EndTime:           0,
				TxHash:            "",
				Bundler:           "0xfdc58263014de1dfe5a08db053ffe67d9faa958b",
				Paymaster:         "",
				Factory:           "",
				Account:           "",
				HashTerm:          "",
			}}, false,
		},
		{
			"GetUserOpsRequest", args{s: GetUserOpsRequest{
				PaginationRequest: NewDefaultPaginationRequest(),
				Network:           "polygon",
				LatestBlockNumber: 0,
				BlockNumber:       0,
				StartTime:         0,
				EndTime:           0,
				TxHash:            "0xfdc58263014de1dfe5a08db053ffe67d9faa958b",
				Bundler:           "",
				Paymaster:         "",
				Factory:           "",
				Account:           "",
				HashTerm:          "",
			}}, true,
		},
		{
			"GetUserOpsRequest", args{s: GetUserOpsRequest{
				PaginationRequest: NewDefaultPaginationRequest(),
				Network:           "polygon",
				LatestBlockNumber: 0,
				BlockNumber:       0,
				StartTime:         0,
				EndTime:           0,
				TxHash:            "0xe9ba942c8c332216b7db9f9511ed61b01f9bc428bc97809f7a40d9e4b677c223",
				Bundler:           "",
				Paymaster:         "",
				Factory:           "",
				Account:           "",
				HashTerm:          "",
			}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil {
				t.Logf("errors = %v", err)
			}
		})
	}
}
