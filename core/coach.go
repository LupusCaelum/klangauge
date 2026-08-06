package core

import (
	"strings"

	"github.com/LupusCaelum/syllabifier/affix"
	"github.com/LupusCaelum/syllabifier/syllabify"
)

// ExpectedSyllable, bir kelimenin beklenen (söylenmesi gereken) tek hecesidir.
type ExpectedSyllable struct {
	Text  string `json:"text"`  // hecenin metni (örn. "mer")
	Vowel string `json:"vowel"` // hecenin ünlüsü (örn. "e")
}

// ExpectedWord, dil kurallarına göre hesaplanan "doğru heceleme"dir.
// Koç, söylenen (Analysis.Syllables) ile bunu karşılaştırır.
type ExpectedWord struct {
	Lang      string             `json:"lang"`
	Word      string             `json:"word"`
	Syllables []ExpectedSyllable `json:"syllables"`
	Stress    int                `json:"stress"` // beklenen vurgu hece indisi (-1 bilinmiyor)
	Approx    bool               `json:"approx"` // vurgu kuralı yaklaşık mı? (de/en için evet)
}

// SyllabifyWord, kelimeyi dile göre heceler ve vurgu konumunu tahmin eder.
// Dil "tr" dışındaki bilinmeyen değerlerde varsayılan Türkçedir.
func SyllabifyWord(lang, word string) *ExpectedWord {
	word = strings.TrimSpace(word)
	if word == "" {
		return nil
	}

	var sf syllabify.Syllabifier
	switch lang {
	case "de":
		sf = syllabify.NewGerman()
	case "en":
		sf = syllabify.NewEnglish()
	default:
		lang = "tr"
		sf = syllabify.NewTurkish()
	}

	res := &ExpectedWord{Lang: lang, Word: word, Stress: -1}
	for _, s := range sf.Syllabify(word) {
		res.Syllables = append(res.Syllables, ExpectedSyllable{Text: s.Text, Vowel: string(s.Vowel)})
	}
	if len(res.Syllables) > 0 {
		res.Stress, res.Approx = expectedStress(lang, word, res.Syllables)
	}
	return res
}

// expectedStress, dil kurallarına göre vurgulu heceyi tahmin eder.
// Dönen ikinci değer, tahminin ne kadar güvenilir olduğunu söyler:
//   - tr: son hece kuralı iyi bir varsayılandır (approx=false).
//   - de: kök-ilk hece + vurgusuz önekleri atma kuralı; sözlüksüz tam doğru değildir (approx=true).
//   - en: ilk hece sezgisi çok kabadır; İngilizce vurgu sözlük gerektirir (approx=true).
func expectedStress(lang, word string, syllables []ExpectedSyllable) (int, bool) {
	n := len(syllables)
	if n == 0 {
		return -1, false
	}

	switch lang {
	case "de":
		// Vurgusuz önekler (be-, ge-, ver-…) kök vurgusundan önce gelir;
		// kökün ilk hecesi vurgulu sayılır. Önek hece sayısı kadar kaydır.
		_, matches := affix.German.Analyze(word)
		prefixSylls := 0
		for _, m := range matches {
			if m.Kind == "prefix" {
				prefixSylls += len(syllabify.NewGerman().Syllabify(m.Text))
			}
		}
		if prefixSylls >= n {
			prefixSylls = 0
		}
		return prefixSylls, true
	case "en":
		// Çok kaba ilk-hece sezgisi. İngilizce vurgu sözlüğe bağlıdır.
		return 0, true
	default: // tr
		return n - 1, false
	}
}
