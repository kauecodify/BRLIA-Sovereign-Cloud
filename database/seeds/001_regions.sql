-- Regioes soberanas iniciais
INSERT INTO regions (id, name, description, location, country, status, is_default)
VALUES
    ('00000000-0000-0000-0001-000000000001', 'br-sul-1', 'Regiao Sul - Porto Alegre', 'Porto Alegre, RS', 'BR', 'active', true),
    ('00000000-0000-0000-0001-000000000002', 'br-suldeste-1', 'Regiao Sudeste - Sao Paulo', 'Sao Paulo, SP', 'BR', 'active', false),
    ('00000000-0000-0000-0001-000000000003', 'br-nordeste-1', 'Regiao Nordeste - Fortaleza', 'Fortaleza, CE', 'BR', 'active', false)
ON CONFLICT (name) DO NOTHING;
