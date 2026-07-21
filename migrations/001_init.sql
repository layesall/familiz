-- Table des utilisateurs (admin pour la V1)
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'admin',
    member_id INTEGER,  -- sera NULL pour l'admin (car il n'est pas encore membre)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table des membres (les gens de la famille)
CREATE TABLE IF NOT EXISTS members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    birth_date TEXT NOT NULL, -- format YYYY-MM-DD
    marital_status TEXT NOT NULL CHECK(marital_status IN ('single', 'married', 'minor')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table des transactions (paiements)
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id INTEGER NOT NULL,
    month INTEGER NOT NULL CHECK(month BETWEEN 1 AND 12),
    year INTEGER NOT NULL,
    amount REAL NOT NULL CHECK(amount > 0),
    paid_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    note TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);

-- Table des événements (mariage, baptême)
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id INTEGER NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('wedding', 'baptism')),
    amount_received REAL NOT NULL CHECK(amount_received >= 0),
    event_date TEXT NOT NULL, -- format YYYY-MM-DD
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE
);

-- Table des montants de cotisation par statut (une seule ligne)
CREATE TABLE IF NOT EXISTS contribution_settings (
    id INTEGER PRIMARY KEY CHECK (id = 1), -- On force une ligne unique
    amount_single REAL NOT NULL DEFAULT 10.0,
    amount_married REAL NOT NULL DEFAULT 15.0,
    amount_minor REAL NOT NULL DEFAULT 5.0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- On insère la ligne par défaut si elle n'existe pas
INSERT OR IGNORE INTO contribution_settings (id, amount_single, amount_married, amount_minor)
VALUES (1, 10.0, 15.0, 5.0);

-- Table des montants par défaut pour les événements
CREATE TABLE IF NOT EXISTS event_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL UNIQUE CHECK(event_type IN ('wedding', 'baptism')),
    default_amount REAL NOT NULL DEFAULT 200.0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insérer les valeurs par défaut
INSERT OR IGNORE INTO event_settings (event_type, default_amount) VALUES ('wedding', 200.0);
INSERT OR IGNORE INTO event_settings (event_type, default_amount) VALUES ('baptism', 200.0);

ALTER TABLE transactions ADD COLUMN is_archived BOOLEAN DEFAULT 0;
ALTER TABLE events ADD COLUMN is_archived BOOLEAN DEFAULT 0;
ALTER TABLE contribution_settings ADD COLUMN current_year INTEGER DEFAULT 2026;