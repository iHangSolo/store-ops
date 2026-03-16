-- 删除索引
DROP INDEX IF EXISTS idx_desktop_queue_store;
DROP INDEX IF EXISTS idx_command_logs_store_id;
DROP INDEX IF EXISTS idx_resource_history_store_time;
DROP INDEX IF EXISTS idx_audit_logs_action_type;
DROP INDEX IF EXISTS idx_audit_logs_store_id;
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_stores_approved;
DROP INDEX IF EXISTS idx_stores_status;

-- 删除表
DROP TABLE IF EXISTS desktop_queue;
DROP TABLE IF EXISTS resource_history;
DROP TABLE IF EXISTS desktop_sessions;
DROP TABLE IF EXISTS software_tasks;
DROP TABLE IF EXISTS command_logs;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS admins;