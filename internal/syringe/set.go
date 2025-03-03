package syringe

import (
	"context"
	"fmt"
	"io"
)

func Set(ctx context.Context, w io.Writer, key string, value string) error {
	w.Write([]byte(fmt.Sprintf("%s %s", key, value)))

	return nil
}
