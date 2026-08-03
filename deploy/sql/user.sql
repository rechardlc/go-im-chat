CREATE TABLE IF NOT EXISTS users (
    id varchar(24) NOT NULL PRIMARY KEY,
    avatar varchar(255) NOT NULL DEFAULT '',
    nickname varchar(20) NOT NULL DEFAULT '',
    phone varchar(255) NOT NULL DEFAULT '',
    password varchar(60) NOT NULL DEFAULT '',
    -- smallint: 相比int数字范围更小，节省空间
    status smallint NOT NULL DEFAULT 0,
    sex smallint NOT NULL DEFAULT 0,
    -- timestamptz: 存储时区敏感的日期时间，自动处理时区转换
    -- CURRENT_TIMESTAMP: 当前时间，自动处理时区转换
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
)   
COMMENT ON COLUMN users.id IS '用户ID';
COMMENT ON COLUMN users.avatar IS '头像';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.phone IS '手机号';
COMMENT ON COLUMN users.password IS '密码 bcrypt 哈希';
COMMENT ON COLUMN users.status IS '状态 0: 正常 1: 禁用';
COMMENT ON COLUMN users.sex IS '性别 0: 未知 1: 男 2: 女';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';

-- 索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users (phone);
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);

-- 自动更新更新时间
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION smart_set_updated_at();

-- drop trigger exists before update users for each row execte function