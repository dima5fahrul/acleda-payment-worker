# Acleda Worker API Documentation

## Create Payment Link

### Request
```bash
curl -X POST http://localhost:8080/api/v1/acleda/payment-links \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "100.00",
    "currency": "USD",
    "description": "Test payment",
    "customer_name": "John Doe",
    "customer_email": "john@example.com",
    "customer_phone": "+85512345678",
    "return_url": "https://example.com/success",
    "callback_url": "https://example.com/callback",
    "merchant": "MERCHANT123"
  }'
```

### Response
```json
{
  "success": true,
  "data": {
    "transaction_id": "ACL-1645678901",
    "payment_url": "http://localhost:8080/payment-page/acleda/ACL-1645678901?sid=abc123&ptid=xyz789",
    "session_id": "abc123",
    "payment_token_id": "xyz789",
    "amount": "100.00",
    "currency": "USD",
    "status": "PENDING",
    "expires_at": "2026-02-24T11:06:00Z",
    "created_at": "2026-02-24T10:06:00Z"
  }
}
```

## Get Payment Status

### Request
```bash
curl -X GET http://localhost:8080/api/v1/acleda/payments/ACL-1645678901/status
```

### Response
```json
{
  "success": true,
  "data": {
    "id": "ACL-1645678901",
    "transaction_id": "ACL-1645678901",
    "merchant_id": "MERCHANT123",
    "session_id": "abc123",
    "payment_token_id": "xyz789",
    "description": "Test payment",
    "amount": 100.00,
    "currency": "USD",
    "invoice_id": "ACL-1645678901",
    "status": "PENDING",
    "expiry_time": 60,
    "created_at": "2026-02-24T10:06:00Z",
    "updated_at": "2026-02-24T10:06:00Z",
    "payment_id": "uuid-here",
    "payment_method_id": "uuid-here",
    "country_id": "uuid-here",
    "merchant_code": "MERCHANT123",
    "currency_id": "uuid-here",
    "purchase_amount": 100.00,
    "purchase_date": 1645678901,
    "quantity": 1,
    "confirm_date": 0,
    "purchase_type": 1,
    "save_token": 0,
    "fee_amount": 2.50,
    "tx_direction": 1,
    "return_url": "https://example.com/success",
    "error_url": "https://example.com/callback",
    "request_json": "...",
    "response_json": "..."
  }
}
```

## Payment Page

### Direct URL
```bash
curl -X GET "http://localhost:8080/payment-page/acleda/ACL-1645678901?sid=abc123&ptid=xyz789"
```

This will return the HTML payment page that auto-submits to Acleda.

## Error Responses

### Bad Request (400)
```json
{
  "error": "Invalid request body",
  "details": "invalid character '}' looking for beginning of object key string"
}
```

### Missing Required Field (400)
```json
{
  "error": "Amount is required"
}
```

### Payment Link Not Found (404)
```json
{
  "error": "Payment link not found",
  "details": "record not found"
}
```

### Internal Server Error (500)
```json
{
  "error": "Failed to create payment link",
  "details": "failed to open session: connection refused"
}
```

## Testing Flow

1. **Create Payment Link**
   ```bash
   curl -X POST http://localhost:8080/api/v1/acleda/payment-links \
     -H "Content-Type: application/json" \
     -d '{"amount": "50.00", "currency": "USD", "merchant": "TEST123"}'
   ```

2. **Get Payment URL from response**
   ```bash
   # Extract payment_url from response and open in browser
   # Example: http://localhost:8080/payment-page/acleda/ACL-1645678901?sid=abc123&ptid=xyz789
   ```

3. **Check Payment Status**
   ```bash
   curl -X GET http://localhost:8080/api/v1/acleda/payments/ACL-1645678901/status
   ```

## Required Fields

- `amount` (string, required) - Payment amount
- `currency` (string, required) - Currency code (USD, KHR, etc)
- `merchant` (string, required) - Merchant ID

## Optional Fields

- `description` (string) - Payment description
- `customer_name` (string) - Customer name
- `customer_email` (string) - Customer email
- `customer_phone` (string) - Customer phone
- `return_url` (string) - Success return URL
- `callback_url` (string) - Callback/notification URL

## Notes

- Transaction ID is auto-generated with format `ACL-{timestamp}`
- Session expires after 60 minutes
- Payment page auto-submits to Acleda after 500ms
- All payment data is stored in YugabyteDB for tracking
