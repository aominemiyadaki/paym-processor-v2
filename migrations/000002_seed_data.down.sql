DELETE FROM paym_processor.currencies WHERE iso_code IN ('RUB', 'USD', 'EUR');
DELETE FROM paym_processor.statuses WHERE name IN ('PENDING', 'SUCCESS', 'FAILED');