-- Organizacao e projeto de demonstracao
INSERT INTO organizations (id, name, description, status)
VALUES ('00000000-0000-0000-0003-000000000001', 'brlia-demo', 'Organizacao de demonstracao', 'active')
ON CONFLICT (name) DO NOTHING;

INSERT INTO projects (id, organization_id, name, description, region_id, status)
VALUES (
    '00000000-0000-0000-0003-000000000002',
    '00000000-0000-0000-0003-000000000001',
    'default',
    'Projeto padrao da organizacao de demonstracao',
    '00000000-0000-0000-0001-000000000001',
    'active'
)
ON CONFLICT DO NOTHING;
