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

	MsgHelp = `🤖 *Pusat Bantuan KasirAI*

Anda bisa mencatat transaksi cukup dengan *chat* biasa (contoh: "Terjual 5 porsi bakso 100rb").

*Perintah Khusus:*
📊 *hari ini* - Ringkasan keuangan hari ini
📈 *minggu ini* - Ringkasan keuangan minggu ini
📑 *laporan* - Buat laporan bulanan lengkap
⭐ *skor* - Cek skor kesiapan KUR Anda
⏪ *batalkan* - Hapus transaksi yang paling terakhir dicatat (gunakan jika ada salah ketik/revisi)

Ada yang bisa saya bantu catat sekarang?`

	MsgAIParseFailure = "Maaf, saya tidak bisa memahami transaksi tersebut. Coba tulis lebih jelas, contoh: 'Terjual 1 porsi bakso 20rb'"

	MsgTransaksiTercatat = "✅ Tercatat! %s Rp%s. %s"
)
