package db

import (
	"DSearch/config"
	"DSearch/logger"
	"github.com/elastic/go-elasticsearch/v7"
	json "github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
)

var elasticsearchClientConfigs map[string]elasticsearch.Config

type elasticsearchClientConfig struct {
	Addresses []string // A list of Elasticsearch nodes to use.
	Username  string   // Username for HTTP Basic Authentication.
	Password  string   // Password for HTTP Basic Authentication.
}

func ConnectElasticsearch() {

	elasticsearchConfig := config.Database.Elasticsearch

	elasticsearchClientConfigs = make(map[string]elasticsearch.Config)
	for k, db := range elasticsearchConfig {
		var dbConf elasticsearchClientConfig
		if err := mapstructure.Decode(db, &dbConf); err != nil {
			logger.SysLog(nil).Panic("elasticsearch日志：" + config.AppMode() + ".elasticsearch[请检查配置是否正确]，err：" + err.Error())
		}

		elasticsearchClientConfigs = map[string]elasticsearch.Config{
			k: {
				Addresses: dbConf.Addresses,
				Username:  dbConf.Username,
				Password:  dbConf.Password,
			},
		}
		es, err := EsKey(k)
		if err != nil {
			logger.SysLog(nil).Panic("elasticsearch日志：" + "链接错误：" + err.Error())
		}

		// 1. Get cluster info
		//
		res, err := es.Info()
		if err != nil {
			logger.SysLog(nil).Panic("elasticsearch日志：" + "获取info信息失败：" + err.Error())
		}

		// Check response status
		if res.IsError() {
			logger.SysLog(nil).Panic("elasticsearch日志：" + "获取info信息失败：" + res.String())
		}

		r := make(map[string]interface{})
		// Deserialize the response into a map.
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			logger.SysLog(nil).Panic("elasticsearch日志：" + "信息解析失败：" + err.Error())
		}

		// Print client and server version numbers.
		logger.SysLog("elasticsearchClient注册成功，ClientV：" + elasticsearch.Version + "ServerV：" + r["version"].(map[string]interface{})["number"].(string)).Info("elasticsearch日志")
		_ = res.Body.Close() // 用完随手关闭是个好习惯
	}
}

func Es() (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearchClientConfigs["default"])
}

// EsKey 获取指定数据链接(切换数据库)
func EsKey(key string) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearchClientConfigs[key])
}
