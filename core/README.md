# Familiz

> Community management system for large families: members, contributions, events & yearly archiving.

**Tech:** Go, SQLite, JWT (bcrypt, chi router)

---

## Quick Start

```bash
git clone <repo>
cd familiz
go mod tidy
cp .env.example .env  # add JWT_SECRET
go run cmd/api/main.go
```

Server runs at `http://localhost:8080`

### Environment

```env
JWT_SECRET=your-secret-key
```

---

## API Endpoints

All endpoints except `/register` and `/login` require:  
`Authorization: Bearer <token>`

### Auth
| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/register` | `email, password, first_name, last_name, birth_date, marital_status` |
| POST | `/login` | `email, password` → returns `token` |

### Members
| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/members` | `first_name, last_name, birth_date, marital_status` |
| GET | `/members` | List all |
| PUT | `/members/{id}` | Update fields |
| DELETE | `/members/{id}` | Delete (cascade) |

### Transactions
| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/transactions` | `member_id, month, year, amount, note` (amount `0` = auto-calc) |
| GET | `/transactions` | `?member_id=X&archived=true` |
| PUT | `/transactions/{id}` | `month, year, amount, note` |
| DELETE | `/transactions/{id}` | Delete |

### Events
| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/events` | `member_id, type, amount_received, event_date` (amount `0` = auto-calc) |
| GET | `/events` | `?member_id=X&archived=true` |
| PUT | `/events/{id}` | `type, amount_received, event_date` |
| DELETE | `/events/{id}` | Delete |

### Profile
| Method | Endpoint | Query |
|--------|----------|-------|
| GET | `/profile/{id}` | `?archived=true` (includes archives) |

### Settings
| Method | Endpoint | Body |
|--------|----------|------|
| GET | `/settings/contributions` | – |
| PUT | `/settings/contributions` | `amount_single, amount_married, amount_minor` |
| GET | `/settings/events` | – |
| PUT | `/settings/events/{type}` | `default_amount` |

### Archiving
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/settings/archive` | Archive current year |
| POST | `/settings/unarchive` | Revert last archive (only if current year empty) |

---

## 🧪 Postman Test Workflow

### Step 1: Register Admin
- **Method:** `POST`
- **URL:** `http://localhost:8080/register`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
```json
{
  "email": "admin@familiz.com",
  "password": "SecurePass123",
  "first_name": "Toto",
  "last_name": "Sall",
  "birth_date": "1990-01-01",
  "marital_status": "married"
}
```
- **Expected Response:** `201 Created`

---

### Step 2: Login → Get Token
- **Method:** `POST`
- **URL:** `http://localhost:8080/login`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**
```json
{
  "email": "admin@familiz.com",
  "password": "SecurePass123"
}
```
- **Expected Response:** `200 OK` with `token`
- **Action:** Copy the `token` from the response. You'll use it in all subsequent requests.

---

### Step 3: Add Member
- **Method:** `POST`
- **URL:** `http://localhost:8080/members`
- **Headers:**
  - `Content-Type: application/json`
  - `Authorization: Bearer <YOUR_TOKEN>`
- **Body (raw JSON):**
```json
{
  "first_name": "Mamadou",
  "last_name": "Diallo",
  "birth_date": "1985-05-20",
  "marital_status": "single"
}
```
- **Expected Response:** `201 Created` with `member_id` (note this ID for next steps)

---

### Step 4: Add Transaction (Auto-calculated)
- **Method:** `POST`
- **URL:** `http://localhost:8080/transactions`
- **Headers:**
  - `Content-Type: application/json`
  - `Authorization: Bearer <YOUR_TOKEN>`
- **Body (raw JSON):**
```json
{
  "member_id": 2,
  "month": 7,
  "year": 2026,
  "amount": 0,
  "note": "July 2026 contribution"
}
```
- **Expected Response:** `201 Created` with the transaction details (amount auto-filled based on member's marital status)

---

### Step 5: Add Event (Auto-calculated)
- **Method:** `POST`
- **URL:** `http://localhost:8080/events`
- **Headers:**
  - `Content-Type: application/json`
  - `Authorization: Bearer <YOUR_TOKEN>`
- **Body (raw JSON):**
```json
{
  "member_id": 2,
  "type": "wedding",
  "amount_received": 0,
  "event_date": "2026-07-22"
}
```
- **Expected Response:** `201 Created` with the event details (amount auto-filled from event settings)

---

### Step 6: View Member Profile
- **Method:** `GET`
- **URL:** `http://localhost:8080/profile/2?archived=true`
- **Headers:** `Authorization: Bearer <YOUR_TOKEN>`
- **Expected Response:** `200 OK` with complete member profile (info + transactions + events)

---

### Step 7: Archive Current Year
- **Method:** `POST`
- **URL:** `http://localhost:8080/settings/archive`
- **Headers:** `Authorization: Bearer <YOUR_TOKEN>`
- **Expected Response:** `200 OK` with confirmation of archived year

---

### Step 8: View Archived Transactions
- **Method:** `GET`
- **URL:** `http://localhost:8080/transactions?archived=true`
- **Headers:** `Authorization: Bearer <YOUR_TOKEN>`
- **Expected Response:** `200 OK` showing all transactions including archived ones

---

### Postman Collection Tips

1. **Create an Environment** in Postman and store:
   - `base_url`: `http://localhost:8080`
   - `token`: Your JWT token from Step 2 (update after each login)

2. **Use Variables** in requests:
   - URLs: `{{base_url}}/members`
   - Headers: `Authorization: Bearer {{token}}`

3. **Collection Structure** (suggested):
   ```
   Familiz API/
   ├── Auth/
   │   ├── Register
   │   └── Login
   ├── Members/
   │   ├── Create Member
   │   ├── List Members
   │   ├── Update Member
   │   └── Delete Member
   ├── Transactions/
   │   ├── Create Transaction
   │   ├── List Transactions
   │   ├── Update Transaction
   │   └── Delete Transaction
   ├── Events/
   │   ├── Create Event
   │   ├── List Events
   │   ├── Update Event
   │   └── Delete Event
   ├── Profile/
   │   └── Get Profile
   ├── Settings/
   │   ├── Get Contribution Settings
   │   ├── Update Contribution Settings
   │   ├── Get Event Settings
   │   └── Update Event Settings
   └── Archiving/
       ├── Archive Year
       └── Unarchive Year
   ```

---

## Project Structure (essential)

```
internal/
├── apps/          # Modules: auth, members, transactions, events, settings, profile
├── services/      # Business logic (calculator.go)
├── database/      # DB connection + migrations
└── utils/         # password, contextkeys
migrations/        # 001_init.sql (schema)
cmd/api/main.go    # Entry point
```

---

## Auto-Calculation Rules

| Module | Trigger | Source |
|--------|---------|--------|
| Transaction | `amount = 0` | `contribution_settings` by member's `marital_status` |
| Event | `amount_received = 0` | `event_settings` by `type` |

---

## Archiving Behavior

- `is_archived` flag added to `transactions` and `events`
- `current_year` stored in `contribution_settings`
- `GET` endpoints exclude archives by default (`?archived=true` to include)
- `PUT`/`DELETE` on archived records are **blocked**

---

## Database Schema

**Core tables:** `users`, `members`, `transactions`, `events`, `contribution_settings`, `event_settings`

**Key fields:**
- `transactions.is_archived` (BOOLEAN)
- `events.is_archived` (BOOLEAN)
- `contribution_settings.current_year` (INTEGER)

---

**License:** MIT