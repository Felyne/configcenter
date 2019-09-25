package config_center

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coreos/etcd/clientv3"
)

const cfgContent = `
a=b
c=1
`

func TestFlow(t *testing.T) {
	var err error
	key := "xx125"

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 10 * time.Second,
	})
	assert.Equal(t, nil, err)
	defer cli.Close()
	cc := NewConfigCenter(cli, "test")

	assert.Equal(t, nil, err)
	err = cc.SetConfig(key, cfgContent)
	assert.Equal(t, nil, err)
	actual, err := cc.GetConfig(key)
	assert.Equal(t, nil, err)
	assert.Equal(t, cfgContent, actual)
	m, err := cc.ListConfig()
	assert.Equal(t, nil, err)
	t.Log(m)
	err = cc.RemoveConfig(key)
	assert.Equal(t, nil, err)
}