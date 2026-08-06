package core

import (
	"math"
	"testing"
)

// synthSteady, sabit genlikli (zarf siz) bir ton üretir — spektrogram testi için.
func synthSteady(sr int, dur float64, freq float64) []float64 {
	n := int(dur * float64(sr))
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(sr)
		out[i] = 0.8 * math.Sin(2*math.Pi*freq*t)
	}
	return out
}

// pikBin, belirli bir frame'de en yüksek dB'nin bulunduğu bin indeksini döndürür.
func pikBin(row []float64) int {
	best := 0
	for i := range row {
		if row[i] > row[best] {
			best = i
		}
	}
	return best
}

func TestSpectrogramFindsTone(t *testing.T) {
	const sr = 44100
	const toneHz = 1000.0
	samples := synthSteady(sr, 0.5, toneHz)

	spec := AnalyzeSpectrogram(samples, DefaultConfig(sr))

	if len(spec.DB) == 0 {
		t.Fatal("spektrogram boş")
	}
	if len(spec.Freqs) == 0 {
		t.Fatal("frekans ekseni boş")
	}
	// Hz cinsinden tüm eksen doğru sıralı mı?
	if spec.Freqs[0] != 0 || spec.Freqs[len(spec.Freqs)-1] > maxSpectrumHz {
		t.Fatalf("frekans ekseni hatalı: ilk=%f son=%f", spec.Freqs[0], spec.Freqs[len(spec.Freqs)-1])
	}

	// Tonun tam ortasındaki frame'i al (zarf yok, enerji sabit).
	mid := len(spec.DB) / 2
	bin := pikBin(spec.DB[mid])
	got := spec.Freqs[bin]

	// Frekans çözünürlüğü = sampleRate/FrameSize ≈ 21.5 Hz; %5 tolerans yeterli.
	if !yakinMi(got, toneHz, 5) {
		t.Fatalf("beklenen pik ~%f Hz, bulunan=%f Hz", toneHz, got)
	}
	if spec.DB[mid][bin] < -30 {
		t.Fatalf("ton şiddeti beklenenden düşük: %f dB", spec.DB[mid][bin])
	}
}

func TestSpectrogramSilenceIsDark(t *testing.T) {
	const sr = 44100
	samples := make([]float64, int(0.5*float64(sr)))
	spec := AnalyzeSpectrogram(samples, DefaultConfig(sr))

	maxdB := math.Inf(-1)
	for _, row := range spec.DB {
		for _, d := range row {
			if d > maxdB {
				maxdB = d
			}
		}
	}
	if maxdB > -70 {
		t.Fatalf("sessizlikte spektrogram parlamamalı, max=%.1f dB", maxdB)
	}
}
