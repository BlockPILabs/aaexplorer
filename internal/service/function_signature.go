package service

import (
	"context"
	"errors"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/memo"
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"time"
)

type functionSignatureService struct {
}

var FunctionSignatureService = &functionSignatureService{}

func (s *functionSignatureService) GetMethodBySignature(ctx context.Context, tx *ent.Client, signature string) (f *ent.FunctionSignature, err error) {
	if !utils.IsHexSting(signature) && len(signature) < 8 {
		return nil, errors.New("not signature")
	}
	key := "f:" + signature
	mf, ok := memo.Get(key)
	if ok {
		switch t := mf.(type) {
		case error:
			return nil, t
		case *ent.FunctionSignature:
			return t, nil
		}
	}

	f, err = tx.FunctionSignature.Get(ctx, signature)
	if err != nil {

		memo.SetWithTTL(key, err, 1, time.Hour)
		return nil, err

	} else {
		memo.Set(key, f, 2)
	}
	return
}
