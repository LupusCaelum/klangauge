package core

import (
	"math"
	"sort"
)

// Config, analizin pencereleme parametreleridir.
type Config struct {
	FrameSize  int // her penceredeki örnek sayısı (örn. 2048)
	Hop        int // pencerelerin kayma miktarı (örn. 512)
	SampleRate int // örnekleme hızı
}

// DefaultConfig, konuşma analizi için uygun varsayılanlardır.
// FrameSize 2048 @ 44100 Hz ≈ 46 ms pencere; 75% örtüşme (hop 512).
func DefaultConfig(sampleRate int) Config {
	return Config{FrameSize: 2048, Hop: 512, SampleRate: sampleRate}
}

// FrameCount, verilen örnek sayısından kaç pencere çıkacağını hesaplar.
func (c Config) FrameCount(n int) int {
	if n < c.FrameSize {
		return 1
	}
	return (n-c.FrameSize)/c.Hop + 1
}

// RMS, bir pencerenin kök-ortalama-kare değeridir (şiddet göstergesi).
// Sıfıra yakınsa sessiz, büyükse gürültülü/şiddetli.
func RMS(frame []float64) float64 {
	sum := 0.0
	for _, s := range frame {
		sum += s * s
	}
	return math.Sqrt(sum / float64(len(frame)))
}

// toDB, RMS değerini desibele çevirir (20*log10, referans 1.0).
func toDB(rms float64) float64 {
	if rms <= 1e-9 {
		return -90
	}
	return 20 * math.Log10(rms)
}

// AnalyzeIntensity, ses örneklerini pencerelere böler ve her pencereye
// bir şiddet değeri (dB) atar. Böylece şiddetin zaman içindeki eğrisini elde ederiz.
func AnalyzeIntensity(samples []float64, cfg Config) []FrameValue {
	frames := cfg.FrameCount(len(samples))
	out := make([]FrameValue, frames)
	for i := 0; i < frames; i++ {
		start := i * cfg.Hop
		if start+cfg.FrameSize > len(samples) {
			start = len(samples) - cfg.FrameSize
		}
		db := toDB(RMS(samples[start : start+cfg.FrameSize]))
		if db < -70 {
			db = -70 // tam sessizliği değil, zemin gürültüsünü temsil etsin
		}
		out[i] = FrameValue{
			Time:  float64(start) / float64(cfg.SampleRate),
			Value: db,
		}
	}
	return out
}

// noiseFloor, kaydın gürültü tabanını tahmin eder.
// En düşük %20 şiddet değerinin ortalaması = sessizlik seviyesi.
func noiseFloor(frames []FrameValue) float64 {
	sorted := make([]float64, len(frames))
	for i, f := range frames {
		sorted[i] = f.Value
	}
	sort.Float64s(sorted)
	n := len(sorted) / 5
	if n < 1 {
		n = 1
	}
	sum := 0.0
	for _, v := range sorted[:n] {
		sum += v
	}
	return sum / float64(n)
}

// peakLevel, kaydın en yüksek %5 şiddet değerinin ortalamasıdır.
// "Konuşma zirvesi" olarak kullanılır — sessizlik hissini ayırt etmek için.
func peakLevel(frames []FrameValue) float64 {
	sorted := make([]float64, len(frames))
	for i, f := range frames {
		sorted[i] = f.Value
	}
	sort.Float64s(sorted)
	n := len(sorted) / 20
	if n < 1 {
		n = 1
	}
	sum := 0.0
	for _, v := range sorted[len(sorted)-n:] {
		sum += v
	}
	return sum / float64(n)
}

// MarkVoiced, her frame'in "sesli" (konuşma içeren) olup olmadığını belirler.
// İki koşuldan güçlü olanı geçerli eşiktir:
//   1) gürültü tabanının 12 dB üzeri — sessiz bir kayıtta nefes/zemin hissini ayıklar
//   2) konuşma zirvesinin 30 dB altı — zemin hissi düşükse bile "sessizlik" sayılmasın.
func MarkVoiced(frames []FrameValue) []bool {
	threshold := noiseFloor(frames) + 12
	if peak := peakLevel(frames); peak-30 > threshold {
		threshold = peak - 30
	}
	voiced := make([]bool, len(frames))
	for i, f := range frames {
		voiced[i] = f.Value > threshold
	}
	return voiced
}

