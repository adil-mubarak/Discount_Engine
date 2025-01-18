# DISCOUNT ENGINE
## Overview
This project implements a Discount Engine that applies dynamic discount rules to orders based on their amount and customer type. The engine supports concurrent processing, priority-based rule application, and various discount types.

**Run the HTTP server**
-starts the server:
```
go run server/server.go
```

-This server will be available at `http://localhost:8080`,

**API Testing**;
-send a post request with json data to `/discount`;
```json
{
  "order_amount": 150,
  "customer_type": "regular"
}
```
-The response will include the discount applied,final amount,applied rules;
```json
{
  
    "discount_applied": 15,
    "final_amount": 135,
    "applied_rules": [
        "10% off for orders over $100"
    ]
}
```
**Unit test**
- Run test with:

- 
- ```
   cd unit_test,
   go test
  ```
