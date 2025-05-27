# MyWallet API 💸

**MyWallet API** is personal finance tracker API, that allows users create their own wallets, record transactions and compose reports. 

## Handlers

### 1. **User Management**

- `POST /signin` — Create a new user.
- `POST /login` — Authenticate a user and return a JWT Token
- `GET /profile` — Fetch user profile.

### 2. Transaction **Management**

- `POST /add_income` — Add a new income entry to user account.
- `POST /add_expense` — Add a new expense entry to user account.
- `GET /transactions` — List all user transactions
- `GET /transactions/:date` — List all user transactions by date
- `GET /transactions/:id` — View a specific user transaction entry.
- `PUT /transactions/:id` — Update an transaction entry.
- `DELETE /transactions/:id` — delete an transaction entry.

### 3. **Reports**

- `GET /summary?month=2025-04` — Return month income, expenses, and balance.

## Installation