-- 000002_seed_categories.up.sql
INSERT INTO transaction_categories (id, name, type, icon) VALUES
    (gen_random_uuid(), 'Penjualan Produk', 'income', 'pi pi-shopping-bag'),
    (gen_random_uuid(), 'Penjualan Jasa', 'income', 'pi pi-briefcase'),
    (gen_random_uuid(), 'Pendapatan Lain', 'income', 'pi pi-plus-circle'),
    (gen_random_uuid(), 'Pembelian Bahan Baku', 'expense', 'pi pi-box'),
    (gen_random_uuid(), 'Gaji Karyawan', 'expense', 'pi pi-users'),
    (gen_random_uuid(), 'Sewa Tempat', 'expense', 'pi pi-home'),
    (gen_random_uuid(), 'Utilitas (Listrik/Air)', 'expense', 'pi pi-bolt'),
    (gen_random_uuid(), 'Transportasi', 'expense', 'pi pi-car'),
    (gen_random_uuid(), 'Pemasaran', 'expense', 'pi pi-megaphone'),
    (gen_random_uuid(), 'Peralatan', 'expense', 'pi pi-wrench'),
    (gen_random_uuid(), 'Pengeluaran Lain', 'expense', 'pi pi-minus-circle');
