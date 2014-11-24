package shopify

import (
  
    "encoding/json"
  
    "fmt"
  
    "time"
  
)

type Transaction struct {
  
    Amount time.Time `json:amount`
  
    Authorization string `json:authorization`
  
    CreatedAt time.Time `json:created_at`
  
    DeviceId string `json:device_id`
  
    Gateway string `json:gateway`
  
    SourceName string `json:source_name`
  
    PaymentDetails string `json:payment_details`
  
    Id string `json:id`
  
    Kind string `json:kind`
  
    OrderId int64 `json:order_id`
  
    Receipt string `json:receipt`
  
    Status string `json:status`
  
    Test string `json:test`
  
    UserId string `json:user_id`
  
    Currency string `json:currency`
  
}


func (api *API) Transaction_index() (*[]Transaction, error) {
  res, status, err := api.request("/admin/transactions.json", "GET", nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := &map[string][]Transaction{}
  err = json.NewDecoder(res).Decode(r)

  fmt.Printf("things are: %v\n\n", *r)

  result := (*r)["products"]

	if err != nil {
		return nil, err
  }

  return &result, nil
}


// TODO implement Transaction.count

// TODO implement Transaction.show

// TODO implement Transaction.create

