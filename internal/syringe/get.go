package syringe

import (
	"context"
	"io"
)

func Get(ctx context.Context, w io.Writer, key string) error {
	w.Write([]byte{})

	return nil
}
