# Go Validation - Testing Best Practices

## Overview

This guide provides patterns and best practices for writing effective Go tests that meet validation requirements.

> **Scope:** Default examples use stdlib `testing` (and `go-cmp` where noted). When companion Go Test rules (stem `go-test`) are installed, follow that companion for test authoring. Optional testify sections below apply only when the package already adopted testify.

## Test Structure

### Basic Test Structure

```go
package mypackage_test

import "testing"

func setup(t *testing.T) *TestContext {
    t.Helper()
    ctx := &TestContext{}

    t.Cleanup(func() {
        ctx.Close()
    })

    return ctx
}

func TestMyFunction(t *testing.T) {
    input := "test input"
    want := "expected output"

    got := MyFunction(input)
    if got != want {
        t.Fatalf("MyFunction(%q) = %q, want %q", input, got, want)
    }
}
```

## Table-Driven Tests

### Basic Table-Driven Pattern

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {name: "positive numbers", a: 2, b: 3, want: 5},
        {name: "negative numbers", a: -1, b: -2, want: -3},
        {name: "zero", a: 0, b: 0, want: 0},
        {name: "mixed", a: -5, b: 10, want: 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Fatalf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### Advanced Table-Driven Pattern with Error Cases

```go
var (
    errTestInvalidFormat = errors.New("invalid email format")
    errTestEmptyEmail    = errors.New("email cannot be empty")
)

func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr error
    }{
        {name: "valid email", email: "user@example.com"},
        {name: "missing @", email: "userexample.com", wantErr: errTestInvalidFormat},
        {name: "empty email", email: "", wantErr: errTestEmptyEmail},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr != nil {
                if err == nil {
                    t.Fatalf("ValidateEmail(%q) error = nil, want %v", tt.email, tt.wantErr)
                }
                if !errors.Is(err, tt.wantErr) {
                    t.Fatalf("ValidateEmail(%q) error = %v, want %v", tt.email, err, tt.wantErr)
                }
                return
            }
            if err != nil {
                t.Fatalf("ValidateEmail(%q) unexpected error: %v", tt.email, err)
            }
        })
    }
}
```

## Test Helpers

### Using t.Helper()

```go
func assertUserValid(t *testing.T, user *User) {
    t.Helper()
    if user == nil {
        t.Fatal("user is nil")
    }
    if user.ID == "" {
        t.Fatal("user.ID is empty")
    }
    if user.Name == "" {
        t.Fatal("user.Name is empty")
    }
}

func createTestUser(t *testing.T, name string) *User {
    t.Helper()

    user := &User{
        ID:   generateID(),
        Name: name,
    }

    t.Cleanup(func() {
        cleanupUser(user)
    })

    return user
}
```

### Setup and Cleanup

```go
func TestDatabaseOperations(t *testing.T) {
    db := setupTestDB(t)

    user := &User{Name: "Test"}
    if err := db.SaveUser(user); err != nil {
        t.Fatalf("SaveUser() error = %v", err)
    }
}

func setupTestDB(t *testing.T) *Database {
    t.Helper()

    db, err := OpenDatabase(":memory:")
    if err != nil {
        t.Fatalf("OpenDatabase() error = %v", err)
    }

    t.Cleanup(func() {
        _ = db.Close()
    })

    return db
}
```

## Mocking and Test Doubles

### Interface-Based Stubbing

```go
type DataStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

type stubDataStore struct {
    getFunc func(key string) (string, error)
    getCalls []string
}

func (s *stubDataStore) Get(key string) (string, error) {
    s.getCalls = append(s.getCalls, key)
    if s.getFunc != nil {
        return s.getFunc(key)
    }
    return "", nil
}

func (s *stubDataStore) Set(key, value string) error {
    return nil
}

func TestMyService(t *testing.T) {
    store := &stubDataStore{
        getFunc: func(key string) (string, error) {
            if key == "test" {
                return "value", nil
            }
            return "", errors.New("not found")
        },
    }

    service := NewService(store)
    got, err := service.Process("test")
    if err != nil {
        t.Fatalf("Process(%q) error = %v", "test", err)
    }
    want := "processed value"
    if got != want {
        t.Fatalf("Process(%q) = %q, want %q", "test", got, want)
    }
    if len(store.getCalls) != 1 {
        t.Fatalf("Get calls = %d, want 1", len(store.getCalls))
    }
}
```

### Optional: testify/mock (package-adopted testify only)

Use only when the package already uses testify. Do not mix with the default stdlib plus go-cmp stack.

```go
import (
    "github.com/stretchr/testify/mock"
)

type MockDataStore struct {
    mock.Mock
}

func (m *MockDataStore) Get(key string) (string, error) {
    args := m.Called(key)
    return args.String(0), args.Error(1)
}

func (m *MockDataStore) Set(key, value string) error {
    args := m.Called(key, value)
    return args.Error(0)
}

// Usage in test
func TestWithMock(t *testing.T) {
    mockStore := new(MockDataStore)

    // Set expectations
    mockStore.On("Get", "key1").Return("value1", nil)
    mockStore.On("Set", "key2", "value2").Return(nil)

    // Test code
    service := NewService(mockStore)
    err := service.DoSomething()

    require.NoError(t, err)
    mockStore.AssertExpectations(t)
}
```

## Coverage Strategies

### Improving Coverage

1. **Test all public APIs**:

```go
// Ensure all exported functions have tests
func TestPublicAPI(t *testing.T) {
    tests := []struct {
        name string
        test func(t *testing.T)
    }{
        {"NewClient", testNewClient},
        {"Client.Connect", testClientConnect},
        {"Client.Send", testClientSend},
        {"Client.Close", testClientClose},
    }

    for _, tt := range tests {
        t.Run(tt.name, tt.test)
    }
}
```

2. **Test error paths**:

```go
func TestErrorHandling(t *testing.T) {
    tests := []struct {
        name      string
        input     string
        wantErr   bool
        setupStub func(*stubDataStore)
    }{
        {
            name:    "success",
            input:   "valid",
            wantErr: false,
            setupStub: func(s *stubDataStore) {
                s.getFunc = func(string) (string, error) { return "value", nil }
            },
        },
        {
            name:    "database error",
            input:   "test",
            wantErr: true,
            setupStub: func(s *stubDataStore) {
                s.getFunc = func(string) (string, error) { return "", errors.New("db error") }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            store := &stubDataStore{}
            tt.setupStub(store)

            _, err := Process(store, tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Fatalf("Process(%q) error = nil, want error", tt.input)
                }
            } else if err != nil {
                t.Fatalf("Process(%q) unexpected error: %v", tt.input, err)
            }
        })
    }
}
```

3. **Test edge cases**:

```go
func TestEdgeCases(t *testing.T) {
    tests := []struct {
        name  string
        input []int
        want  int
    }{
        {"empty slice", []int{}, 0},
        {"single element", []int{5}, 5},
        {"negative numbers", []int{-1, -5, -3}, -1},
        {"mixed", []int{-5, 0, 10, -3}, 10},
        {"all zeros", []int{0, 0, 0}, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FindMax(tt.input)
            if got != tt.want {
                t.Fatalf("FindMax(%v) = %d, want %d", tt.input, got, tt.want)
            }
        })
    }
}
```

### Coverage Analysis

```bash
# Generate coverage profile
go test -coverprofile=/tmp/coverage.out ./...

# View coverage by function
go tool cover -func=/tmp/coverage.out

# View coverage in browser
go tool cover -html=/tmp/coverage.out -o /tmp/coverage.html

# Check coverage threshold
go test -coverprofile=/tmp/coverage.out ./...
go tool cover -func=/tmp/coverage.out | grep total | awk '{print $3}'
```

## Testing Concurrent Code

### Testing with Race Detector

```go
func TestConcurrentAccess(t *testing.T) {
    counter := &SafeCounter{}
    const numGoroutines = 100
    const numIncrements = 1000

    var wg sync.WaitGroup
    wg.Add(numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < numIncrements; j++ {
                counter.Increment()
            }
        }()
    }

    wg.Wait()

    want := numGoroutines * numIncrements
    if got := counter.Value(); got != want {
        t.Fatalf("counter.Value() = %d, want %d", got, want)
    }
}

// Run with: go test -race
```

### Testing Channels

```go
func TestChannelCommunication(t *testing.T) {
    results := make(chan int, 10)
    done := make(chan bool)

    go func() {
        for i := 0; i < 10; i++ {
            results <- i
        }
        close(results)
        done <- true
    }()

    var received []int
    for val := range results {
        received = append(received, val)
    }

    <-done
    if len(received) != 10 {
        t.Fatalf("received len = %d, want 10", len(received))
    }
}
```

### Testing Timeouts

```go
func TestWithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    result := make(chan string, 1)
    go func() {
        time.Sleep(50 * time.Millisecond)
        result <- "success"
    }()

    select {
    case res := <-result:
        if res != "success" {
            t.Fatalf("result = %q, want %q", res, "success")
        }
    case <-ctx.Done():
        t.Fatal("test timed out")
    }
}
```

## Benchmarking

### Basic Benchmarks

```go
func BenchmarkMyFunction(b *testing.B) {
    input := generateTestData()

    b.ResetTimer() // Reset timer after setup

    for i := 0; i < b.N; i++ {
        MyFunction(input)
    }
}

// Run with: go test -bench=. -benchmem
```

### Table-Driven Benchmarks

```go
func BenchmarkStringOperations(b *testing.B) {
    benchmarks := []struct {
        name  string
        input []string
    }{
        {"small", generateStrings(10)},
        {"medium", generateStrings(100)},
        {"large", generateStrings(1000)},
    }

    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ConcatenateStrings(bm.input)
            }
        })
    }
}
```

## Test Organization

### File Naming

```
mypackage/
├── user.go
├── user_test.go       # Tests for user.go
├── auth.go
├── auth_test.go       # Tests for auth.go
└── testdata/          # Test fixtures
    └── sample.json
```

### Test Package Naming

```go
// Internal testing (access to private members)
package mypackage

import "testing"

func TestInternalFunction(t *testing.T) {
    // Can access private functions
}

// External testing (public API only)
package mypackage_test

import (
    "testing"
    "myapp/mypackage"
)

func TestPublicAPI(t *testing.T) {
    // Can only access exported functions
}
```

## Common Patterns

### Golden Files

```go
func TestOutputFormat(t *testing.T) {
    result := GenerateReport(testData)

    goldenFile := "testdata/report.golden"

    if *update {
        if err := os.WriteFile(goldenFile, []byte(result), 0644); err != nil {
            t.Fatalf("WriteFile() error = %v", err)
        }
    }

    expected, err := os.ReadFile(goldenFile)
    if err != nil {
        t.Fatalf("ReadFile() error = %v", err)
    }

    if string(expected) != result {
        t.Fatalf("result mismatch:
want %q
got  %q", string(expected), result)
    }
}

var update = flag.Bool("update", false, "update golden files")
```

### Testing HTTP Handlers

```go
func TestHTTPHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/users", nil)
    rec := httptest.NewRecorder()

    handler := NewUserHandler(mockStore)
    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
    }

    var users []User
    if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
        t.Fatalf("Unmarshal() error = %v", err)
    }
    if len(users) == 0 {
        t.Fatal("users is empty")
    }
}
```

## Summary

Effective Go testing requires:

- Clear test structure with explicit `got` / `want` messages
- Table-driven tests for multiple scenarios
- Proper use of helpers and cleanup
- Hand-written stubs/fakes (or package-adopted mocks) for dependencies
- Comprehensive coverage of edge cases and error paths
- Race detection for concurrent code
- Benchmarking for performance-critical code

Follow these patterns to achieve and maintain ≥ 80% test coverage. When companion Go Test rules (stem `go-test`) are present, prefer that companion for authoring conventions.
