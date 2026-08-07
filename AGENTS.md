# Klangauge — AGENTS.md

Çevrimdışı konuşma analizi + telaffuz koçu. Go/WASM çekirdek + Vue 3. TR/DE/EN.
Canlı: https://lupuscaelum.github.io/klangauge/ · Repo: github.com/LupusCaelum/klangauge

## Komutlar

- Core test: `go test ./core/...`
- WASM derle: `GOOS=js GOARCH=wasm go build -o web/public/klangauge.wasm ./cmd/wasm`
- Web build: `cd web && pnpm install --frozen-lockfile && pnpm build`
- E2E: `/tmp/opencode/e2e/` (headless Chromium, koç akışı dahil)

## Mimari

- `core/` saf Go (stdlib), `Analyze()` tüm zinciri çalıştırır: dsp → pitch → spectrogram → coach.
- Hece algılama: `DetectSyllables(smooth radius 3, minGap 0.12, minProminence 3.0)`; seslilik eşiği `max(noiseFloor+12, peak-30)`.
- Deploy: `.github/workflows/pages.yml` → GitHub Pages (upload-pages-artifact, `web/dist`).

## ⚠️ SON OTURUMDAN KALANLAR (açılışta önce bunlar)

1. ✅ **Pages deploy tamamlandı (07.08.2026).** Hece düzeltmesi canlıda:
   canlı wasm = yerel wasm = 3445344 bayt, sha256 `5d6e09c5...`.
   Ders: push-to-main otomatik tetikleme çalışmıyor (anomali) — manuel
   `gh workflow run "Deploy to GitHub Pages" --ref main` gerekli; ayrıca GitHub
   hosted runner kıtlığı yaşanabiliyor, denemeler tekrarlanmalı.

2. **Sonraki konu: beyin fırtınası fikirleri.** Kullanıcı "diğer beyin fırtınasındaki
   fikirlere bakalım" dedi. Fikir notu dosyası bulunamadı — kullanıcıya notun nerede
   olduğunu sor (veya roadmap tartışmasına geç). Sunum dokümanında geçen adaylar:
   CMUdict tabanlı EN vurgu sözlüğü, YIN/cspektral F0, obfuscation + lisans, WAV dışı
   formatlar, hedef telaffuz sesli referans, öğrenci ilerleme takibi.

## Tamamlananlar (referans)

- v1+v2 yayında (deploy bekleyen düzeltme hariç), portföy kartı + 3 dil blog push edildi (4b70c8c).
- Hece düzeltmesi `core/dsp.go` + regresyon testi `TestWordCountIgnoresConsonantBumps` geçti (13/13).
- Sunum dosyaları: `docs/sunum.md` (repo) + Desktop'ta `Klangauge_Sunum_Teknik.md` / `Klangauge_Sunum_Kullanici.md`.
