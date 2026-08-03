-- 自动更新更新时间
CREATE OR REPLACE FUNCTION smart_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.updated_at IS NOT DISTINCT FROM OLD.updated_at THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;