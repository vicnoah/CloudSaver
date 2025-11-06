-- 初始化数据库脚本

-- 创建users表
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

CREATE INDEX idx_users_uuid ON users(uuid);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 创建global_settings表
CREATE TABLE IF NOT EXISTS global_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    http_proxy_host VARCHAR(255) NOT NULL DEFAULT '127.0.0.1',
    http_proxy_port INTEGER NOT NULL DEFAULT 7890,
    is_proxy_enabled BOOLEAN NOT NULL DEFAULT 0,
    common_user_code INTEGER DEFAULT 9527,
    admin_user_code INTEGER NOT NULL DEFAULT 230713,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认全局设置
INSERT OR IGNORE INTO global_settings (id, http_proxy_host, http_proxy_port, is_proxy_enabled, common_user_code, admin_user_code)
VALUES (1, '127.0.0.1', 7890, 0, 9527, 230713);

-- 创建user_settings表
CREATE TABLE IF NOT EXISTS user_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uuid VARCHAR(36) UNIQUE NOT NULL,
    cloud115_cookie TEXT,
    quark_cookie TEXT,
    cloud115_user_id VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE INDEX idx_user_settings_uuid ON user_settings(user_uuid);
