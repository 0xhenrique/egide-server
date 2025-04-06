CREATE TABLE IF NOT EXISTS sites_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    domain TEXT NOT NULL,
    protection_mode TEXT NOT NULL CHECK(protection_mode IN ('simple', 'hardened')),
    active BOOLEAN DEFAULT TRUE,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, domain)
);

INSERT INTO sites_new 
SELECT id, user_id, domain, protection_mode, active, FALSE as verified, created_at, updated_at FROM sites;

DROP TABLE sites;
ALTER TABLE sites_new RENAME TO sites;

CREATE INDEX idx_sites_user_id ON sites(user_id);
CREATE INDEX idx_sites_domain ON sites(domain);
