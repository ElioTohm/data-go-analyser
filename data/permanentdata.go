package data

type ProcessedData struct {
	Type              string  `parquet:"name=type, type=UTF8"`
	CustomerReference string  `parquet:"name=customer_reference, type=UTF8"`
	NumberOfOrders    int32   `parquet:"name=number_of_orders, type=INT32"`
	TotalAmountSpent  float32 `parquet:"name=total_amount_spent, type=FLOAT"`
}

type ErrorMessage struct {
	Type           string `parquet:"name=type, type=BYTE_ARRAY, encoding=PLAIN"`
	OrderReference string `parquet:"name=order_reference, type=BYTE_ARRAY, encoding=PLAIN"`
	Message        string `parquet:"name=message, type=BYTE_ARRAY, encoding=PLAIN"`
}
