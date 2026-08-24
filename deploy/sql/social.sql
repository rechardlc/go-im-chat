-- 好友关系表
CREATE TABLE IF NOT EXISTS friends (
    id serial PRIMARY KEY,
    user_id varchar(64) NOT NULL,
    friend_uid varchar(64) NOT NULL,
    remark varchar(255) DEFAULT NULL,
    add_source smallint DEFAULT NULL,
    created_at timestamptz DEFAULT now()
);
-- 已有库：CREATE TABLE IF NOT EXISTS 不会改列默认值
ALTER TABLE friends ALTER COLUMN created_at SET DEFAULT now();

COMMENT ON TABLE friends IS '好友关系';
COMMENT ON COLUMN friends.id IS '主键';
COMMENT ON COLUMN friends.user_id IS '用户ID';
COMMENT ON COLUMN friends.friend_uid IS '好友用户ID';
COMMENT ON COLUMN friends.remark IS '备注';
COMMENT ON COLUMN friends.add_source IS '添加来源';
COMMENT ON COLUMN friends.created_at IS '创建时间';

-- 好友申请表
CREATE TABLE IF NOT EXISTS friend_requests (
    id serial PRIMARY KEY,
    user_id varchar(64) NOT NULL,
    req_uid varchar(64) NOT NULL,
    req_msg varchar(255) DEFAULT NULL,
    req_time timestamptz NOT NULL,
    handle_result smallint DEFAULT NULL,
    handle_msg varchar(255) DEFAULT NULL,
    handled_at timestamptz DEFAULT NULL
);

COMMENT ON TABLE friend_requests IS '好友申请';
COMMENT ON COLUMN friend_requests.id IS '主键';
COMMENT ON COLUMN friend_requests.user_id IS '被申请用户ID';
COMMENT ON COLUMN friend_requests.req_uid IS '申请人用户ID';
COMMENT ON COLUMN friend_requests.req_msg IS '申请留言';
COMMENT ON COLUMN friend_requests.req_time IS '申请时间';
COMMENT ON COLUMN friend_requests.handle_result IS '处理结果';
COMMENT ON COLUMN friend_requests.handle_msg IS '处理留言';
COMMENT ON COLUMN friend_requests.handled_at IS '处理时间';

-- groups 在 PostgreSQL 15+ 是保留字（GROUP BY GROUPS），必须加双引号
CREATE TABLE IF NOT EXISTS "groups" (
    id varchar(24) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    icon varchar(255) NOT NULL,
    status smallint DEFAULT NULL,
    creator_uid varchar(64) NOT NULL,
    group_type integer NOT NULL,
    is_verify boolean NOT NULL,
    notification varchar(255) DEFAULT NULL,
    notification_uid varchar(64) DEFAULT NULL,
    created_at timestamptz DEFAULT NULL,
    updated_at timestamptz DEFAULT NULL
);

COMMENT ON TABLE "groups" IS '群';
COMMENT ON COLUMN "groups".id IS '群ID';
COMMENT ON COLUMN "groups".name IS '群名称';
COMMENT ON COLUMN "groups".icon IS '群头像';
COMMENT ON COLUMN "groups".status IS '群状态';
COMMENT ON COLUMN "groups".creator_uid IS '创建者用户ID';
COMMENT ON COLUMN "groups".group_type IS '群类型';
COMMENT ON COLUMN "groups".is_verify IS '入群是否需要审核';
COMMENT ON COLUMN "groups".notification IS '群公告';
COMMENT ON COLUMN "groups".notification_uid IS '公告发布者用户ID';
COMMENT ON COLUMN "groups".created_at IS '创建时间';
COMMENT ON COLUMN "groups".updated_at IS '更新时间';

DROP TRIGGER IF EXISTS trg_groups_updated_at ON "groups";
CREATE TRIGGER trg_groups_updated_at
BEFORE UPDATE ON "groups"
FOR EACH ROW EXECUTE FUNCTION smart_set_updated_at();

-- 群成员表
CREATE TABLE IF NOT EXISTS group_members (
    id serial PRIMARY KEY,
    group_id varchar(64) NOT NULL,
    user_id varchar(64) NOT NULL,
    role_level smallint NOT NULL,
    join_time timestamptz DEFAULT NULL,
    join_source smallint DEFAULT NULL,
    inviter_uid varchar(64) DEFAULT NULL,
    operator_uid varchar(64) DEFAULT NULL
);

COMMENT ON TABLE group_members IS '群成员';
COMMENT ON COLUMN group_members.id IS '主键';
COMMENT ON COLUMN group_members.group_id IS '群ID';
COMMENT ON COLUMN group_members.user_id IS '成员用户ID';
COMMENT ON COLUMN group_members.role_level IS '角色等级';
COMMENT ON COLUMN group_members.join_time IS '入群时间';
COMMENT ON COLUMN group_members.join_source IS '入群来源';
COMMENT ON COLUMN group_members.inviter_uid IS '邀请人用户ID';
COMMENT ON COLUMN group_members.operator_uid IS '操作人用户ID';

-- 入群申请表
CREATE TABLE IF NOT EXISTS group_requests (
    id serial PRIMARY KEY,
    req_id varchar(64) NOT NULL,
    group_id varchar(64) NOT NULL,
    req_msg varchar(255) DEFAULT NULL,
    req_time timestamptz DEFAULT NULL,
    join_source smallint DEFAULT NULL,
    inviter_user_id varchar(64) DEFAULT NULL,
    handle_user_id varchar(64) DEFAULT NULL,
    handle_time timestamptz DEFAULT NULL,
    handle_result smallint DEFAULT NULL
);

COMMENT ON TABLE group_requests IS '入群申请';
COMMENT ON COLUMN group_requests.id IS '主键';
COMMENT ON COLUMN group_requests.req_id IS '申请人用户ID';
COMMENT ON COLUMN group_requests.group_id IS '群ID';
COMMENT ON COLUMN group_requests.req_msg IS '申请留言';
COMMENT ON COLUMN group_requests.req_time IS '申请时间';
COMMENT ON COLUMN group_requests.join_source IS '入群来源';
COMMENT ON COLUMN group_requests.inviter_user_id IS '邀请人用户ID';
COMMENT ON COLUMN group_requests.handle_user_id IS '处理人用户ID';
COMMENT ON COLUMN group_requests.handle_time IS '处理时间';
COMMENT ON COLUMN group_requests.handle_result IS '处理结果';
