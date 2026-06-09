-- 000001_init_schema.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Tabel UMKM (pemilik usaha)
CREATE TABLE umkms (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    phone_number  VARCHAR(20)  NOT NULL UNIQUE, -- format: 628XXXXXXXXXX
    email         VARCHAR(255) UNIQUE,
    address       TEXT,
    business_type VARCHAR(100), -- contoh: 'kuliner', 'fashion', 'jasa', 'pertanian'
    password_hash VARCHAR(255) NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Tabel kategori transaksi
CREATE TABLE transaction_categories (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL, -- 'Penjualan', 'Pembelian Bahan', 'Gaji', dll
    type       VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    icon       VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabel transaksi
CREATE TABLE transactions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    umkm_id       UUID NOT NULL REFERENCES umkms(id),
    category_id   UUID REFERENCES transaction_categories(id),
    amount        BIGINT NOT NULL CHECK (amount > 0), -- dalam satuan Rupiah, BUKAN desimal
    type          VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    description   TEXT NOT NULL,           -- deskripsi asli dari AI
    raw_message   TEXT NOT NULL,           -- pesan WhatsApp asli dari UMKM
    transaction_date DATE NOT NULL,        -- tanggal transaksi (bisa berbeda dari created_at)
    source        VARCHAR(20) NOT NULL DEFAULT 'whatsapp' CHECK (source IN ('whatsapp', 'manual')),
    ai_confidence DECIMAL(3,2),            -- skor kepercayaan AI 0.00-1.00
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Tabel KUR score history
CREATE TABLE kur_scores (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    umkm_id             UUID NOT NULL REFERENCES umkms(id),
    score               INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    level               VARCHAR(20) NOT NULL CHECK (level IN ('rendah', 'sedang', 'baik', 'sangat_baik')),
    monthly_income_avg  BIGINT NOT NULL,   -- rata-rata pemasukan 3 bulan terakhir
    monthly_expense_avg BIGINT NOT NULL,
    profit_margin       DECIMAL(5,2),      -- persentase margin laba
    consistency_score   INTEGER,           -- konsistensi pencatatan (0-100)
    months_of_data      INTEGER NOT NULL,  -- berapa bulan data yang dianalisis
    recommendations     TEXT[],            -- array rekomendasi dalam Bahasa Indonesia
    calculated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_transactions_umkm_id ON transactions(umkm_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_deleted_at ON transactions(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_kur_scores_umkm_id ON kur_scores(umkm_id);
CREATE INDEX idx_umkms_phone ON umkms(phone_number);
