package db

import "github.com/olivere/elastic"

var EsClient *elastic.Client

func NewEsClient(address []string, username, password string) (esClient *elastic.Client, err error) {
	optList := []elastic.ClientOptionFunc{elastic.SetSniff(false)}

	optList = append(optList, elastic.SetURL(address...))

	if username != "" || password != "" {
		optList = append(optList, elastic.SetBasicAuth(username, password))
	}

	esClient, err = elastic.NewSimpleClient(optList...)

	return
}
