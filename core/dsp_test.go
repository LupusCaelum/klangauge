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
