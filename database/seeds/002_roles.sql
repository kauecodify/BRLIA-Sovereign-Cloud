-- Roles do sistema
INSERT INTO roles (id, name, description, permissions, is_system)
VALUES
    ('00000000-0000-0000-0002-000000000001', 'admin', 'Administrador da plataforma', '["*"]', true),
    ('00000000-0000-0000-0002-000000000002', 'operator', 'Operador de infraestrutura', '["compute:*", "storage:*", "network:*"]', true),
    ('00000000-0000-0000-0002-000000000003', 'developer', 'Desenvolvedor de aplicacoes', '["compute:create", "storage:use", "network:use"]', true),
    ('00000000-0000-0000-0002-000000000004', 'viewer', 'Leitura de estado', '["*:read"]', true)
ON CONFLICT (name) DO NOTHING;
