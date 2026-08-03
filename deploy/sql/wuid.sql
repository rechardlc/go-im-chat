-- WUID 高位 (H28) 序列
-- 多实例共享同一 PostgreSQL，通过 nextval 分配 h28，保证分布式 ID 全局唯一。
-- 每个 h28 段本地可发约 550 亿个 ID，续租时才访问本序列（热路径不打 PG）。

CREATE SCHEMA IF NOT EXISTS chat;

CREATE SEQUENCE IF NOT EXISTS chat.wuid_h_seq
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    -- WUID monolithic 模式 VerifyH28 上限为 0x07FFFFFF（134217727），不是 2^28-1（268435455）。
    -- 若 MAXVALUE 超过 134217727，renew 时 LoadH28WithCallback 会报错。
    MAXVALUE 134217727
    -- 默认 NO CYCLE：序号耗尽后报错，避免 h28 复用带来 ID 冲突风险（实际上几乎不可能耗尽）。
    CACHE 1;

COMMENT ON SEQUENCE chat.wuid_h_seq IS 'WUID h28 序列；上限 0x07FFFFFF=134217727，勿用 2^28-1，勿 CYCLE';
