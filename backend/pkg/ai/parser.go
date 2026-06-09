package ai

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// ParsedTransaction adalah hasil parsing AI dari pesan WhatsApp.
type ParsedTransaction struct {
	IsTransaction   bool    `json:"is_transaction"`
	Type            string  `json:"type"`            // income | expense | ""
	Amount          int64   `json:"amount"`          // dalam Rupiah
	Description     string  `json:"description"`     // deskripsi singkat Bahasa Indonesia
	CategoryHint    string  `json:"category_hint"`   // nama kategori
	TransactionDate string  `json:"transaction_date"` // YYYY-MM-DD
	Confidence      float64 `json:"confidence"`      // 0.0 - 1.0
	ReplyMessage    string  `json:"reply_message"`   // pesan balasan untuk user
}

// Parser menggunakan AI Client untuk parsing pesan WhatsApp menjadi transaksi terstruktur.
type Parser struct {
	client *Client
}

func NewParser(client *Client) *Parser {
	return &Parser{client: client}
}

// Parse mengirim pesan ke DeepSeek API dan mengembalikan ParsedTransaction.
func (p *Parser) Parse(rawMessage string) (*ParsedTransaction, error) {
	today := time.Now().Format("2006-01-02")

	userPrompt := fmt.Sprintf(
		"Tanggal hari ini: %s\nPesan WhatsApp: \"%s\"",
		today, rawMessage,
	)

	response, err := p.client.Chat(systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("ai.Parser.Parse: %w", err)
	}

	// Bersihkan response dari markdown code fences
	response = cleanJSON(response)

	var parsed ParsedTransaction
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		slog.Error("gagal unmarshal response AI", "response", response, "error", err)
		return nil, fmt.Errorf("ai.Parser.Parse: unmarshal: %w", err)
	}

	return &parsed, nil
}

// cleanJSON menghapus markdown code fences (```json ... ```) dari response AI.
func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

// systemPrompt adalah prompt sistem untuk DeepSeek API.
// WAJIB digunakan persis seperti yang didefinisikan di AGENTS.md.
const systemPrompt = `Kamu adalah asisten keuangan UMKM Indonesia. Tugasmu adalah mengekstrak informasi transaksi dari pesan WhatsApp pelaku UMKM.

Balas HANYA dengan JSON valid, tidak ada teks lain di luar JSON.

Format JSON yang harus dikembalikan:
{
  "is_transaction": boolean,
  "type": "income" | "expense" | null,
  "amount": integer (dalam Rupiah, tanpa desimal) | null,
  "description": "deskripsi singkat dalam Bahasa Indonesia" | null,
  "category_hint": "nama kategori yang paling cocok dari daftar" | null,
  "transaction_date": "YYYY-MM-DD" | null,
  "confidence": float (0.0 - 1.0),
  "reply_message": "pesan balasan singkat dan ramah dalam Bahasa Indonesia untuk dikirim ke pengguna"
}

Aturan parsing:
- Jika pesan bukan transaksi (salam, pertanyaan, dll): is_transaction=false, reply_message berisi respons yang membantu
- Tanggal: jika tidak disebutkan, gunakan hari ini
- Kata "jual/dapat/masuk/terima" → income
- Kata "beli/bayar/keluar/bayarin/belanja" → expense
- Format angka Indonesia: "5rb" = 5000, "2jt" = 2000000, "500k" = 500000
- Jika ambigu, confidence rendah (< 0.6) dan minta klarifikasi di reply_message

Daftar kategori: Penjualan Produk, Penjualan Jasa, Pendapatan Lain, Pembelian Bahan Baku, Gaji Karyawan, Sewa Tempat, Utilitas (Listrik/Air), Transportasi, Pemasaran, Peralatan, Pengeluaran Lain`
