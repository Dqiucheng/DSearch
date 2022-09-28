package controller

import (
	"DSearch/core"
	"DSearch/db"
	"fmt"
)

type Api struct {
}

type Ass struct {
	Aa   string `validate:"required"`
	Ut   string `validate:"required"`
	Ut_a string `validate:"required"`
}

func (Api) Test(ctx *core.Context) {
	indexName := "test_index" //要创建的索引名
	type mi = map[string]interface{}

	mapping := mi{
		"settings": mi{
			"index": mi{
				"analysis.analyzer.default.type": "ik_max_word", // 设置默认分词器
			},
		},
		"mappings": mi{
			"properties": mi{
				"id": mi{ //整形字段, 允许精确匹配
					"type": "integer",
				},
				"name": mi{
					"type":            "text",     //字符串类型且进行分词, 允许模糊匹配
					"analyzer":        "ik_smart", //设置字段分词器
					"search_analyzer": "ik_smart",
					"fields": mi{ //当需要对模糊匹配的字符串也允许进行精确匹配时假如此配置
						"keyword": mi{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"date_field": mi{ //时间类型, 允许精确匹配
					"type": "date",
				},
				"keyword_field": mi{ //字符串类型, 允许精确匹配
					"type": "keyword",
				},
				"nested_field": mi{ //嵌套类型
					"type": "nested",
					"properties": mi{
						"id": mi{
							"type": "integer",
						},
						"start_time": mi{ //长整型, 允许精确匹配
							"type": "long",
						},
					},
				},
			},
		},
	}
	fmt.Println(db.Es().CreateIndex(indexName).BodyJson(mapping).Do(ctx.Context)) // 创建索引
	fmt.Println(db.Es().IndexExists(indexName).Do(ctx.Context))                   // 判断索引是否存在

	fmt.Println(db.Es().PutMapping().
		Index(indexName).
		BodyJson(mi{
			"properties": mi{
				"id": mi{ //整形字段, 允许精确匹配
					"type": "integer",
				},
			},
		}).
		Do(ctx.Context))	// 更新索引

	fmt.Println(db.Es().IndexAnalyze().Analyzer("ik_max_word").Text("哈撒给").Do(ctx.Context)) // 预览关键词分词效果
	fmt.Println(db.Es().DeleteIndex(indexName).Do(ctx.Context))                             // 删除索引

	ctx.Json(0, "ok", nil)
	return
}
