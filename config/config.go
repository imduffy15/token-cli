package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/imduffy15/token-cli/client"
)

func Dir() string {
	return path.Join(userHomeDir(), ".token-cli")
}

func Path() string {
	return path.Join(Dir(), "config.json")
}

func Read() client.Config {
	c := client.NewConfig()

	data, err := ioutil.ReadFile(Path())
	if err != nil {
		fmt.Print(fmt.Errorf(err.Error()))
		return c
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		fmt.Print(fmt.Errorf(err.Error()))
		return c
	}

	return c
}

func Write(c client.Config) error {
	err := makeDirectory()
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	path := Path()
	return ioutil.WriteFile(path, data, 0600)
}
