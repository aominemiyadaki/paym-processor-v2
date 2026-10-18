-- Сиды для валют
INSERT INTO paym_processor.currencies (name, iso_code)
VALUES 
    ('Russian Ruble', 'RUB'),
    ('US Dollar', 'USD'),
    ('Euro', 'EUR')
ON CONFLICT (iso_code) DO NOTHING;

-- Сиды для статусов платежей
INSERT INTO paym_processor.statuses (name)
VALUES 
    ('PENDING'),
    ('SUCCESS'),
    ('FAILED')
ON CONFLICT (name) DO NOTHING;