package consumer_data

import (
	"context"
	"fmt"
	"github.com/1340691923/xwl_bi/engine/db"
	"github.com/1340691923/xwl_bi/engine/logs"
	"github.com/olivere/elastic"
	"go.uber.org/zap"

	"strings"
)

type ClientReportData struct {
	Data    string
	TableId string
	Date    string
}

func (this *ClientReportData) Name() string {
	return "client_report_data" + this.TableId
}

func (this *ClientReportData) GetReportData() *elastic.BulkIndexRequest {
	return elastic.NewBulkIndexRequest().Index(this.CreateReportName()).Type(this.getTyp()).Doc(this.Data)
}

func (this *ClientReportData) CreateIndex() (err error) {

	indexName := this.CreateReportName()
	indexExists, err := db.EsClient.IndexExists(indexName).Do(context.Background())

	if err != nil {
		return
	}
	if !indexExists {
		db.EsClient.CreateIndex(indexName).Body(this.createIndexStr()).Do(context.Background())

		_, err = db.EsClient.Alias().Add(indexName, this.GetAliasName()).Do(context.Background())
		if err != nil {
			logs.Logger.Error("别名创建失败", zap.Error(err))
		}
	}
	return
}

//type
func (this *ClientReportData) getTyp() string {
	return "_doc"
}

//index
func (this *ClientReportData) GetAliasName() string {
	return fmt.Sprintf("%s%s", this.Name(), "_index")
}

func (this *ClientReportData) CreateReportName() string {
	return fmt.Sprintf("%v_%v", this.Name(), this.Date)
}

func (this *ClientReportData) GetList(ctx context.Context, searchKw, date string) (*elastic.SearchResult, error) {

	search := db.EsClient.Search(this.GetAliasName())

	dateArr := strings.Split(date, ",")

	q := elastic.NewBoolQuery()

	if len(dateArr) == 2 {
		q = q.Must(elastic.NewRangeQuery("create_time").
			Gte(dateArr[0]).
			Lte(dateArr[1]).IncludeLower(false).IncludeUpper(false))
	}

	if searchKw != "" {
		q = q.Must(elastic.NewMatchQuery("data", searchKw))
		highlight := elastic.NewHighlight().Field("data").PreTags("<b style='color:red'>").PostTags("</b>").NumOfFragments(0)
		search = search.Highlight(highlight)
	}

	return search.Query(q).Sort("create_time", false).From(0).Size(1000).Do(ctx)
}

//创建索引字符串
func (this *ClientReportData) createIndexStr() string {
	s := `
		{
		  "settings": {
			"number_of_replicas": 0,
			"number_of_shards": 1
		  },
		  "mappings" : {
			  "_doc" : {
				"dynamic" : "false",
				"properties": {
				"event_name": {
					"type": "keyword"
				},
				"create_time": {
					"format": "yyyy-MM-dd HH:mm:ss",
					"type": "date"
				},
				"data": {
					"type": "text",
					"analyzer":"english",
					"fields": {
							"keyword": {
							  "type": "keyword"
							}
						}
				}
			}
			  }
			}
		}
	`
	return s
}
