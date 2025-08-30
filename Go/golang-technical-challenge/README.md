# 🧾 Invoice API – CRUD & Excel Import

A Go REST API for managing invoices and their products using PostgreSQL, with Excel import support.

---

## ⚙️ Project Setup

### 1️⃣ Create PostgreSQL Database

Ensure PostgreSQL is installed and running. Then, create a new database:

```sql
CREATE DATABASE invoice_db;
```

### 2️⃣ Create `.env` File

In the root directory of the project, create a `.env` file with the following structure:

```env
# APP CONFIG
APP_NAME=invoice-service
APP_ENV=development
APP_PORT=3000
WEB_PREFORK=false

# LOG CONFIG
LOG_LEVEL=6

# DATABASE CONFIG (POSTGRES)
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_postgres_user
DB_PASSWORD=your_password
DB_NAME=invoice_db
DB_SSLMODE=disable

# Database Connection Pool
DB_POOL_IDLE=5
DB_POOL_MAX=20
DB_POOL_LIFETIME=300
```

> ✅ **Tip**: You may copy this to a `.env.example` file for team sharing and exclude `.env` in `.gitignore`.

---

### 3️⃣ Install Go Dependencies

Ensure you’re inside the project root directory, then run:

```bash
go mod tidy
```

This will download and install all required dependencies listed in go.mod.

---

### 4️⃣ Run Database Migrations

Install `golang-migrate` if you haven’t:

```bash
brew install golang-migrate        # for macOS
# or
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Then run migrations:

```bash
migrate -path db/migrations -database "postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/DB_NAME?sslmode=disable" up

```

🔁 Example:

```bash
migrate -path db/migrations -database "postgres://postgres:1234@localhost:5432/invoice_db?sslmode=disable" up
```

---

## 🔧 Tech Stack

- **Language**: Go 1.24.5
- **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL (via `gorm.io/gorm`)
- **ORM**: [GORM](https://gorm.io/)
- **Excel Import**: [Excelize](https://github.com/xuri/excelize)
- **Logging**: [Logrus](https://github.com/sirupsen/logrus)
- **Validation**: [Validator v10](https://github.com/go-playground/validator)
- **Decimal Support**: [shopspring/decimal](https://github.com/shopspring/decimal)
- **Nullable Types**: [guregu/null](https://github.com/guregu/null)
- **Configuration**: [Viper](https://github.com/spf13/viper)

---

## 📌 Base URL

```
http://localhost:3000/api/invoices
```

---

## 📥 1. Import Invoices from Excel

**POST** `/import`

Uploads an `.xlsx` file containing invoices and product data.

### ✅ Postman
- Method: `POST`
- URL: `http://localhost:3000/api/invoices/import`
- Body: `form-data`
  - Key: `file` (type: file)
  - Value: Upload `2. InvoiceImport.xlsx`

### 🌀 curl

```bash
curl -X POST http://localhost:3000/api/invoices/import   -H "Content-Type: multipart/form-data"   -F "file=@2. InvoiceImport.xlsx"
```

---

## 📄 2. Get Invoices (Read)

**GET** `/?date=YYYY-MM-DD&page=1&size=5`

Returns paginated invoice list, total profit, and total cash transactions for a given date.

### ✅ Postman
- Method: `GET`
- URL:  
  ```
  http://localhost:3000/api/invoices?date=2025-08-25&page=1&size=5
  ```

### 🌀 curl

```bash
curl "http://localhost:3000/api/invoices?date=2025-08-25&page=1&size=5"
```

---

## 🆕 3. Create Invoice

**POST** `/`

Creates a new invoice with products.

### ✅ Postman
- Method: `POST`
- URL: `http://localhost:3000/api/invoices`
- Body: `raw` JSON

```json
{
  "invoice_no": "INV-1005",
  "date": "2025-08-25",
  "customer_name": "Edwardo Samosir",
  "salesperson_name": "Andi Wijaya",
  "payment_type": "CASH",
  "notes": "First order",
  "products": [
    {
      "item_name": "iPhone 15 Pro",
      "quantity": 2,
      "total_cost": 25000000,
      "total_price": 30000000
    },
    {
      "item_name": "MacBook Pro M3",
      "quantity": 1,
      "total_cost": 35000000,
      "total_price": 42000000
    }
  ]
}
```

### 🌀 curl

```bash
curl -X POST http://localhost:3000/api/invoices   -H "Content-Type: application/json"   -d '{
    "invoice_no": "INV-1005",
    "date": "2025-08-25",
    "customer_name": "Edwardo Samosir",
    "salesperson_name": "Andi Wijaya",
    "payment_type": "CASH",
    "notes": "First order",
    "products": [
      {
        "item_name": "iPhone 15 Pro",
        "quantity": 2,
        "total_cost": 25000000,
        "total_price": 30000000
      },
      {
        "item_name": "MacBook Pro M3",
        "quantity": 1,
        "total_cost": 35000000,
        "total_price": 42000000
      }
    ]
  }'
```

---

## 🔁 4. Update Invoice

**PUT** `/:invoiceNo`

Updates an existing invoice by `invoice_no`.

### ✅ Postman
- Method: `PUT`
- URL: `http://localhost:3000/api/invoices/INV-1005`
- Body: `raw` JSON

```json
{
  "date": "2025-08-30",
  "customer_name": "Edwardo S",
  "salesperson_name": "Budi Santoso",
  "payment_type": "CREDIT",
  "notes": "Updated order",
  "products": [
    {
      "item_name": "MacBook Pro M3 Max",
      "quantity": 10,
      "total_cost": 40000000,
      "total_price": 45000000
    },
    {
      "item_name": "iPad Pro 13 M4",
      "quantity": 2,
      "total_cost": 30000000,
      "total_price": 35000000
    }
  ]
}
```

### 🌀 curl

```bash
curl -X PUT http://localhost:3000/api/invoices/INV-1005   -H "Content-Type: application/json"   -d '{
    "date": "2025-08-30",
    "customer_name": "Edwardo S",
    "salesperson_name": "Budi Santoso",
    "payment_type": "CREDIT",
    "notes": "Updated order",
    "products": [
      {
        "item_name": "MacBook Pro M3 Max",
        "quantity": 10,
        "total_cost": 40000000,
        "total_price": 45000000
      },
      {
        "item_name": "iPad Pro 13 M4",
        "quantity": 2,
        "total_cost": 30000000,
        "total_price": 35000000
      }
    ]
  }'
```

---

## ❌ 5. Delete Invoice

**DELETE** `/:invoiceNo`

Deletes an invoice by `invoice_no`.

### ✅ Postman
- Method: `DELETE`
- URL:  
  ```
  http://localhost:3000/api/invoices/INV-1005
  ```

### 🌀 curl

```bash
curl -X DELETE http://localhost:3000/api/invoices/INV-1005
```

---

## ✅ Validation Rules

- `invoice_no`, `date`, `customer_name`, `salesperson_name`, `payment_type` → **required**
- `payment_type` must be either: `"CASH"` or `"CREDIT"`
- Each product must contain:
  - `item_name` (string)
  - `quantity` (integer)
  - `total_cost` and `total_price` (numeric)

---

## 📂 Excel Import Format

Ensure your `.xlsx` file includes **two sheets**:

- `invoice`
- `product sold`

Refer to the sample file: `InvoiceImport.xlsx`

---