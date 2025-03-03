package cmd

import (
	"io"
)

func List(w io.Writer) error {
	w.Write([]byte{})

	return nil
}
