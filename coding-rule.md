# Golang コーディング規約

## 1. 基本原則

### 1.1 フォーマットとスタイル
- `gofmt`の使用は必須
  - コミット前に必ず`gofmt -s -w .`を実行
  - CIパイプラインで`gofmt`チェックを実施
- `golint`による静的解析の実施
  - すべての警告を解決すること
  - 特定の警告を無視する場合は理由をコメントで記述

### 1.2 コードスタイル参照
- [Effective Go](https://go.dev/doc/effective_go)に準拠
- [Code Review Comments](https://go.dev/wiki/CodeReviewComments)のガイドラインに従う
- 以下の追加ツールの使用を推奨：
  - `golangci-lint`
  - `staticcheck`

## 2. コード構造

### 2.1 パッケージ設計
```go
// Good
package user

// パッケージ名は簡潔で意味を表す
// main以外のパッケージ名は単数形
// ディレクトリ名とパッケージ名は一致させる
```

### 2.2 インターフェース定義
```go
// Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

// インターフェースは使用される場所で定義
// 小さなインターフェースを優先
// 必要なメソッドのみを定義
```

### 2.3 構造体定義
```go
type User struct {
    // 大文字から始まるフィールドは公開
    ID        string    `json:"id" validate:"required"`
    Email     string    `json:"email" validate:"email"`
    
    // 小文字から始まるフィールドは非公開
    password  string    // パスワードハッシュ
    lastLogin time.Time // 最終ログイン時刻
}
```

## 3. エラーハンドリング

### 3.1 エラー定義
```go
// カスタムエラー型の定義
type ErrorCode int

const (
    ErrNotFound ErrorCode = iota + 1
    ErrInvalidInput
    ErrDatabase
)

type AppError struct {
    Code    ErrorCode
    Message string
    Err     error
    Context map[string]interface{}
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}
```

### 3.2 エラーハンドリングパターン
```go
func (s *service) CreateUser(ctx context.Context, user *User) error {
    // バリデーション
    if err := user.Validate(); err != nil {
        return &AppError{
            Code:    ErrInvalidInput,
            Message: "invalid user data",
            Err:     err,
            Context: map[string]interface{}{
                "user_email": user.Email,
                "validation_errors": err.Error(),
            },
        }
    }

    // データベース操作
    if err := s.repo.Create(ctx, user); err != nil {
        // データベースエラーのラッピング
        return &AppError{
            Code:    ErrDatabase,
            Message: "failed to create user",
            Err:     err,
            Context: map[string]interface{}{
                "user_id": user.ID,
            },
        }
    }

    return nil
}
```

### 3.3 エラーチェック
```go
// Good
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Bad
if err != nil {
    return err // コンテキスト情報が失われる
}
```

## 4. ドキュメンテーション

### 4.1 パッケージドキュメント
```go
// Package user implements user management functionality.
// It provides operations for user CRUD, authentication, and authorization.
//
// Usage:
//
//     svc := user.NewService(repo, logger)
//     user, err := svc.CreateUser(ctx, &CreateUserRequest{...})
package user
```

### 4.2 型とメソッドのドキュメント
```go
// UserService provides user management operations.
// It handles all user-related business logic and ensures data consistency.
type UserService interface {
    // CreateUser creates a new user in the system.
    // It returns ErrInvalidInput if the user data is invalid,
    // or ErrDatabase if there was a problem with the database operation.
    CreateUser(ctx context.Context, user *User) error

    // GetUser retrieves a user by their ID.
    // It returns ErrNotFound if the user does not exist.
    GetUser(ctx context.Context, id string) (*User, error)
}
```

## 5. テスト

### 5.1 テストファイル構成
```go
func TestUserService_CreateUser(t *testing.T) {
    // テストケースの定義
    tests := []struct {
        name    string
        user    *User
        mock    func(*MockRepository)
        wantErr error
    }{
        {
            name: "valid user",
            user: &User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            mock: func(r *MockRepository) {
                r.EXPECT().
                    Create(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            wantErr: nil,
        },
        // 他のテストケース...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テストの実装
        })
    }
}
```

### 5.2 テストヘルパー
```go
// テストヘルパー関数の命名は "test" プレフィックスを使用
func testSetup(t *testing.T) (*UserService, *MockRepository) {
    ctrl := gomock.NewController(t)
    repo := NewMockRepository(ctrl)
    svc := NewUserService(repo)
    return svc, repo
}
```

## 6. パフォーマンスとリソース管理

### 6.1 メモリ使用
```go
// Good: 事前にキャパシティを確保
users := make([]User, 0, len(ids))

// Bad: 動的な拡張が必要
var users []User
```

### 6.2 ゴルーチン管理
```go
func ProcessItems(ctx context.Context, items []Item) error {
    eg, ctx := errgroup.WithContext(ctx)
    
    for _, item := range items {
        item := item // ループ変数のキャプチャ
        eg.Go(func() error {
            return processItem(ctx, item)
        })
    }

    return eg.Wait()
}
```

## 7. レビューチェックリスト

### 7.1 コード品質
- [ ] `gofmt -s -w .`が適用されている
- [ ] `golint`の警告がない
- [ ] `golangci-lint run`がパスする
- [ ] Effective Goのガイドラインに従っている
- [ ] Code Review Commentsの推奨事項に従っている

### 7.2 エラー処理
- [ ] すべてのエラーが適切にハンドリングされている
- [ ] エラーに十分なコンテキスト情報が含まれている
- [ ] エラーチェーンが適切に維持されている
- [ ] カスタムエラー型が適切に使用されている

### 7.3 テスト
- [ ] テストカバレッジが80%以上
- [ ] エッジケースがテストされている
- [ ] モックが適切に使用されている
- [ ] テストが読みやすく保守可能