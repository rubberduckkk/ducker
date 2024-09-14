package aidoc

type QueryParam struct {
	NumDocuments int
}

type QueryOption func(*QueryParam)

func WithNumDocuments(numDocuments int) QueryOption {
	return func(o *QueryParam) {
		o.NumDocuments = numDocuments
	}
}

func defaultQueryParam() *QueryParam {
	return &QueryParam{
		NumDocuments: 3,
	}
}
