package db

import (
	"DSearch/config"
	"DSearch/logger"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/olivere/elastic/v7"
)

type EsWriter struct {
}

func (w EsWriter) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v)
}

var elasticsearchClientConfigs map[string]*elastic.Client

type elasticsearchClientConfig struct {
	Addresses []string // A list of Elasticsearch nodes to use.
	Username  string   // Username for HTTP Basic Authentication.
	Password  string   // Password for HTTP Basic Authentication.
}

func ConnectElasticsearch() {
	elasticsearchClientConfigs = make(map[string]*elastic.Client)
	elasticsearchConfig := config.Database.Elasticsearch

	for k, db := range elasticsearchConfig {
		var dbConf elasticsearchClientConfig
		if err := mapstructure.Decode(db, &dbConf); err != nil {
			logger.SysLog(nil).Panic("elasticsearch日志：" + config.AppMode() + ".elasticsearch[请检查配置是否正确]，err：" + err.Error())
		}

		clientOptionFunc := make([]elastic.ClientOptionFunc, 0)
		clientOptionFunc = append(clientOptionFunc, elastic.SetSniff(false)) // SetSniff启用或禁用集群嗅探器（默认情况下启用）。
		clientOptionFunc = append(clientOptionFunc, elastic.SetURL(dbConf.Addresses...))
		clientOptionFunc = append(clientOptionFunc, elastic.SetErrorLog(EsWriter{}))
		clientOptionFunc = append(clientOptionFunc, elastic.SetInfoLog(EsWriter{}))
		//clientOptionFunc = append(clientOptionFunc, elastic.SetTraceLog(EsWriter{}))
		if dbConf.Username != "" && dbConf.Password != "" {
			clientOptionFunc = append(clientOptionFunc, elastic.SetBasicAuth(dbConf.Username, dbConf.Password))
		}

		var err error
		elasticsearchClientConfigs[k], err = elastic.NewClient(clientOptionFunc...)
		if err != nil {
			logger.SysLog(nil).Panic("elasticsearch日志：" + "链接错误：" + err.Error())
		}
	}

	if len(elasticsearchConfig) > 0 {
		logger.SysLog("elasticsearchClient注册成功").Info("elasticsearch日志")
	}
}

func Es() *elastic.Client {
	return elasticsearchClientConfigs["default"]
}

// EsKey 获取指定数据链接(切换数据库)
func EsKey(key string) *elastic.Client {
	return elasticsearchClientConfigs[key]
}
