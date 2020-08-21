package htsconstants

/*
 * ENUM for acceptable HTTP Header names/keys
 */

type HttpHeaderName int

const (
	ContentTypeHeader HttpHeaderName = 0
)

const (
	contentTypeHeaderString = "Content-Type"
)

var httpHeaderNameStringMap = map[HttpHeaderName]string{
	ContentTypeHeader: contentTypeHeaderString,
}

func (e HttpHeaderName) String() string {
	return httpHeaderNameStringMap[e]
}

/*
 * ENUM for acceptable Content-Type Header values
 */

type ContentTypeHeaderValue int

const (
	ContentTypeHeaderHtsgetJSON ContentTypeHeaderValue = 0
)

const (
	contentTypeHeaderHtsgetJSONString string = "application/vnd.ga4gh.htsget.v1.2.0+json; charset=utf-8"
)

var contentTypeStringMap = map[ContentTypeHeaderValue]string{
	ContentTypeHeaderHtsgetJSON: contentTypeHeaderHtsgetJSONString,
}

func (e ContentTypeHeaderValue) String() string {
	return contentTypeStringMap[e]
}
