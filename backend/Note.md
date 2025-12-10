What context is in Go

context is a standard library mechanism for controlling the lifecycle of operations. It carries three categories of information:

Cancellation signals
Allows a parent operation to stop all child operations.

Deadlines / timeouts
Ensures work does not run longer than intended.

Request-scoped values (lightweight metadata)
Trace IDs, user IDs, request IDs, etc.
Not meant for business data.

Why Go needs context

Go uses concurrency heavily (goroutines, servers, background workers). Problems arise if:

A request is cancelled but goroutines keep running.

A DB query hangs without a timeout.

A worker starts something but has no way to stop it when upstream shuts down.

context solves these by giving every goroutine a shared, chainable cancellation system.

How context works (engineering perspective)

1. ctx := context.Background()

Root context. It is never cancelled.
Used to start chains in main(), initialization, or tests.

2. ctx, cancel := context.WithCancel(parent)

Creates a child context with cancellation ability.
Calling cancel() stops downstream work.

Use case:
Stop DB queries, goroutines, or network calls when a request ends.

3. ctx, cancel := context.WithTimeout(parent, 2\*time.Second)

Automatically triggers cancellation after a time limit.
Common around network calls or DB queries.

4. ctx, cancel := context.WithDeadline(parent, time.Time)

Same idea but with a fixed point in time.

5. ctx = context.WithValue(parent, key, value)

Attaches metadata.
Typical: request IDs, correlation IDs, auth data.

Not recommended for passing domain objects or large data.

Why it's powerful (propagation model)

Cancellation flows down to all children:

client request cancelled
↓
HTTP handler context cancelled
↓
DB query context cancelled
↓
All goroutines spawned using that ctx stop

This makes system behavior predictable, resource-safe, and easy to test.

---

Client → chi Router → Handler.Method → Service → Repo → DB

# 2) Responsibilities of Each Layer (very important)

## 2.1 Handler (HTTP layer)

**Role:**

- Reads HTTP request JSON
- Validates basic input
- Calls service functions
- Sends JSON response back

**Handler does NOT:**

- Hash passwords
- Talk to the database
- Implement business rules

Handlers only handle HTTP.

**You can think:
Handler ≈ “Controller”**

## 2.2 Service (Business logic layer)

**Role:**

- Hash password
- Validate business rules
- Check email existing or not
- Call repository
- Build final model

Service does NOT know anything about HTTP.

**Service ≈ “Brain of your feature”**

## 2.3 Repository (Database Access Layer)

**Role:**

- Executes SQL queries
- Inserts / updates / fetches rows
- Converts DB rows to Go structs

Repository does NOT know passwords, validation, HTTP, or rules.

**Repository ≈ “Database worker”**

## 2.4 Database

- Stores the final teacher row.

# 3) Flow for Creating a Teacher (simple storytelling)

### Step 1 — Handler receives the HTTP request

Your handler reads JSON from the HTTP body:

```go
var req CreateTeacherRequest
json.NewDecoder(r.Body).Decode(&req)
```

It converts JSON → Go struct.

### Step 2 — Handler sends the data to Service

Handler calls service:

```go
id, err := teacherService.Register(ctx, req)
```

And now the service takes over.

### Step 3 — Service hashes password

Example:

```go
hash := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
```

Why service?
Because password hashing is **business logic**, not HTTP and not DB work.

### Step 4 — Service calls Repository

Service builds a `model.Teacher` struct and sends it to the repository:

```go
repo.Create(ctx, &teacher)
```

### Step 5 — Repository writes SQL

Repository executes SQL:

```sql
INSERT INTO teachers (...) VALUES (...);
```

It returns the new teacher ID.

---

### Step 6 — Service returns result to Handler

Handler receives the ID and creates JSON response:

```json
{
  "id": "uuid-value",
  "message": "teacher created"
}
```

### Step 7 — Handler sends JSON back to client

The client sees the result.

# 4) Why this architecture is professional (for your portfolio)

| Layer      | Why it matters                                |
| ---------- | --------------------------------------------- |
| Handler    | Keeps HTTP clean and isolated.                |
| Service    | Business logic becomes readable and testable. |
| Repository | SQL is contained and modular.                 |
| DB         | Clean schema design.                          |

This structure is exactly what senior teams use in real Go projects and what interviewers expect to see.
