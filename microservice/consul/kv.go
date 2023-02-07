package consul

import (
	"errors"

	"github.com/hashicorp/consul/api"
)

type KVClient struct {
	*api.KV
}

func NewKvClient(c *ConsulClient) *KVClient {
	return &KVClient{
		c.client.KV(),
	}
}

func (c *KVClient) PutKV(key, value string) error {
	p := &api.KVPair{Key: key, Value: []byte(value)}
	_, err := c.Put(p, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *KVClient) GetKV(key string) (string, error) {
	p, _, err := c.Get(key, nil)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", errors.New("key does not exist")
	}
	return string(p.Value), nil
}
