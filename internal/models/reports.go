package models

type ServiceMonthlyReportReq struct {
	DateFrom string `json:"date-from"`
	DateTo   string `json:"date-to"`
}

type ServiceMonthlyReportResponse struct {
	FileLink string `json:"csv-file-link"`
}

type GetTransactionsReq struct {
	UserId    int64  `json:"user-id"`
	DateFrom  string `json:"date-from"`
	DateTo    string `json:"date-to"`
	SortBy    string `json:"sort-by"`
	SortOrder string `json:"sort-order"`
	Page      int    `json:"page"`
}

type GetTransactionsResponse struct {
	FileLink string `json:"csv-file-link"`
}
