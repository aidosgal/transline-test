package entity

type (
	CreateReq struct {
		Route    string            `json:"route"`
		Price    int               `json:"price"`
		Customer CreateCustomerReq `json:"customer"`
	}

	CreateCustomerReq struct {
		IDN string `json:"idn"`
	}

	CreateResp struct {
		Shipment
	}
)