// smooth, şiddet eğrisini üçgen ağırlıklı pencereyle yumuşatır.
// Ünsüz atakları (s, f, k…) ve zarf içi mikro dalgalanmalar ayrı hece
// gibi görünmesin diye tepe tespitinden ÖNCE uygulanır.
func smooth(vals []float64, radius int) []float64 {
	if len(vals) < 3 || radius < 1 {
		return vals
	}
	weights := make([]float64, radius*2+1)
	wsum := 0.0
	for i := -radius; i <= radius; i++ {
		w := float64(radius + 1 - abs(i))
		weights[i+radius] = w
		wsum += w
	}
	out := make([]float64, len(vals))
	for i := range vals {
		s := 0.0
		for j := -radius; j <= radius; j++ {
			idx := i + j
			if idx < 0 || idx >= len(vals) {
				continue
			}
			s += vals[idx] * weights[j+radius]
		}
		out[i] = s / wsum
	}
	return out
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DetectSyllables, şiddet eğrisindeki tepeleri bularak her hecenin zamanını
// ve öne çıkanlığını (prominence) döndürür. Her hecenin bir "enerji tepesi"
// vardır (a/e/o gibi sesli harflerde ses en güçlüdür).
// Eğri önce yumuşatılır; böylece ünsüz atakları heceye eklenmez.
// minGap: iki hece arası minimum süre (120ms) — hızlı konuşmada bile gerçek sınır.
// minProminence: tepenin iki yanındaki en yüksek "vadi"ye göre öne çıkanlığı (dB).
func DetectSyllables(times, vals []float64, voiced []bool, minGap, minProminence float64) []Syllable {
	vals = smooth(vals, 3)
	var peaks []int
	for i := 1; i < len(vals)-1; i++ {
		if !voiced[i] {
			continue
		}
		// yerel tepe: iki komşusundan da küçük değil (düz tepeler de yakalanır)
		if vals[i] >= vals[i-1] && vals[i] >= vals[i+1] {
			peaks = append(peaks, i)
		}
	}

	selected := []int{}
	for _, p := range peaks {
		if peakProminence(vals, p) < minProminence {
			continue
		}
		// Son seçilen tepeye çok yakınsa ya onu yükselt ya da bu tepedeni atla.
		if len(selected) > 0 && times[p]-times[selected[len(selected)-1]] < minGap {
			if vals[p] > vals[selected[len(selected)-1]] {
				selected[len(selected)-1] = p
			}
			continue
		}
		selected = append(selected, p)
	}

	out := make([]Syllable, 0, len(selected))
	for _, p := range selected {
		out = append(out, Syllable{Time: times[p], Prominence: peakProminence(vals, p)})
	}
	return out
}

// CountSyllables, algılanan hece sayısıdır (DetectSyllables'ın sayı hali).
func CountSyllables(times, vals []float64, voiced []bool, minGap, minProminence float64) int {
	return len(DetectSyllables(times, vals, voiced, minGap, minProminence))
}

// peakProminence, bir tepenin iki yanındaki en yüksek vadiye göre öne çıkanlığını
// hesaplar. Bu, düz zarf tepelerinde bile gerçek "enerji patlamasını" ölçer.
func peakProminence(vals []float64, p int) float64 {
	leftValley := vals[p]
	for i := p - 1; i >= 0 && vals[i] <= vals[p]; i-- {
		if vals[i] < leftValley {
			leftValley = vals[i]
		}
	}
	rightValley := vals[p]
	for i := p + 1; i < len(vals) && vals[i] <= vals[p]; i++ {
		if vals[i] < rightValley {
			rightValley = vals[i]
		}
	}
	valley := math.Max(leftValley, rightValley)
	return vals[p] - valley
}

// waveformBuckets: dalga biçimi çizimi için hedef çift (min,max) sayısı.
// Tüm ham örnekleri (44.1k/sn) değil, bu kadar noktayı tarayıcıya göndeririz.
const waveformBuckets = 1500

// DownsampleWaveform, sesi N parçaya böler ve her parçanın [min,max]
// değerlerini sırayla döndürür (2N uzunlukta). Canvas bunları dikey
// çizgilerle birleştirerek dolgulu dalga biçimini çizer.
func DownsampleWaveform(samples []float64, buckets int) []float64 {
	if len(samples) == 0 {
		return nil
	}
	if buckets < 1 {
		buckets = 1
	}
	if buckets > len(samples) {
		buckets = len(samples)
	}
	size := len(samples) / buckets
	if size < 1 {
		size = 1
	}
	out := make([]float64, 0, buckets*2)
	for b := 0; b < buckets; b++ {
		start := b * size
		end := start + size
		if end > len(samples) {
			end = len(samples)
		}
		if start >= end {
			break
		}
		lo, hi := samples[start], samples[start]
		for _, s := range samples[start:end] {
			if s < lo {
				lo = s
			}
			if s > hi {
				hi = s
			}
		}
		out = append(out, lo, hi)
	}
	return out
}

// AnalyzeRhythm, şiddet eğrisi, sesli bölgeler, algılanan heceler ve özet
// istatistikleri üretir. Bu, analizin "ritim" katmanıdır; perde ve spektrogram
// ayrı katmanlardır.
func AnalyzeRhythm(samples []float64, cfg Config) (intensity []FrameValue, voiced []bool, syllables []Syllable, stats *SummaryStats) {
	intensity = AnalyzeIntensity(samples, cfg)
	voiced = MarkVoiced(intensity)

	times := make([]float64, len(intensity))
	vals := make([]float64, len(intensity))
	for i, f := range intensity {
		times[i] = f.Time
		vals[i] = f.Value
	}

	syllables = DetectSyllables(times, vals, voiced, 0.12, 3.0)

	// Sesli süre: her sesli pencere hop/sampleRate saniye katkı verir.
	// MeanIntensity YALNIZCA sesli frame'lerden hesaplanır — sessizlik
	// (-70 dB) ortalamayı aşağı çekip "kısık konuşuyorsun" yanlış alarmı
	// üretmesin diye. Yani "konuşurken ne kadar gürsün" ölçüsüdür.
	speechDur := 0.0
	voicedIntensitySum := 0.0
	voicedCount := 0
	for i, v := range voiced {
		if v {
			speechDur += float64(cfg.Hop) / float64(cfg.SampleRate)
			voicedCount++
			voicedIntensitySum += vals[i]
		}
	}
	totalDur := float64(len(samples)) / float64(cfg.SampleRate)

	meanIntensity := 0.0
	if voicedCount > 0 {
		meanIntensity = voicedIntensitySum / float64(voicedCount)
	}

	stats = &SummaryStats{
		VoicedRatio:      float64(voicedCount) / float64(len(voiced)),
		MeanIntensity:    meanIntensity,
		SyllableCount:    len(syllables),
		SpeechDuration:   speechDur,
		SpeechRate:       float64(len(syllables)) / (totalDur / 60), // hece/dakika
		ArticulationRate: float64(len(syllables)) / speechDur,       // hece/saniye
	}
	return intensity, voiced, syllables, stats
}
