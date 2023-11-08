package service

import (
	"context"
	"errors"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/functionsignature"
	"github.com/BlockPILabs/aaexplorer/internal/memo"
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"github.com/BlockPILabs/aaexplorer/service"
	"time"
)

type functionSignatureService struct {
}

var FunctionSignatureService = &functionSignatureService{}

func (s *functionSignatureService) GetMethodBySignature(ctx context.Context, tx *ent.Client, signature string) (f *ent.FunctionSignature, err error) {
	if !utils.IsHexSting(signature) && len(signature) < 8 {
		return nil, errors.New("not signature")
	}

	f, err = tx.FunctionSignature.Get(ctx, signature)
	if err != nil {
		// db not found , www.4byte.directory
		key := "f:" + signature
		ferr, ok := memo.Get(key)
		if ok {
			err, ok := ferr.(error)
			if ok {
				return nil, err
			}
		}
		functionSignature, err := service.GetSignature(signature)
		if err != nil {
			memo.SetWithTTL(key, err, 1, time.Hour)
			return nil, err
		}

		f = &ent.FunctionSignature{
			ID:         functionSignature.HexSignature,
			Name:       functionSignature.Name,
			Text:       functionSignature.TextSignature,
			Bytes:      functionSignature.BytesSignature,
			CreateTime: time.Now(),
		}

		tx.FunctionSignature.Create().
			SetID(f.ID).
			SetName(f.Name).
			SetText(f.Text).
			SetBytes(f.Bytes).
			SetCreateTime(f.CreateTime).
			OnConflictColumns(functionsignature.FieldID).
			Update(func(upsert *ent.FunctionSignatureUpsert) {
				upsert.
					UpdateName().
					UpdateBytes().
					UpdateText().
					UpdateText()
			}).
			ExecX(ctx)

	}
	return
}
