// Klangauge koç motoru — beklenen (heceleme) ile söyleneni (analiz) karşılaştırır
// ve öğrenci dostu yorum üretir. Metinler UI diline göre üç dilde (tr/de/en).
// expected: core.SyllabifyWord JSON'u, analysis: core.Analyze JSON'u.

const langName = {
  tr: { tr: 'Türkçe', de: 'Almanca', en: 'İngilizce' },
  de: { tr: 'Türkisch', de: 'Deutsch', en: 'Englisch' },
  en: { tr: 'Turkish', de: 'German', en: 'English' },
}

const msg = {
  tr: {
    typeWord: 'Önce yukarıdan bir kelime yaz — beklenen heceler çizilsin.',
    noData: 'Kayıt yok — kelimeyi söyleyip koçtan yorum al.',
    countOk: 'Hece sayısı doğru: {d}/{e}.',
    countMore: 'Beklenen {e} heceydi, {d} algılandı. Vurgusuz bir hece eklemiş olabilirsin — heceleri hafifçe kısalt.',
    countLess: 'Beklenen {e} heceydi, {d} algılandı. Hece yutuyor olabilirsin — son heceyi net söyle.',
    countZero: 'Konuşma algılanamadı — mikrofona yaklaşıp kelimeyi net söylemeyi dene.',
    stressOk: 'Vurgu doğru hecede: {i}. hece.',
    stressDiff: 'Beklenen vurgu {e}. hecedeydi, sen {d}. heceyi vurguladın.',
    stressApprox: 'Not: {lang} vurgu kuralı sözlüksüz yaklaşıktır — sözlüğe bakmak doğru sonucu verir.',
    stressUnknown: 'Söylenen vurgu hecesi tespit edilemedi.',
  },
  de: {
    typeWord: 'Schreib zuerst oben ein Wort — die erwarteten Silben werden angezeigt.',
    noData: 'Keine Aufnahme — sprich das Wort und hol dir den Kommentar des Coaches.',
    countOk: 'Silbenzahl richtig: {d}/{e}.',
    countMore: 'Erwartet waren {e} Silben, erkannt {d}. Du hast evtl. eine unbetonte Silbe hinzugefügt — kürze die Silben etwas.',
    countLess: 'Erwartet waren {e} Silben, erkannt {d}. Du verschluckst evtl. Silben — sprich die letzte Silbe deutlich.',
    countZero: 'Kein Sprechen erkannt — geh näher ans Mikrofon und sprich das Wort deutlich.',
    stressOk: 'Betonung auf der richtigen Silbe: {i}. Silbe.',
    stressDiff: 'Erwartet war die {e}. Silbe betont, du hast die {d}. betont.',
    stressApprox: 'Hinweis: Die Betonung für {lang} ist ohne Wörterbuch nur eine Näherung — im Wörterbuch nachsehen gibt das richtige Ergebnis.',
    stressUnknown: 'Betonte Silbe konnte nicht erkannt werden.',
  },
  en: {
    typeWord: 'Type a word above first — the expected syllables will appear.',
    noData: 'No recording yet — say the word and get feedback from the coach.',
    countOk: 'Syllable count is right: {d}/{e}.',
    countMore: '{e} syllables were expected, {d} detected. You may have added an unstressed syllable — shorten the syllables a bit.',
    countLess: '{e} syllables were expected, {d} detected. You may be swallowing syllables — say the last syllable clearly.',
    countZero: 'No speech detected — move closer to the mic and say the word clearly.',
    stressOk: 'Stress on the right syllable: {i}.',
    stressDiff: 'Expected stress on syllable {e}, you stressed syllable {d}.',
    stressApprox: 'Note: stress rules for {lang} are approximate without a dictionary — checking a dictionary gives the exact result.',
    stressUnknown: 'Could not detect the stressed syllable.',
  },
}

function fill(tpl, vars) {
  return tpl.replace(/\{(\w+)\}/g, (_, k) => vars[k] ?? '?')
}

// Koç yorumları. Seviyeler: 'ok' (yeşil) | 'warn' (cyan).
export function evaluateCoach(expected, analysis, uiLang = 'tr') {
  const M = msg[uiLang] ?? msg.tr
  const notes = []

  if (!expected || !expected.syllables || expected.syllables.length === 0) {
    notes.push({ level: 'warn', text: M.typeWord })
    return notes
  }
  if (!analysis) {
    notes.push({ level: 'warn', text: M.noData })
    return notes
  }

  const e = expected.syllables.length
  const d = analysis.stats.syllableCount

  if (d === 0) {
    notes.push({ level: 'warn', text: M.countZero })
  } else if (d === e) {
    notes.push({ level: 'ok', text: fill(M.countOk, { d, e }) })
  } else if (d > e) {
    notes.push({ level: 'warn', text: fill(M.countMore, { e, d }) })
  } else {
    notes.push({ level: 'warn', text: fill(M.countLess, { e, d }) })
  }

  const spoken = analysis.stressIndex
  const exp = expected.stress
  if (spoken >= 0 && exp >= 0) {
    if (spoken === exp) {
      notes.push({ level: 'ok', text: fill(M.stressOk, { i: spoken + 1 }) })
    } else {
      notes.push({ level: 'warn', text: fill(M.stressDiff, { e: exp + 1, d: spoken + 1 }) })
    }
    if (expected.approx) {
      notes.push({ level: 'warn', text: fill(M.stressApprox, { lang: langName[uiLang]?.[expected.lang] ?? expected.lang }) })
    }
  } else if (d > 0 && spoken < 0) {
    notes.push({ level: 'warn', text: M.stressUnknown })
  }

  return notes
}

// Beklenen kelimeyi hece-hece, vurgulu hece işaretli olarak çizmek için yardımcı.
// Dönüş: [{ text, stressed }]
export function formatExpected(expected) {
  if (!expected || !expected.syllables || expected.syllables.length === 0) return []
  return expected.syllables.map((s, i) => ({ text: s.text, stressed: i === expected.stress }))
}
