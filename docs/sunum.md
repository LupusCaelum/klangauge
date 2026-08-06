# Klangauge — Sunum Bilgileri

Canlı: https://lupuscaelum.github.io/klangauge/ · Repo: https://github.com/LupusCaelum/klangauge

Bu doküman iki bölümden oluşur:

1. **Teknik/Bilimsel** — mimari, sinyal işleme algoritmaları, parametreler, test stratejisi. Mühendis/meslektaş ortamlarında kullanılır.
2. **Kullanıcı Sunumu** — ürünün değer önerisi, özellikler, demoya hazırlık. Son kullanıcıya/demo dinleyicisine anlatırken kullanılır.

---

## 1. Teknik / Bilimsel Bölüm

### 1.1 Konumlandırma (tek cümle)

> Tarayıcıda çalışan, **tamamen çevrimdışı** bir konuşma analizi terminali ve telaffuz koçu. Ses verisi hiçbir sunucuya gönderilmez; tüm DSP (sayısal sinyal işleme) Go çekirdeği içinde yürür ve **WASM** olarak tarayıcıya derlenir. LLM kullanmaz.

### 1.2 Mimari

```
                    ┌──────────────────────────────┐
                    │          Web UI (Vue 3)      │
                    │  App.vue · CoachView.vue     │
                    │  LiveCanvas · AnalysisCanvas │
                    └──────────────┬───────────────┘
                                   │  JSON  (WebAssembly.instantiate)
                    ┌──────────────▼───────────────┐
                    │   Go core  →  wasm32 JS      │
                    │  core/ (saf Go, sıfır bağım- │
                    │  lıklık, testli)             │
                    │  analyze · dsp · pitch ·     │
                    │  spectrogram · coach         │
                    └──────────────┬───────────────┘
                                   │
                    ┌──────────────▼───────────────┐
                    │   Audio pipeline (JS)        │
                    │  getUserMedia / WAV decode → │
                    │  Float32 PCM @ 44.1 kHz      │
                    └──────────────────────────────┘
```

- **Go çekirdeği saf Go'dur** (yalnızca stdlib) → kolay test, kolay WASM.
- Tek istisna: heceleme/vurgu için `github.com/LupusCaelum/syllabifier` paketi (aynı yazarın önceki projesi) — TR/DE/EN kural tabanlı heceleyici.
- Web tarafı ortak `useRecorder.js` composable'ı kullanır: kayıt sırasında `AnalyserNode` frekans verisini LiveCanvas'a akıtır, kayıt bitince Float32 PCM'i WASM'e verir.
- Deploy: GitHub Pages + Actions (wav derlenmiş repo'da → CI'da Go derlemesi gerekmez), `vite.config.js` `base:'./'` sayesinde alt yol (subpath) üzerinde çalışır.

### 1.3 DSP Zinciri (core/dsp.go)

Girdi: `Float32 PCM @ 44100 Hz` → `Analysis` struct'ı (JSON).

**Çerçeveleme**
- Pencere: `FrameSize = 2048` örnek (~46 ms), `Hop = 512` (~11.6 ms kayma). ~9 pencereler/saniye.
- Her çerçevede RMS şiddet dB cinsinden hesaplanır → `Intensity` eğrisi.

**Gürültü tabanı (`noiseFloor`)**
- Kaydın en düşük %20 şiddet çerçevesinin ortalaması → gürültü tabanı (dB).
- Farklı mikrofon/kayıt ortamlarına adapte olur (mutlak eşik yerine göreli).

**Seslilik sınıflandırması (`MarkVoiced`)** — adaptif eşik:
```
threshold = max(noiseFloor + 12 dB,  peakLevel − 30 dB)
```
- `peakLevel` = en yüksek %5 çerçevenin ortalama şiddeti (tepe/kırpma darbelerine dayanıklı).
- `threshold` üzerindeki çerçeveler "sesli" (konuşma içeriyor).
- Neden çift koşul: (a) gürültülü kayıtlarda tabanın 12 dB üstü zayıf ünlüleri kaybettirmez, (b) çok sessiz kayıtlarda `peak−30` tavanı nefes/friketifleri yanlışlıkla "hece" yapmaz.

