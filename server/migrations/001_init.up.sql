-- 创建管理员表
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建门店表
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    device_id UUID UNIQUE NOT NULL,
    device_fingerprint VARCHAR(64) NOT NULL,
    rustdesk_id VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'offline',
    approved BOOLEAN NOT NULL DEFAULT false,
    last_seen TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID REFERENCES stores(id) ON DELETE SET NULL,
    action_type VARCHAR(50) NOT NULL,
    operator VARCHAR(50) NOT NULL,
    details JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建命令日志表
CREATE TABLE IF NOT EXISTS command_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID REFERENCES stores(id) ON DELETE SET NULL,
    operator VARCHAR(50) NOT NULL,
    command TEXT NOT NULL,
    output TEXT,
    status VARCHAR(20) NOT NULL,
    executed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration_ms INT
);

-- 创建软件任务表
CREATE TABLE IF NOT EXISTS software_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    operator VARCHAR(50) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    install_args VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- 创建远程桌面会话表
CREATE TABLE IF NOT EXISTS desktop_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    operator VARCHAR(50) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_seconds INT,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

-- 创建资源历史表
CREATE TABLE IF NOT EXISTS resource_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    cpu_percent DECIMAL(5,2) NOT NULL,
    memory_percent DECIMAL(5,2) NOT NULL,
    memory_available_gb DECIMAL(10,2) NOT NULL,
    disk_data JSONB,
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建远程桌面排队表
CREATE TABLE IF NOT EXISTS desktop_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    operator VARCHAR(50) NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_stores_status ON stores(status);
CREATE INDEX IF NOT EXISTS idx_stores_approved ON stores(approved);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_store_id ON audit_logs(store_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action_type ON audit_logs(action_type);
CREATE INDEX IF NOT EXISTS idx_resource_history_store_time ON resource_history(store_id, recorded_at);
CREATE INDEX IF NOT EXISTS idx_command_logs_store_id ON command_logs(store_id);
CREATE INDEX IF NOT EXISTS idx_desktop_queue_store ON desktop_queue(store_id, position);

-- 插入默认管理员 (密码: admin123)
INSERT INTO admins (username, password_hash)
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.iW8jTpTqMxVwK9FzFu')
ON CONFLICT (username) DO NOTHING;