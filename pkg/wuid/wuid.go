package wuid

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	callbackwuid "github.com/edwingeng/wuid/callback/wuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	mu   sync.RWMutex
	Wuid *callbackwuid.WUID
	db   *sql.DB
)

const (
	wuidName    = "go-im-chat"
	h28Sequence = "chat.wuid_h_seq"
)

// ErrNotInitialized 当 WUID 未初始化时返回
var ErrNotInitialized = errors.New("wuid is not initialized or has been closed")

/**
 * Init 使用 PostgreSQL 作为 WUID 高位协调数据源并初始化生成器。
 * @param dsn PostgreSQL 连接串，例如 postgres://chat:123456@127.0.0.1:5432/chat?sslmode=disable
 */
func Init(dsn string) error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return nil
	}

	if dsn == "" {
		return errors.New("dsn cannot be empty")
	}

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)

	if err = conn.Ping(); err != nil {
		_ = conn.Close()
		return fmt.Errorf("ping postgres: %w", err)
	}

	db = conn
	Wuid = callbackwuid.NewWUID(wuidName, nil)

	// Init 已持有 mu.Lock，回调不能再抢 RLock，否则会死锁。
	if err = Wuid.LoadH28WithCallback(func() (int64, func(), error) {
		return nextH28FromConn(conn)
	}); err != nil {
		_ = db.Close()
		db = nil
		Wuid = nil
		return fmt.Errorf("load h28 from postgres: %w", err)
	}

	return nil
}

/**
 * loadH28FromPostgres 从 PostgreSQL 加载 H28 序列。
 */
func loadH28FromPostgres() (int64, func(), error) {
	mu.RLock()
	conn := db
	mu.RUnlock()

	if conn == nil {
		return 0, nil, errors.New("postgres connection is not initialized")
	}

	return nextH28FromConn(conn)
}

func nextH28FromConn(conn *sql.DB) (int64, func(), error) {
	var h28 int64
	query := "SELECT nextval($1)"
	if err := conn.QueryRow(query, h28Sequence).Scan(&h28); err != nil {
		return 0, nil, err
	}

	return h28, nil, nil
}

/**
 * Next 生成下一个全局唯一 ID（string 形式）。
 */
func Next() string {
	mu.RLock()
	w := Wuid
	mu.RUnlock()

	if w == nil {
		panic("wuid: Next() called before Init() or after Close()")
	}
	return strconv.FormatInt(w.Next(), 10)
}

/**
 * SafeNext 生成下一个 ID；未初始化时返回 ErrNotInitialized，不 panic。
 */
func SafeNext() (string, error) {
	mu.RLock()
	w := Wuid
	mu.RUnlock()

	if w == nil {
		return "", ErrNotInitialized
	}
	return strconv.FormatInt(w.Next(), 10), nil
}

/**
 * Close 关闭 PostgreSQL 连接，通常在服务退出时调用。
 */
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db == nil {
		return nil
	}

	err := db.Close()
	db = nil
	Wuid = nil
	return err
}
