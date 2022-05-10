package request

type PageRequest struct {
	Page     int64    `json:"page"`
	PageSize int64    `json:"pageSize"`
	Filter   []string `json:"filter"`
}

type PageRequestCommon struct {
	Order    string        `json:"order"`
	Page     int64         `json:"page"`
	PageSize int64         `json:"pageSize"`
	Filter   []interface{} `json:"filter"`
}

// CabinInReceive Cabin structure for input parameters
type CabinInReceive struct {
	PType    string `json:"pType"`    //
	Source   string `json:"source"`   //
	Resource string `json:"resource"` // 路径
	Domain   string `json:"domain"`   //
	Method   string `json:"method"`   //
}

type DataDelete struct {
	Ids string `json:"ids"`
}
