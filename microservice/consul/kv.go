package consul

import (
	"errors"

	"github.com/hashicorp/consul/api"
	"github.com/hoquangnam45/pharmacy-common-go/util"
)

type KVClient struct {
	*api.KV
}

func NewKvClient(c *Client) *KVClient {
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
		util.SugaredLogger.Infof("key %s in consul not exist", key)
		return "", errors.New("key does not exist")
	}
	return string(p.Value), nil
}
