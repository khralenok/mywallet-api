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
<summary><code>POST /signin</code><a> - Create a new user.</a></summary><br>


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
<summary><code>POST /login</code><a> - Authenticate a user and return a JWT Token.</a></summary><br>


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
<summary><code>GET /profile</code><a> - Fetch user profile.</a></summary><br>


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
<summary><code>POST /add_income</code><a> - Add a new income entry to user account.</a></summary><br>

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
<summary><code>POST /add_expense</code><a> - Add a new expense entry to user account.</a></summary><br>

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
<summary><code>GET /transactions</code><a> - List all user transactions.</a></summary><br>

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
<summary><code>GET /transactions/:date</code><a> - List all user transactions by date</a></summary><br>



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
<summary><code>GET /transactions/:id</code><a> - View a specific user transaction entry.</a></summary><br>



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
<summary><code>PUT /transactions/:id</code><a> - Update an transaction entry.</a></summary><br>



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
<summary><code>DELETE /transactions/:id</code><a> - Delete an transaction entry.</a></summary><br>

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
<summary><code>GET /summary?month=2025-04</code><a> - Return month income, expenses, and balance.</a></summary><br>



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

1. Clone the repo
2. Create `.env` file with this template
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=mywallet
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
JWT_SECRET=your_jwt_secret
```
3. Run Docker:
```
bash
docker-compose up -d
```
4. Run app
```
go run .
```

5. Use tests from test folders with REST Client for VS Code. 

6. Example path to follow: Signin -> Login -> Add income -> See profile -> Add expense -> See Transactions -> See summary 

## Future plans

This project was originally created for learning purposes, but I see potential for it to evolve into the foundation for future products.

It can be useful for:

- Developers who want to build wallet frontends without writing the backend logic
- Integration with Telegram bots or other lightweight frontends
- Experimenting with analytics, budgeting tools, or financial summaries
- Further improvements I’m considering:

Adding transaction categories and labels

- Exporting reports to CSV or PDF
- Multi-currency support
- Admin panel or user dashboard
- Thanks for checking it out!