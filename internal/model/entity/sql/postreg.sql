-- PostgreSQL 创建语句
-- 注意：PostgreSQL 对枚举类型的处理与 MySQL 不同

-- 1. 创建用户信息表
CREATE TABLE IF NOT EXISTS user_info (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_name TEXT NOT NULL,
    password TEXT,
    role TEXT NOT NULL,
    is_locked INTEGER DEFAULT 0,
    extra_info TEXT,
    CONSTRAINT uk_user_info_user_name UNIQUE (user_name)
);

CREATE INDEX IF NOT EXISTS idx_user_info_deleted_at ON user_info(deleted_at);

-- 2. 创建用户操作记录表
CREATE TABLE IF NOT EXISTS user_operation (
    id BIGSERIAL PRIMARY KEY,
    user_name TEXT NOT NULL,
    operation TEXT NOT NULL,
    module TEXT NOT NULL,
    target_id BIGINT,
    target_name TEXT,
    ip TEXT,
    user_agent TEXT,
    extra_info TEXT NOT NULL,
    create_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. 创建点云地图表
CREATE TABLE IF NOT EXISTS pcd_file (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name TEXT NOT NULL,
    area TEXT NOT NULL,
    path TEXT NOT NULL,
    user_name TEXT NOT NULL,
    size INTEGER,
    minio_path TEXT,
    extra_info TEXT
);

CREATE INDEX IF NOT EXISTS idx_pcd_file_deleted_at ON pcd_file(deleted_at);

-- 4. 创建语义地图表
CREATE TABLE IF NOT EXISTS semantic_map (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    pcd_file_id BIGINT NOT NULL,
    user_name TEXT NOT NULL,
    semantic_info TEXT,
    extra_info TEXT,
    CONSTRAINT fk_semantic_map_pcd_file 
        FOREIGN KEY (pcd_file_id) REFERENCES pcd_file(id)
);

CREATE INDEX IF NOT EXISTS idx_semantic_map_deleted_at ON semantic_map(deleted_at);
CREATE INDEX IF NOT EXISTS idx_semantic_map_pcd_file_id ON semantic_map(pcd_file_id);

-- 5. 创建任务编排表
CREATE TABLE IF NOT EXISTS task (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    semantic_map_id BIGINT NOT NULL,
    user_name TEXT NOT NULL,
    task_info TEXT,
    status TEXT DEFAULT 'pending',
    extra_info TEXT,
    CONSTRAINT fk_task_semantic_map 
        FOREIGN KEY (semantic_map_id) REFERENCES semantic_map(id)
);

CREATE INDEX IF NOT EXISTS idx_task_deleted_at ON task(deleted_at);
CREATE INDEX IF NOT EXISTS idx_task_semantic_map_id ON task(semantic_map_id);

-- 6. 创建设备表
CREATE TABLE IF NOT EXISTS device (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    type TEXT NOT NULL,
    company TEXT NOT NULL,
    ip TEXT,
    port INTEGER,
    user_name TEXT,
    password TEXT,
    status TEXT DEFAULT 'offline',
    extra_info TEXT
);

CREATE INDEX IF NOT EXISTS idx_device_deleted_at ON device(deleted_at);
