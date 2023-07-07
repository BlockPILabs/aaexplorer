package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type bundleService struct {
}

var BundleService = bundleService{}

func (*bundleService) GetBundlers(ctx context.Context, response vo.GetNetworksResponse) {

}
