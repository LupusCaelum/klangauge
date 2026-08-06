package core

// Analysis, tek bir ses kaydının tüm analiz sonucudur.
// Bu struct, web tarafına JSON olarak gönderilecek.
type Analysis struct {
	SampleRate int     `json:"sampleRate"` // örnekleme hızı (örn. 44100)
	Duration   float64 `json:"duration"`   // toplam süre (saniye)
	FrameSize  int     `json:"frameSize"`  // her pencere kaç örnek (örn. 2048)
	Hop        int     `json:"hop"`        // pencereler kaç örnek kayarak ilerler (örn. 512)

	Waveform    []float64     `json:"waveform"`    // çizim için küçültülmüş [min,max] çiftleri
	WaveformPts int           `json:"waveformPts"` // kaç çift (min,max) olduğu
	Intensity   []FrameValue  `json:"intensity"`   // şiddet eğrisi (dB)
	Voiced      []bool        `json:"voiced"`      // hangi frame'ler sesli
	Pitch       []PitchFrame  `json:"pitch"`       // perde eğrisi (Hz)
	Spectrogram *Spectrogram  `json:"spectrogram"` // frekans-zaman haritası
	Syllables   []Syllable    `json:"syllables"`   // algılanan hece tepeleri
	StressIndex int           `json:"stressIndex"` // en belirgin hece indisi (-1 yoksa)
	Stats       *SummaryStats `json:"stats"`       // özet istatistikler
}

// Syllable, algılanan tek bir hecedir (şiddet eğrisindeki enerji tepesi).
// Koç, beklenen hecelerle karşılaştırırken zamanı ve öne çıkanlığı kullanır.
type Syllable struct {
	Time       float64 `json:"t"`          // tepenin zamanı (saniye)
	Prominence float64 `json:"prominence"` // tepenin öne çıkanlığı (dB)
}

// FrameValue, zaman eksenindeki tek bir ölçüm noktasıdır.
type FrameValue struct {
	Time  float64 `json:"t"`     // kaydın başından itibaren saniye
	Value float64 `json:"value"` // ölçülen değer
}

// PitchFrame, belirli bir anda tespit edilen perde bilgisidir.
type PitchFrame struct {
	Time       float64 `json:"t"`          // saniye
	Frequency  float64 `json:"hz"`         // tespit edilen temel frekans
	Voiced     bool    `json:"voiced"`     // sesli mi (konuşma var mı)?
	Confidence float64 `json:"confidence"` // 0..1, tespit ne kadar güvenilir
}

// Spectrogram, iki boyutlu frekans-zaman haritasıdır.
type Spectrogram struct {
	Times []float64   `json:"times"` // her sütunun zamanı (saniye)
	Freqs []float64   `json:"freqs"` // her satırın frekansı (Hz)
	DB    [][]float64 `json:"db"`    // [frame][frekans] -> desibel (dB)
}

// SummaryStats, kaydın genel özetidir.
type SummaryStats struct {
	VoicedRatio   float64 `json:"voicedRatio"`   // konuşma geçen süre oranı (0..1)
	MeanPitch     float64 `json:"meanPitch"`     // ortalama perde (Hz)
	PitchRange    float64 `json:"pitchRange"`    // perde aralığı (max-min, Hz)
	MeanIntensity float64 `json:"meanIntensity"` // ortalama şiddet (dB)

	SyllableCount    int     `json:"syllableCount"`    // algılanan hece sayısı
	SpeechDuration   float64 `json:"speechDuration"`   // konuşulan (sesli) toplam süre (s)
	SpeechRate       float64 `json:"speechRate"`       // hece/dakika (toplam süreye göre)
	ArticulationRate float64 `json:"articulationRate"` // hece/saniye (sadece sesli kısma göre)
}
