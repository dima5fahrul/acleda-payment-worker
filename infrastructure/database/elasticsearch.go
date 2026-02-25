package database

// Mock ElasticsearchClient to fix compilation errors until fully implemented.
// This matches the usage `database.ElasticsearchClient.Index("...").Request(&model).Do(ctx)`
type elasticsearchClientMock struct{}

var ElasticsearchClient *elasticsearchClientMock = nil

func (m *elasticsearchClientMock) Index(index string) *elasticsearchIndexRequestMock {
	return &elasticsearchIndexRequestMock{}
}

type elasticsearchIndexRequestMock struct{}

func (req *elasticsearchIndexRequestMock) Request(data interface{}) *elasticsearchIndexRequestMock {
	return req
}

func (req *elasticsearchIndexRequestMock) Do(ctx interface{}) (interface{}, error) {
	return nil, nil
}