**Hece algılama (`DetectSyllables`)**
1. Sesli şiddet eğrisi **üçgen pencere ile yumuşatılır** (`smooth`, radius = 3 → ~81 ms komşuluk). Amacı: ünsüz geçişlerinin/burun tınısının mikro tepelerini ve zarf dalgalanmalarını bastırmak.
2. Yumuşatılmış eğride yerel tepeler bulunur.
3. Her aday tepe filtrelenir:
   - **`minProminence = 3.0 dB`** — tepenin iki yanındaki en yüksek "vadi"ye göre öne çıkanlığı. < 3 dB öne çıkan tepeler gürültü/artefakt sayılır.
   - **`minGap = 0.12 s`** — ardışık kabul edilen tepeler en az 120 ms aralıklı olmalı. Gerçek bir hece bu süreye ulaşır; hızlı konuşmada bile güvenli sınır.
4. Sonuç: her hecenin **zamanı** ve **öne çıkanlığı (dB)**. Öne çıkanlık, "konuşmacı hangi heceyi vurguladı?" sorusunun cevabıdır → `StressIndex`.

**Neden bu parametreler? (tartışma)**
- Önceki sürüm: `minGap=0.08`, `minProminence=1.5`, yumuşatma yok. Sonuç: "selam" gibi bir kelime 4–5 "hece" üretiyordu (sessizlik, /s/ atağı, ünlü zarf tümsekleri ayrı ayrı sayılıyordu).
- Düzeltme üç cephede: (1) adaptif sesli eşiği sessizliği kesti, (2) yumuşatma mikro tepeleri yuttu, (3) daha seçici `minGap`+`minProminence` tepe-eleme yaptı.
- "selam-benzeri" sentetik kelime (frikatif atağı + tümsekli ünlüler + sessizlik): artık tam **2** hece.

**Ritim istatistikleri**
- `SpeechRate` = hece/dakika (toplam kayıt süresine göre).
- `ArticulationRate` = hece/saniye (**yalnızca sesli** kısma göre) — duraklamalardan bağımsız "ağız hızı" ölçüsü.
- Hedef dil TR/DE/EN bantları bunlara uygulanır (`explainStats`).

### 1.4 Perde Takibi (core/pitch.go)

- **Otokorelasyon** tabanlı temel frekans (F0) tespiti.
- Yalnızca "sesli" çerçeveler işlenir; sessizlik perdesi `0 / unvoiced` → gürültüde sahte perde üretilmez.
- Her çerçevede `Hz + confidence (0..1)` döner. `MeanPitch` ve `PitchRange` (max−min) özet istatistiklerdir — "monoton mu konuşuyorsun?" sorusunun nesnel karşılığı.
- Testler: `TestPitchDetectsMelody` (saf ton 165 Hz doğru bulunur), `TestPitchSilenceIsUnvoiced` (sessizlik unvoiced).

### 1.5 Spektrogram (core/spectrogram.go)

- Pencereli **FFT** (önceki pencere + cos yumuşatma), `10·log10(P)` ile dB.
- Frekans ekseni: sınırlı aralıkta örneklenir, çizim için satır başına frekans etiketi.
- Zaman ekseni: `Hop` hizasında her sütun.
- Testler: `TestSpectrogramFindsTone` (tam frekans satırında enerji tepe), `TestSpectrogramSilenceIsDark` (sessizlikte enerji yok).

### 1.6 Telaffuz Koçu (core/coach.go)

- Beklenen heceleme `SyllabifyWord(lang, word)`: TR (ünlü-ünsüz kuralları), DE (`NewGerman`), EN (`NewEnglish`) heceleyicileri.
- **Vurgu tahmini `expectedStress`:**
  - TR: son hece vurgusu — güvenilir kural (`approx=false`).
  - DE: kök-ilk hece + vurgusuz önekleri (be-, ge-, ver-…) affix analiziyle atlama — iyi ama sözlüksüz tam değil (`approx=true`).
  - EN: ilk hece sezgisi çok kabadır, İngilizce vurgu gerçekte sözlüğe bağlıdır (`approx=true`).
