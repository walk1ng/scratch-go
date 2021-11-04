package utils

import (
	"encoding/json"
	"fmt"
	"io"
)

func PrintJsonPretty(out io.Writer, v interface{}) {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Printf("error ocurred in PrintJsonPretty: %s\n", err.Error())
		return
	}

	fmt.Fprint(out, string(data))
}
