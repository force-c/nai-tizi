-- +goose Up
DELETE FROM m_role_api_permission
WHERE permission_id IN (SELECT id FROM s_api_permission WHERE code = 'role.permission');

DELETE FROM m_user_api_permission
WHERE permission_id IN (SELECT id FROM s_api_permission WHERE code = 'role.permission');

DELETE FROM s_api_permission WHERE code = 'role.permission';

DROP TABLE IF EXISTS casbin_rule;

