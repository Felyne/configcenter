package config_center

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var ErrNoConfig = errors.New("no config in etcd")

const (
	PathSeparator  = "/"
	ConfigBasePath = "/config_center"
	contextTimeout = 15 * time.Second
)

type ConfigCenter struct {
	etcdClient *clientv3.Client
	envName    string
}

func New(client *clientv3.Client, envName string) *ConfigCenter {
	envName = strings.Trim(envName, " \t\r\n")
	if "" == envName {
		envName = "_"
	}
	return &ConfigCenter{
		etcdClient: client,
		envName:    envName,
	}
}

func (cc *ConfigCenter) GetConfig(cfgName string) (string, error) {
	cfgPath := cc.genPath(cfgName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		contextTimeout)
	resp, err := cc.etcdClient.Get(ctx, cfgPath)
	cancel()
	if nil != err {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", ErrNoConfig
	}
	return string(resp.Kvs[0].Value), err
}

func (cc *ConfigCenter) SetConfig(cfgName string, content string) error {
	cfgPath := cc.genPath(cfgName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		contextTimeout)
	_, err := cc.etcdClient.Put(ctx, cfgPath, content)
	cancel()
	return err
}

func (cc *ConfigCenter) RemoveConfig(cfgName string) error {
	cfgPath := cc.genPath(cfgName)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		contextTimeout)
	_, err := cc.etcdClient.Delete(ctx, cfgPath)
	cancel()
	return err
}

func (cc *ConfigCenter) ListConfig() (map[string]string, error) {
	basePath := strings.Join(
		[]string{ConfigBasePath, cc.envName, ""}, PathSeparator)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		contextTimeout)

	resp, err := cc.etcdClient.Get(ctx, basePath,
		clientv3.WithPrefix())
	cancel()
	if nil != err {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	retMap := make(map[string]string)
	for _, ev := range resp.Kvs {
		key := string(ev.Key)
		pos := strings.LastIndex(key, PathSeparator)
		if pos == len(key)-1 {
			continue
		}
		name := key[pos+1:]
		retMap[name] = string(ev.Value)
	}
	return retMap, nil
}

func (cc *ConfigCenter) genPath(cfgName string) string {
	return strings.Join(
		[]string{ConfigBasePath, cc.envName, cfgName}, PathSeparator)
}
