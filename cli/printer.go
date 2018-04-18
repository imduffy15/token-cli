package cli

import (
	"encoding/json"
)

type Printer interface {
	Print(interface{}) error
}

type JSONPrinter struct {
	Log Logger
}

func NewJSONPrinter(log Logger) JSONPrinter {
	return JSONPrinter{log}
}
func (jp JSONPrinter) Print(obj interface{}) error {
	j, err := json.MarshalIndent(&obj, "", "  ")
	if err != nil {
		return err
	}
	jp.Log.Robots(string(j))
	return nil
}
