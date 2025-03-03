package syringe

import (
	"context"
	"io"
)

func List(ctx context.Context, w io.Writer) error {
	w.Write([]byte{})

	return nil
}
