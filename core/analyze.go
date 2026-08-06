package core

// Analyze, ham ses örneklerini alıp tüm analiz katmanlarını çalıştırır:
// ritim (şiddet + hece), perde ve spektrogram.
func Analyze(samples []float64, sampleRate int) *Analysis {
	cfg := DefaultConfig(sampleRate)

	intensity, voiced, syllables, stats := AnalyzeRhythm(samples, cfg)

	pitch := AnalyzePitch(samples, cfg, voiced)
	stats.MeanPitch, stats.PitchRange = pitchStats(pitch)

	spectrogram := AnalyzeSpectrogram(samples, cfg)

	// Söylenen vurgu: en yüksek öne çıkanlığa sahip hece.
	stressIndex := -1
	for i, s := range syllables {
		if stressIndex == -1 || s.Prominence > syllables[stressIndex].Prominence {
			stressIndex = i
		}
	}

	return &Analysis{
		SampleRate:  sampleRate,
		Duration:    float64(len(samples)) / float64(sampleRate),
		FrameSize:   cfg.FrameSize,
		Hop:         cfg.Hop,
		Waveform:    DownsampleWaveform(samples, waveformBuckets),
		WaveformPts: waveformBuckets,
		Intensity:   intensity,
		Voiced:      voiced,
		Pitch:       pitch,
		Spectrogram: spectrogram,
		Syllables:   syllables,
		StressIndex: stressIndex,
		Stats:       stats,
	}
}
