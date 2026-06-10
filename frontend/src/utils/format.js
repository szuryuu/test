/**
 * Format angka ke Rupiah — "Rp1.500.000"
 */
export const formatRupiah = (amount) => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(amount);
};

/**
 * Format tanggal ISO ke "8 Juni 2026"
 */
export const formatDate = (dateStr) => {
  if (!dateStr) return "";
  const date = new Date(dateStr);
  return date.toLocaleDateString("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
};

/**
 * Parse teks WhatsApp ke HTML untuk tampilan web.
 * - Escape HTML entities agar aman untuk v-html
 * - Hapus emoji berlebih
 * - Konversi *tebal* → <strong>, _miring_ → <em>, ~coret~ → <del>
 * - Konversi \n → <br>
 *
 * Data asli (report.report_text) tidak diubah — hanya output rendering.
 */
export const parseWhatsAppText = (text) => {
  if (!text) return "";

  let out = text;

  // 1. Escape HTML entities (keamanan untuk v-html)
  out = out
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

  // 2. Hapus emoji (kecuali simbol numerik/fungsi: 1-9, #, *, dsb tidak termasuk)
  //    Regex ini menangkap sebagian besar emoji Unicode
  out = out.replace(/[\u{1F600}-\u{1F64F}\u{1F300}-\u{1F5FF}\u{1F680}-\u{1F6FF}\u{1F1E0}-\u{1F1FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}\u{FE00}-\u{FE0F}\u{1F900}-\u{1F9FF}\u{1FA00}-\u{1FA6F}\u{1FA70}-\u{1FAFF}]/gu, "");

  // 3. Konversi WhatsApp markdown → HTML
  //    *tebal* → <strong> (WA pakai satu *, bukan ** seperti Markdown standar)
  out = out.replace(/\*(.+?)\*/g, "<strong>$1</strong>");
  //    _miring_ → <em>
  out = out.replace(/_(.+?)_/g, "<em>$1</em>");
  //    ~coret~ → <del>
  out = out.replace(/~(.+?)~/g, "<del>$1</del>");

  // 4. Konversi newline → <br> (white-space: pre-wrap di CSS sudah cukup,
  //    tapi dengan v-html lebih aman eksplisit)
  out = out.replace(/\n/g, "<br>");

  return out;
};
