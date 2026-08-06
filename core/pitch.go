package core

import (
	"math"

	"github.com/FreibergVlad/go-yinfft"
)

// pitchConfidenceThreshold: YIN confidence (1 - yinMin).
// Temiz perdede ~0.95+, boğuk seste ~0.7 civarı. Altını "perde yok" sayarız.
const pitchConfidenceThreshold = 0.7

// AnalyzePitch, her frame'in temel frekansını (perde) tespit eder.
// intensityVoiced: ritim katmanının seslilik bilgisi — gürültüdeki yanlış
// tespitleri engellemek için perde bulunsa bile şiddet "sesli" demiyorsa
// frame'i sessiz kabul ederiz.
func AnalyzePitch(samples []float64, cfg Config, intensityVoiced []bool) []PitchFrame {
	frames := cfg.FrameCount(len(samples))

	detector, err := yinfft.New(yinfft.Params{
		FrameSize:         cfg.FrameSize,
		SampleRate:        float64(cfg.SampleRate),
		ShouldInterpolate: true,    // fraksiyonel perde hassasiyeti (ondalık Hz)
		Tolerance:         1,       // 1 = tüm sonuçları kabul et, eşiği kendimiz uyguluyoruz
		WeightingType:     "EMPTY", // ağırlıksız — konuşma bandını bozmaz
		MinFrequency:      60,      // konuşma perdesi alt sınırı
		MaxFrequency:      500,     // kadın/çocuk sesi üst sınırı
	})
	if err != nil {
		return nil
	}

	out := make([]PitchFrame, frames)
	frame := make([]float64, cfg.FrameSize)
	for i := 0; i < frames; i++ {
		start := i * cfg.Hop
		if start+cfg.FrameSize > len(samples) {
			start = len(samples) - cfg.FrameSize
		}
		// ÖNEMLİ: YINFFT frame'i in-place pencereleyip bozuyor.
		// Go'da slice diziye bir görünüm (pointer) olduğu için,
		// `samples[start:start+size]`'ı doğrudan geçersek tüm örtüşen
		// frame'ler çoktan bozulmuş veriyi görür. Bu yüzden kopyalıyoruz.
		copy(frame, samples[start:start+cfg.FrameSize])
		freq, conf, err := detector.DetectFromFrame(frame)

		voiced := err == nil && conf >= pitchConfidenceThreshold && intensityVoiced[i]
		if !voiced {
			freq = 0
		}
		out[i] = PitchFrame{
			Time:       float64(start) / float64(cfg.SampleRate),
			Frequency:  freq,
			Voiced:     voiced,
			Confidence: conf,
		}
	}
	return out
}

// pitchStats, sesli frame'lerden ortalama perde ve perde aralığını hesaplar.
func pitchStats(pitch []PitchFrame) (mean, rng float64) {
	min := math.Inf(1)
	max := math.Inf(-1)
	sum := 0.0
	n := 0
	for _, p := range pitch {
		if !p.Voiced {
			continue
		}
		sum += p.Frequency
		n++
		if p.Frequency < min {
			min = p.Frequency
		}
		if p.Frequency > max {
			max = p.Frequency
		}
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), max - min
}