- Karşılaştırma: algılanan hece **sayısı** ve en yüksek öne çıkanlığın **zamanı** → beklenen vurgu hece indisiyle eşleşir. Eşleşme varsa "vurgu doğru" yoksa "vurgu beklenen hecede değil" + hangi hece olduğu.
- Koç notları, "söylediğin kelimeyi ayrıştır" mantığıyla `coach.js` tarafında TR/DE/EN üretilir.

### 1.7 Test Stratejisi

Go tarafında **saf sinyal sentezi** ile birim testler — deterministik, mikrofon gerektirmez:

| Test | Sinyal | Doğrulama |
|---|---|---|
| `TestRhythmCountsSyllables` | sentetik hece dizisi | hece sayısı, periyot |
| `TestSilenceHasNoSyllables` | saf gürültü | 0 hece |
| `TestWordCountIgnoresConsonantBumps` | "selam" modeli (frikatif + tümsekli ünlüler) | **2** hece (regresyon testi) |
| `TestPitchDetectsMelody` | 165 Hz saf ton | F0 ≈ 165 Hz |
| `TestPitchSilenceIsUnvoiced` | gürültü | unvoiced |
| `TestSpectrogramFindsTone` | saf ton | doğru frekans satırında enerji |
| `TestSyllabifyWord*` | kelime | TR/DE/EN hece + vurgu konumu |

Sonuç: `go test ./core/...` → **13 test, tamamı geçer.**

Uçtan uca (E2E): headless Chromium (`osascript`/Playwright) ile canlı site — 3 dilde kayıt+dosya yükleme akışı, koç adımları, canvas render, DOM varlığı doğrulanır.

### 1.8 Güvenlik / Gizlilik Argümanı

- Kayıt `getUserMedia` ile cihazda alınır, **ağdan hiç çıkmaz**.
- WASM modülü ve Go çekirdeği istemcide koşar; sunucu yok, telemetri yok, çerez yok.
- Bu mimari, GDPR/MKVKK hassasiyeti olan ses verisi için önemli bir satış noktasıdır.

### 1.9 Bilinen Sınırlar (dürüstlük maddesi)

- EN/DE vurgu tahmini sözlüksüz yaklaşıktır (`approx=true`); tam vurgu için CMUdict benzeri bir sözlük eklenebilir.
- Hece algılama enerji tepelerine dayanır; çok gürültülü ortamda `noiseFloor` adaptasyonu yetmeyebilir.
- F0 otokorelasyonu, arka planda başka ses varken kirlenebilir (uzun vadede YIN veya cepstral yöntemler düşünülebilir).

---

## 2. Kullanıcı Sunumu Bölümü

### 2.1 Tek cümlelik özet

> **Klangauge, dil öğrenirken konuşmanı görünür kılan çevrimdışı bir telaffuz koçudur.** Ses kaydını tarayıcıda analiz eder, ne kadar hızlı konuştuğunu, hangi heceyi vurguladığını ve sesinin perdelenmesini sana gösterir.

### 2.2 Ne işe yarar? (problem → çözüm)

- **Problem:** Dil öğrenirken "kelimeleri doğru mu söylüyorum?" sorusuna objektif cevap zor. Dinleme-ezber döngüsü yavaş; ses kaydınızı kimse puanlamıyor.
- **Çözüm:** Anlık, görsel, kişisel bir koç: konuş, hece sayısı/vurgu perde/süre anında ekranda çizilsin.

### 2.3 Özellikler (demo sırasına göre)

