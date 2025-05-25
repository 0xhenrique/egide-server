CREATE TABLE health_checks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    response_time_ms INTEGER NOT NULL,
    status_code INTEGER,
    success BOOLEAN NOT NULL,
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_health_checks_timestamp ON health_checks(timestamp);
CREATE INDEX idx_health_checks_success ON health_checks(success);
