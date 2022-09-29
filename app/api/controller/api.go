package controller

import (
	"DSearch/app/elasticUtil"
	"DSearch/core"
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

	aa := elasticUtil.Elastic{
		Ctx: ctx.Context,
	}
	aa.CreateIndex(indexName, mapping)

	bb, _ := aa.GetMapping(indexName)

	aa.DeleteIndex(indexName)
	ctx.Json(0, "ok", bb)
	return
}
