package cmd

import (
	"fmt"
	"io"
)

func Set(key string, value string, out io.Writer) error {
	out.Write([]byte(fmt.Sprintf("%s %s", key, value)))

	return nil
}
