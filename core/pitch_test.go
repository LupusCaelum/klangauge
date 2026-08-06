package core

import (
	"math"
	"testing"
)

// synthTone, verilen frekansta ve zarf ile tek bir ton üretir.
func synthTone(sr int, dur float64, freq float64) []float64 {
	n := int(dur * float64(sr))
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(sr)
		env := math.Sin(math.Pi * float64(i) / float64(n))
		out[i] = 0.8 * env * math.Sin(2*math.Pi*freq*t)
	}
	return out
}

// synthMelody, üç farklı perdede hece üretir: 150, 200, 300 Hz.
func synthMelody(sr int) []float64 {
	var out []float64
	for _, f := range []float64{150, 200, 300} {
		out = append(out, synthTone(sr, 0.5, f)...)
		out = append(out, make([]float64, int(0.3*float64(sr)))...)
	}
	return out
}

// yakinMi, iki frekansın yüzde olarak yakın olup olmadığını kontrol eder.
func yakinMi(got, want, yuzde float64) bool {
	if want == 0 {
		return false
	}
	return math.Abs(got-want)/want < yuzde/100
}

func TestPitchDetectsMelody(t *testing.T) {
	const sr = 44100
	samples := synthMelody(sr)

	_, voiced, _, _ := AnalyzeRhythm(samples, DefaultConfig(sr))
	pitch := AnalyzePitch(samples, DefaultConfig(sr), voiced)

	wants := []float64{150, 200, 300}
	found := make([]bool, len(wants))

	for _, p := range pitch {
		if !p.Voiced || p.Frequency == 0 {
			continue
		}
		for i, w := range wants {
			if yakinMi(p.Frequency, w, 5) {
				found[i] = true
			}
		}
	}
	for i, f := range found {
		if !f {
			t.Fatalf("%.0f Hz perdesi algılanamadı", wants[i])
		}
	}

	mean, rng := pitchStats(pitch)
	if mean < 120 || mean > 350 {
		t.Fatalf("beklenen ortalama ~200, bulunan=%f", mean)
	}
	if rng < 100 {
		t.Fatalf("beklenen perde aralığı >100, bulunan=%f", rng)
	}
}

func TestPitchSilenceIsUnvoiced(t *testing.T) {
	const sr = 44100
	samples := make([]float64, int(1.0*float64(sr)))
	_, voiced, _, _ := AnalyzeRhythm(samples, DefaultConfig(sr))
	pitch := AnalyzePitch(samples, DefaultConfig(sr), voiced)
	for _, p := range pitch {
		if p.Voiced {
			t.Fatalf("sessizlikte perde algılanmamalı, t=%f", p.Time)
		}
	}
}
