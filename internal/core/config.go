package core

import (
	"context"

	"github.com/tremj/lbx/internal/parser"
	"github.com/tremj/lbx/internal/storage"
)

func SaveConfig(ctx context.Context, name string, data []byte) error {
	err := parser.ValidateConfig(data)
	if err != nil {
		return err
	}
	return storage.SaveConfig(ctx, name, data)
}

func DeleteConfig(ctx context.Context, name string) error {
	return storage.DeleteConfig(ctx, name)
}
