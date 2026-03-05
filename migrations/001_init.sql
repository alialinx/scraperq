-- Kullanıcılar
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       api_key VARCHAR(64) UNIQUE NOT NULL,
                       is_active BOOLEAN DEFAULT true,
                       daily_limit INT DEFAULT 100,
                       monthly_limit INT DEFAULT 2000,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW()
);

-- Refresh token (güvenli auth)
CREATE TABLE refresh_tokens (
                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                token VARCHAR(255) UNIQUE NOT NULL,
                                expires_at TIMESTAMP NOT NULL,
                                is_revoked BOOLEAN DEFAULT false,
                                created_at TIMESTAMP DEFAULT NOW()
);

-- Scrape işleri
CREATE TABLE jobs (
                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      user_id UUID NOT NULL REFERENCES users(id),
                      url TEXT NOT NULL,
                      status VARCHAR(20) DEFAULT 'pending',
                      retry_count INT DEFAULT 0,
                      max_retries INT DEFAULT 3,
                      error TEXT,
                      created_at TIMESTAMP DEFAULT NOW(),
                      updated_at TIMESTAMP DEFAULT NOW()
);

-- Scrape sonuçları
CREATE TABLE results (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         job_id UUID UNIQUE NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
                         status_code INT,
                         body_size INT,
                         title TEXT,
                         headers JSONB,
                         scraped_at TIMESTAMP DEFAULT NOW()
);

-- Rate limiting
CREATE TABLE rate_limits (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             request_count INT DEFAULT 0,
                             window_start TIMESTAMP DEFAULT NOW(),
                             window_type VARCHAR(10) NOT NULL,
                             UNIQUE(user_id, window_type)
);

-- Audit log
CREATE TABLE audit_logs (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            user_id UUID REFERENCES users(id),
                            action VARCHAR(50) NOT NULL,
                            resource VARCHAR(50),
                            resource_id UUID,
                            ip_address INET,
                            details JSONB,
                            created_at TIMESTAMP DEFAULT NOW()
);

-- Indexler (performans için)
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);
CREATE INDEX idx_results_job_id ON results(job_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_rate_limits_user_window ON rate_limits(user_id, window_type);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Updated_at otomatik güncelleme
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER jobs_updated_at BEFORE UPDATE ON jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();