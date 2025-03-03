package cmd

import (
	"io"
)

func Get(key string, out io.Writer) error {
	out.Write([]byte{})

	return nil
}
