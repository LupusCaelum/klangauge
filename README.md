# Klangauge 🔊

> Çevrimdışı konuşma analizi + telaffuz koçu — Go çekirdeği WASM ile tarayıcıda çalışır, hiçbir ses kaydı cihazdan çıkmaz.
> Offline speech analysis + pronunciation coach — Go core compiled to WASM, every recording stays on your device.

**Canlı dene / Try it live:** <https://lupuscaelum.github.io/klangauge/>

Klangauge bir web tarayıcısında tamamen yerel çalışan bir konuşma analizi terminalidir. Mikrofonla kaydeder veya bir ses dosyası yükler; dalga biçimi, şiddet, perde ve spektrogramı çizer; hece sayısı, konuşma hızı, artikülasyon, sesli oranı, perde aralığı ve şiddet gibi metrikleri kural tabanlı yorumlarla birleştirir. Arayüz Türkçe, Almanca ve İngilizce'dir.

## Özellikler / Features

- **Dalga biçimi, şiddet eğrisi, perde çizgisi ve spektrogram** — YIN ile perde tespiti (`go-yinfft`)
- **Hece tespiti** — şiddet eğrisindeki enerji tepelerinden
- **Otomatik kural tabanlı yorum** — dil ayarlı hız/artikülasyon bantları (TR hece-zamanlı, DE/EN vurgu-zamanlı)
- **Perde karşılaştırma** — aynı kelimeyi iki kez söyle, iki perde çizgisi üst üste binsin
- **Telaffuz koçu** — kelime yaz, beklenen heceleri ve vurguyu gör, söyle, koç hece sayısı ve vurgu yerini karşılaştırsın ([Syllabifier](https://github.com/LupusCaelum/syllabifier) motoruyla)
- **Canlı seviye + canlı dalga** — kayıt sırasında anlık gösterge
- **%100 çevrimdışı** — sunucu yok, API yok, takip yok; PNG dışa aktarım

## Mimarı / Architecture

```
┌─────────────┐   compile (GOOS=js)   ┌──────────────────────────┐
│  core (Go)  │ ────────────────────► │  klangauge.wasm           │
│  analyze    │                       │  (tarayıcıda çalışır)     │
│  pitch/spec │                       └────────────┬─────────────┘
│  syllables  │                                    │ global JS API
└─────────────┘                                    ▼
                                       ┌──────────────────────────┐
                                       │  Vue 3 + Tailwind 4 UI   │
                                       │  canlı gösterge · koç    │
                                       └──────────────────────────┘
```

- `core/` — Go analiz motoru (spectrogram, YIN perde, şiddet, hece algılama, koç kuralları)
- `cmd/wasm/` — WASM sınırı: `klangaugeAnalyze(samples, sampleRate)` ve `klangaugeSyllabify(text, lang)` global JS API'leri
- `web/` — Vue 3 + Tailwind CSS 4 + Vite arayüzü

## Geliştirme / Development

Go 1.26 ve pnpm gerekir.

```bash
# 1) WASM modülünü derle
GOOS=js GOARCH=wasm go build -o web/public/klangauge.wasm ./cmd/wasm

# 2) Arayüzü çalıştır / derle
cd web
pnpm install
pnpm dev        # geliştirme
pnpm build      # üretim (dist/)
```

Testler:

```bash
go test ./core/...
```

## Dağıtım / Deployment

`push` → `main` (veya `workflow_dispatch`) → GitHub Actions `web/dist`'i GitHub Pages'e yayınlar. WASM dosyası repo'da commit'li olduğundan CI'da Go gerekmez. Action'lar node24 major sürümlerinde çalışır.

## Teknoloji / Stack

- **Go 1.26** — analiz çekirdeği: `mjibson/go-dsp`, `FreibergVlad/go-yinfft`, `go-audio`
- **WASM** — `syscall/js`, `wasm_exec.js`
- **Vue 3** + **Tailwind CSS 4** + **Vite 8**
- **GitHub Actions** — Pages deploy

## Lisans / License

Özel (tescilli) proje. Her hakkı saklıdır — kodun yeniden kullanımı için izin alın.
