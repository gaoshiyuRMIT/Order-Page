#### GET /api/orders/search

##### URL query parameters
- PartOfName: string
- DateFrom: string, UTC, RFC3339
- DateTill: string, UTC, RFC3339
- PageSize: int, default 5
- PageNo: int, default 1

##### response (application/json)
###### fields
- OrderName: string
- OrderID: int
- CustomerID: string
- CustomerName: string
- CustomerCompany: string
- TotalAmount: float
- DeliveredAmount: float
- CreatedAt: string, UTC, RFC3339