package core

import "testing"

func TestSyllabifyWordTurkishFinalStress(t *testing.T) {
	w := SyllabifyWord("tr", "merhaba")
	if w == nil {
		t.Fatal("kelime boş olmamalı")
	}
	if len(w.Syllables) != 3 {
		t.Fatalf("merhaba=3 hece, bulunan=%d (%s)", len(w.Syllables), texts(w))
	}
	if w.Stress != 2 {
		t.Fatalf("Türkçede son hece vurgulu, beklenen=2 bulunan=%d", w.Stress)
	}
	if w.Approx {
		t.Fatal("Türkçe vurgu kuralı yaklaşık değil, son hece")
	}
}

func TestSyllabifyWordGermanPrefixStress(t *testing.T) {
	// verstehen → önek "ver" vurgusuz, kökün ilk hecesi vurgulu
	w := SyllabifyWord("de", "verstehen")
	if len(w.Syllables) != 3 {
		t.Fatalf("verstehen=3 hece, bulunan=%d (%s)", len(w.Syllables), texts(w))
	}
	if w.Stress != 1 {
		t.Fatalf("verstehen vurgusu 2. hecede (indeks 1), bulunan=%d", w.Stress)
	}
	if !w.Approx {
		t.Fatal("Almanca vurgu sözlüksüz yaklaşıktır")
	}

	// entschuldigung → önek "ent", vurgu "schul" hecesinde
	w = SyllabifyWord("de", "entschuldigung")
	if w.Stress != 1 {
		t.Fatalf("entschuldigung vurgusu 'schul' (indeks 1), bulunan=%d", w.Stress)
	}
}

func TestSyllabifyWordGermanNoPrefix(t *testing.T) {
	w := SyllabifyWord("de", "kaufen")
	if len(w.Syllables) != 2 {
		t.Fatalf("kaufen=2 hece, bulunan=%d (%s)", len(w.Syllables), texts(w))
	}
	if w.Stress != 0 {
		t.Fatalf("öneksiz Almanca kelime kök-ilk hece (indeks 0), bulunan=%d", w.Stress)
	}
}

func TestSyllabifyWordEnglishHeuristic(t *testing.T) {
	w := SyllabifyWord("en", "hello")
	if len(w.Syllables) != 2 {
		t.Fatalf("hello=2 hece, bulunan=%d (%s)", len(w.Syllables), texts(w))
	}
	if w.Stress != 0 {
		t.Fatalf("İngilizce kaba kural ilk hece (indeks 0), bulunan=%d", w.Stress)
	}
	if !w.Approx {
		t.Fatal("İngilizce vurgu yaklaşıktır")
	}
}

func TestSyllabifyWordEmpty(t *testing.T) {
	if SyllabifyWord("tr", "   ") != nil {
		t.Fatal("boş kelime nil döndürmeli")
	}
}

// texts, test hata mesajları için hece metinlerini "-" ile birleştirir.
func texts(w *ExpectedWord) string {
	out := ""
	for i, s := range w.Syllables {
		if i > 0 {
			out += "-"
		}
		out += s.Text
	}
	return out
}
