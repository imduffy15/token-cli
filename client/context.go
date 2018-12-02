package client

import (
	"encoding/json"

	"github.com/zalando/go-keyring"
)

type Config struct {
	Verbose          bool
	Targets          map[string]*Target
	ActiveTargetName string
}

type Target struct {
	AuthorizationEndpoint   string
	TokenEndpoint           string
	Name                    string
	SkipSSLValidation       bool
	ActiveClientContextName string
	ClientContexts          map[string]struct{}
}

type ClientContext struct {
	ClientID string `json:"client_id"`
	Token
}

func NewConfig() Config {
	c := Config{}
	c.Targets = map[string]*Target{}
	return c
}

func (c *Config) AddTarget(newTarget Target) {
	c.Targets[newTarget.Name] = &newTarget
}

func (c *Config) GetContext(contextName string) (ClientContext, error) {
	strToken, err := keyring.Get(c.ActiveTargetName, contextName)
	if err != nil {
		return ClientContext{}, err
	}

	var token Token
	err = json.Unmarshal([]byte(strToken), &token)

	if err != nil {
		return ClientContext{}, err
	}

	return ClientContext{
		ClientID: contextName,
		Token:    token,
	}, nil
}

func (c *Config) AddContext(newClientContext ClientContext) error {
	return c.Targets[c.ActiveTargetName].AddClientContext(newClientContext)
}

func (c Config) GetTarget(targetName string) *Target {
	return c.Targets[targetName]
}

func (c Config) TargetExists(targetName string) bool {
	_, found := c.Targets[targetName]
	return found
}

func (c Config) DeleteTarget(targetName string) error {
	for k := range c.Targets[targetName].ClientContexts {
		err := keyring.Delete(targetName, k)
		if err != nil {
			return err
		}
		delete(c.Targets[targetName].ClientContexts, k)
	}
	delete(c.Targets, targetName)
	return nil
}

func (c Config) ListTargets() []string {
	keys := make([]string, len(c.Targets))

	i := 0
	for k := range c.Targets {
		keys[i] = k
		i++
	}

	return keys
}

func (c Config) GetActiveTarget() *Target {
	return c.Targets[c.ActiveTargetName]
}

func (t Target) ClientContextExists(clientName string) bool {
	_, found := t.ClientContexts[clientName]
	return found
}

func (t *Target) AddClientContext(newClientContext ClientContext) error {
	payload, err := json.Marshal(newClientContext.Token)
	if err != nil {
		return err
	}

	if err := keyring.Set(t.Name, newClientContext.ClientID, string(payload)); err != nil {
		return err
	}
	if t.ClientContexts == nil {
		t.ClientContexts = make(map[string]struct{})
	}
	t.ClientContexts[newClientContext.ClientID] = struct{}{}
	return nil
}