1. **Üç dil desteği** — Türkçe, Almanca, İngilizce. Arayüz de aynı üç dilde.
2. **"Konuştuğun dil" seçici** — analiz bantlarını (hız, artikülasyon) seçtiğin dile göre ayarlar; arayüz dilinden bağımsızdır.
3. **Canlı görselleştirme** — konuşurken dalga formu gerçek zamanlı çizilir (bu, demo'da en etkileyici an).
4. **Tek kelimeyle telaffuz koçu:**
   - Kelimeyi yaz → sistem TR/DE/EN kurallarıyla heceler ve vurgulu heceyi işaretler.
   - Söyle → algılanan hece sayısını "beklenen/gerçek" olarak karşılaştırır.
   - Vurgu: "vurgu son hecede değil, ilk hecede" gibi somut geri bildirim.
5. **Ölçümler:**
   - **Konuşma hızı** (hece/dakika) ve **artikülasyon hızı** (duraklamalar hariç hece/saniye).
   - **Perde aralığı** — monoton konuşup konuşmadığını söyler.
   - **Spektrogram** — sesin frekans haritası, ünlü formantlarının görseli.
6. **Ses dosyası yükleme** — kaydettiysen WAV dosyanı da analiz ettirebilirsin.
7. **Perde karşılaştırma** — iki kaydı üst üste bindirip farkını görürsün (altın renkli "referans" çizgisi).

### 2.4 Gizlilik vurgusu (önemli satış argümanı)

> Kaydın tarayıcıdan **asla çıkmaz**. Sunucu yok, analiz cihazında çalışır. Sesinle ilgili hiçbir veri bir yere gönderilmez.

### 2.5 Demo akışı (5 dakika)

1. Ana sayfayı göster: terminal görünümü, üç dil.
2. "Konuştuğun dil" → Türkçe seç.
3. Canlı görselleştirmeyi aç → kısa bir kelime söyle → dalga formunun canlı çizildiğini göster.
4. Koç: `merhaba` yaz → "mer-ha-ba, vurgu son hecede" beklentisi. Söyle → "Hece sayısı doğru: 3/3" + vurgu sonucu.
5. Hızlı konuş → artikülasyon hızının yükseldiğini göster.
6. Aynı kelimeyi iki farklı tonlamayla söyle → perde eğrisi farkı + karşılaştırma görünümü.
7. WAV dosyası yükle (hazır örnek) → hece/vurgu sonucu.
8. Kapanış cümlesi: "Tamamı tarayıcıda, veri cihazdan çıkmıyor."

### 2.6 Rakamlar (konuşurken kullanılacak sayılar)

- Kayıt: 44.1 kHz, 16 bit — mikrofon kalitesinde.
- Analiz çözünürlüğü: ~46 ms pencereler, ~12 ms adımlar.
- Hece algılama hassasiyeti: 120 ms minimum hece aralığı, 3 dB vurgu seçiciliği.
- Desteklenen diller: 3 (TR/DE/EN) — arayüz ve analiz.
- Kurulum: yok, tek URL, çevrimdışı çalışır (site bir kez yüklendikten sonra internet gerekmez).

### 2.7 SSS (soru-cevap hazırlığı)

- **"İnternet olmadan çalışır mı?"** — Siteyi bir kez açtıktan sonra tüm analiz cihazında koşar; WASM modülü zaten indirilmiş demektir.
- **"Neden LLM/AI kullanmıyor?"** — Çünkü telaffuz ölçümü deterministik sinyal işlemeyle yapılabilir ve veri gizliliği korunur; "AI" görünümüne değil doğruluğa yatırım yapıldı.
- **"Yanlış hece sayarsa?"** — Enerji tepelerine dayanır; çok gürültülü ortamları etkiler. Eşikler adaptiftir (kaydın kendi gürültü tabanına göre ayarlanır).
- **"Kayıtlar nereye gidiyor?"** — Hiçbir yere. Tarayıcı bellekten silinir.
- **"Almanca/İngilizce vurgu ne kadar güvenilir?"** — Türkçe son-hece kuralı kesin; Almanca/İngilizce sözlüksüz tahmindir, yaklaşık sonuç verir (arayüzde işaretlidir).

### 2.8 Kapanış mesajı önerisi

> "Klangauge, telaffuz öğrenmenin ölçülebilir ve mahremiyete saygılı hale getirilebileceğini gösteren, uçtan uca yazılmış bir proje. Sinyal işlemeden WASM'e, Vue arayüzünden test stratejisine kadar tamamı açık kaynak."
