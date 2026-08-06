package core

import (
	"math"

	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
)

// maxSpectrumHz: spektrogramda saklanacak üst frekans sınırı.
// Konuşma formantları (a/e/i vs. tonu) ~5 kHz altındadır; bu sayede hem
// görsel netleşir hem JSON küçük kalır.
const maxSpectrumHz = 5000

// AnalyzeSpectrogram, her frame için FFT uygulayarak frekans-zaman
// haritasını (dB cinsinden) üretir. Çıktı, canvas'ta renk olarak çizilir.
func AnalyzeSpectrogram(samples []float64, cfg Config) *Spectrogram {
	frames := cfg.FrameCount(len(samples))
	win := window.Hann(cfg.FrameSize)

	// Kaç frekans bandı tutacağız? (0..maxSpectrumHz)
	maxBin := int(float64(maxSpectrumHz) / float64(cfg.SampleRate) * float64(cfg.FrameSize))
	if maxBin > cfg.FrameSize/2 {
		maxBin = cfg.FrameSize / 2
	}

	times := make([]float64, frames)
	freqs := make([]float64, maxBin+1)
	for b := 0; b <= maxBin; b++ {
		// bin b'nin temsil ettiği frekans: b * (örnekleme hızı / pencere boyu)
		freqs[b] = float64(b) / float64(cfg.FrameSize) * float64(cfg.SampleRate)
	}

	db := make([][]float64, frames)
	frame := make([]float64, cfg.FrameSize)
	for i := 0; i < frames; i++ {
		start := i * cfg.Hop
		if start+cfg.FrameSize > len(samples) {
			start = len(samples) - cfg.FrameSize
		}
		// YINFFT dersinden: pencere uygulaması frame'i yerinde değiştirir,
		// o yüzden kopya üzerinde çalışıyoruz (samples'a dokunmayız).
		copy(frame, samples[start:start+cfg.FrameSize])
		for j := range frame {
			frame[j] *= win[j]
		}
		spectrum := fft.FFTReal(frame)

		row := make([]float64, maxBin+1)
		for b := 0; b <= maxBin; b++ {
			mag := math.Hypot(real(spectrum[b]), imag(spectrum[b]))
			// Tam genlikli (1.0) sinüsün pik bini 0 dB olsun diye 2/N ile ölçekliyoruz.
			d := 20 * math.Log10(2*mag/float64(cfg.FrameSize))
			if d < -90 {
				d = -90
			}
			if d > 0 {
				d = 0
			}
			row[b] = d
		}
		db[i] = row
		times[i] = float64(start) / float64(cfg.SampleRate)
	}

	return &Spectrogram{Times: times, Freqs: freqs, DB: db}
}
