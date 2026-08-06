package core

import (
	"math"
	"math/rand"
	"testing"
)

// synthSyllable, belirli sürede, belirli frekansta, ortası en yüksek
// (Hann-benzeri) zarf ile bir "hece" üretir. Tek enerji tepesi vardır.
func synthSyllable(sr int, dur float64, freq float64) []float64 {
	n := int(dur * float64(sr))
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(sr)
		env := math.Sin(math.Pi * float64(i) / float64(n))
		out[i] = 0.8 * env * math.Sin(2*math.Pi*freq*t)
	}
	return out
}

// synthRecording, heceler ve sessizliklerden oluşan bir "cümle" sentezler.
func synthRecording(sr int, syllableDur, gapDur float64, nSyllables int) []float64 {
	var out []float64
	for i := 0; i < nSyllables; i++ {
		out = append(out, synthSyllable(sr, syllableDur, 150.0)...)
		out = append(out, make([]float64, int(gapDur*float64(sr)))...)
	}
	// mikrofon hissi gibi hafif gürültü tabanı ekle
	rng := rand.New(rand.NewSource(42))
	for i := range out {
		out[i] += 0.005 * (rng.Float64()*2 - 1)
	}
	return out
}

func TestRhythmCountsSyllables(t *testing.T) {
	const sr = 44100
	// 3 hece, her biri 0.25s + 0.35s sessizlik
	samples := synthRecording(sr, 0.25, 0.35, 3)

	intensity, voiced, syllables, stats := AnalyzeRhythm(samples, DefaultConfig(sr))

	if stats.SyllableCount != 3 {
		t.Fatalf("beklenen hece=3, bulunan=%d", stats.SyllableCount)
	}
	if len(syllables) != 3 {
		t.Fatalf("beklenen hece listesi=3, bulunan=%d", len(syllables))
	}
	if len(intensity) != len(voiced) {
		t.Fatalf("intensity (%d) ve voiced (%d) uzunlukları aynı olmalı", len(intensity), len(voiced))
	}
	if stats.VoicedRatio < 0.3 || stats.VoicedRatio > 0.6 {
		t.Fatalf("beklenen VoicedRatio ~0.42, bulunan=%f", stats.VoicedRatio)
	}
	if stats.SpeechRate <= 0 || stats.ArticulationRate <= 0 {
		t.Fatal("konuşma hızları pozitif olmalı")
	}
}

func TestSilenceHasNoSyllables(t *testing.T) {
	const sr = 44100
	samples := make([]float64, int(1.5*float64(sr))) // tamamen sessiz
	_, _, syllables, stats := AnalyzeRhythm(samples, DefaultConfig(sr))
	if stats.SyllableCount != 0 {
		t.Fatalf("sessizlikte hece=0 beklenir, bulunan=%d", stats.SyllableCount)
	}
	if len(syllables) != 0 {
		t.Fatalf("sessizlikte hece listesi boş olmalı, bulunan=%d", len(syllables))
	}
}

// synthVowelBump, ana tepeye yakın küçük bir ikincil enerji tümseği eklenmiş
// "hece" üretir. Gerçek konuşmada ünsüz geçişleri / burun tınısı böyle
// mikro tepeler yaratır ve aşırı hassas algılama bunları hece sayar.
// attack: hecenin başına eklenen frikatif benzeri atağı (ünsüz akışı)
// ünlüye SÜREKLİ geçer — arada vadi olmaz.
func synthVowelBump(sr int, dur, freq, amp float64, bumpT float64, attack float64) []float64 {
	n := int(dur * float64(sr))
	out := make([]float64, n)
	rng := rand.New(rand.NewSource(11))
	for i := 0; i < n; i++ {
		t := float64(i) / float64(sr)
		env := math.Sin(math.Pi * float64(i) / float64(n))
		// atağı (frikatif) ünlü zarfının altına göm: sürekli enerji
		attackEnv := attack * (1 - env) * (rng.Float64()*2 - 1)
		out[i] = amp * env * math.Sin(2*math.Pi*freq*t) + attackEnv
	}
	const bumpDur = 0.05
	b0 := int(bumpT * float64(sr))
	bn := int(bumpDur * float64(sr))
	for i := 0; i < bn && b0+i < n; i++ {
		env := math.Sin(math.Pi * float64(i) / float64(bn))
		out[b0+i] += amp * 0.45 * env * math.Sin(2*math.Pi*freq*float64(b0+i)/float64(sr))
	}
	return out
}

// synthNoiseRamp, zarf ile şekillendirilmiş gürültü (ünsüz / nefes taklidi).
func synthNoiseRamp(sr int, dur, amp float64, seed int64) []float64 {
	n := int(dur * float64(sr))
	rng := rand.New(rand.NewSource(seed))
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		env := 0.5 + 0.5*math.Cos(math.Pi*float64(i)/float64(n))
		out[i] = amp * env * (rng.Float64()*2 - 1)
	}
	return out
}

func TestWordCountIgnoresConsonantBumps(t *testing.T) {
	const sr = 44100
	// "selam": öndeki sessizlik, /s/+/e/ (frikatif atağı + tümsekli ünlü),
	// /l/ geçişi (dip), /a/+/m/ (tümsekli ünlü), sondaki sessizlik. Gerçek hece: 2.
	parts := [][]float64{
		synthNoiseRamp(sr, 0.5, 0.006, 1), // öndeki sessizlik
		synthVowelBump(sr, 0.28, 200, 0.8, 0.16, 0.25),
		synthNoiseRamp(sr, 0.07, 0.05, 3), // /l/ geçişi (dip)
		synthVowelBump(sr, 0.34, 150, 0.8, 0.2, 0.1),
		synthNoiseRamp(sr, 0.5, 0.006, 4), // sondaki sessizlik
	}
	var samples []float64
	for _, p := range parts {
		samples = append(samples, p...)
	}
	_, _, syllables, stats := AnalyzeRhythm(samples, DefaultConfig(sr))
	if stats.SyllableCount != 2 {
		ts := make([]float64, len(syllables))
		for i, s := range syllables {
			ts[i] = s.Time
		}
		t.Fatalf("selam-benzeri kelime: beklenen hece=2, bulunan=%d (zamanlar: %v)", stats.SyllableCount, ts)
	}
}
