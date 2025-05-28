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

#### POST /signin
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

#### POST /login
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

#### GET /profile
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

### 2. Transaction **Management**

#### POST /add_income
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
<details>
<summary>POST /add_expense</summary><br>


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
  "user_id": 16,
  "amount": 2000,
  "trx_type": "expense",
  "trx_category": "Subscriptions",
  "created_at": "2025-05-28T02:35:43.009377Z"
}
```

</details>


- `GET /transactions` — List all user transactions
- `GET /transactions/:date` — List all user transactions by date
- `GET /transactions/:id` — View a specific user transaction entry.
- `PUT /transactions/:id` — Update an transaction entry.
- `DELETE /transactions/:id` — delete an transaction entry.

### 3. **Reports**

- `GET /summary?month=2025-04` — Return month income, expenses, and balance.

## Installation