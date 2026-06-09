package fonnte

// Template pesan WhatsApp standar — WAJIB digunakan di seluruh aplikasi.
// Simpan sebagai konstanta untuk konsistensi dan kemudahan maintenance.

const (
	MsgWelcome = `Selamat datang di KasirAI! 👋

Saya adalah asisten keuangan UMKM Anda.

Cara pakai:
📥 Catat pemasukan: "jual nasi goreng 3 porsi 45rb"
📤 Catat pengeluaran: "beli bahan 150rb"
📊 Lihat ringkasan: ketik *laporan*
⭐ Cek skor KUR: ketik *skor*
❓ Bantuan: ketik *bantuan*`

	MsgHelp = `Perintah yang tersedia:
• *laporan* - Ringkasan keuangan bulan ini
• *hari ini* - Transaksi hari ini
• *minggu ini* - Transaksi minggu ini
• *skor* - Cek kesiapan KUR Anda
• *hapus terakhir* - Hapus transaksi terakhir

Untuk mencatat transaksi, cukup tulis seperti biasa:
"jual 5 kaos 250rb" atau "bayar listrik 200rb"`

	MsgAIParseFailure = "Maaf, saya tidak bisa memahami transaksi tersebut. Coba tulis lebih jelas, contoh: 'jual nasi 30rb'"

	MsgTransaksiTercatat = "✅ Tercatat! %s Rp%s. %s"
)
