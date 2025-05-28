# MyWallet API 💸

**MyWallet API** is a personal finance tracker API that allows users to create their own wallets, record transactions, and compose financial reports.

## Features
- User registration and login with JWT authentication
- Income and expense tracking
- Date-filtered transaction queries
- Monthly financial summary reports

## Tech Stack
- Go + Gin Framework
- PostgreSQL (via Docker)
- JWT Authentication
- Docker Compose for local development

## API Endpoints

### 1. **User Management**
<details>
<summary><code>POST /signin</code></summary><br>
Create a new user.

**Request**

```
{
    "username": "dummy_user",
    "password": "dummy_password"
}
```

**Response**

```
{
  "id": 1,
  "username": "dummy_user",
  "password": "dummy_password",
  "created_at": "2025-05-28T02:24:11.64697Z"
}
```
</details>

<details>
<summary><code>POST /login</code></summary><br>
Authenticate a user and return a JWT Token

**Request**

```
{
    "username": "dummy_user",
    "password": "dummy_password"
}
```

**Response**

```
{
  "message": "Success",
  "token": [Your JWT]
}
```
</details>

<details>
<summary><code>GET /profile</code></summary><br>
Fetch user profile.

**Request**

```
Authorization: [Your JWT]
```

**Response**

```
{
  "id": 1,
  "username": "dummy_user",
  "balance_usd": 1000.15
}
```
</details>

### 2. Transaction **Management**
<details>
<summary><code>POST /add_income</code></summary><br>

Add a new income entry to user account.

**Request**

```
Authorization: [Your JWT]

{
    "amount_usd": 5000.00,
    "trx_category": "Salary"
}
```

**Response**

```
{
  "id": 48,
  "user_id": 1,
  "amount": 500000,
  "trx_type": "income",
  "trx_category": "Salary",
  "created_at": "2025-05-28T02:32:45.555731Z"
}
```
</details>



<details>
<summary><code>POST /add_expense</code></summary><br>

Add a new expense entry to user account.

**Request**

```
Authorization: [Your JWT]
{
    "amount_usd": 20,
    "trx_category": "Subscriptions"
}
```

**Response**

```
{
  "id": 49,
  "user_id": 1,
  "amount": 2000,
  "trx_type": "expense",
  "trx_category": "Subscriptions",
  "created_at": "2025-05-28T02:35:43.009377Z"
}
```

</details>

<details>
<summary><code>GET /transactions</code></summary><br>

List all user transactions

**Request**

```
Authorization: [Your JWT]
```

**Response**

```
[
  {
    "id": 48,
    "amount_usd": 5000,
    "trx_category": "Salary",
    "created_at": "2025-05-28T02:52:36.925028Z"
  },
  {
    "id": 49,
    "amount_usd": -20,
    "trx_category": "Subscriptions",
    "created_at": "2025-05-28T02:53:00.432178Z"
  }
]
```

</details>

<details>
<summary><code>GET /transactions/:date</code></summary><br>

List all user transactions by date

**Request**

```
GET http://localhost:8080/transactions/date?from_date=2025-05-23&to_date=2025-05-29
Authorization: [Your JWT]
```

**Response**

```
[
  {
    "id": 48,
    "amount_usd": 5000,
    "trx_category": "Salary",
    "created_at": "2025-05-28T02:52:36.925028Z"
  },
  {
    "id": 49,
    "amount_usd": -20,
    "trx_category": "Subscriptions",
    "created_at": "2025-05-28T02:53:00.432178Z"
  }
]
```

</details>

<details>
<summary><code>GET /transactions/:id</code></summary><br>

View a specific user transaction entry.

**Request**

```
GET http://localhost:8080/transactions/51
Authorization: [Your JWT]
```

**Response**

```
{
  "id": 51,
  "amount_usd": -20,
  "trx_category": "Subscriptions",
  "created_at": "2025-05-28T02:53:00.432178Z"
}
```

</details>

<details>
<summary><code>PUT /transactions/:id</code></summary><br>

Update an transaction entry.

**Request**

```
PUT http://localhost:8080/update_transaction/50
Authorization: [Your JWT]

{
    "amount_usd": 2000.00,
    "trx_category": "Salary"
}
```

**Response**

```
{
  "message": "Transaction updated successfully"
}
```

</details>

<details>
<summary><code>DELETE /transactions/:id</code></summary><br>

Delete an transaction entry.

**Request**

```
DELETE http://localhost:8080/delete_transaction/51
Authorization: [Your JWT]
```

**Response**

```
{
  "message": "Transaction deleted successfully"
}
```

</details>

### 3. **Reports**

<details>
<summary><code>GET /summary?month=2025-04</code></summary><br>

Return month income, expenses, and balance.

**Request**

```
GET http://localhost:8080/month_summary?month=2025-05
Authorization: [Your JWT]
```

**Response**

```
{
  "balance_usd": 2000,
  "expenses_usd": 0,
  "income_usd": 2000
}
```

</details>

## Installation