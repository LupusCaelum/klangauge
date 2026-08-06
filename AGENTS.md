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

1. **Pages deploy yayında DEĞİL — hece düzeltmesi canlıya çıkmadı.**
   - Son iki `workflow_dispatch` run GitHub'ın hosted runner kıtlığı nedeniyle başarısız oldu:
     `"The job was not acquired by Runner of type hosted even after multiple attempts"`.
   - Canlı wasm hâlâ ESKİ: 3442010 bayt, sha256 `530d9bd4...`, last-modified 20:32.
   - Yerel YENİ wasm: 3445344 bayt, sha256 `5d6e09c5...`.
   - Yapılacak: `gh workflow run "Deploy to GitHub Pages" --ref main`, sonra
     `curl -sI https://lupuscaelum.github.io/klangauge/klangauge.wasm` ile
     size=3445344 ve sha256=`5d6e09c5...` eşleştiğini doğrula.
   - Not: push-to-main otomatik tetikleme 25c849b'de ÇALIŞMADI (anomali) — manuel dispatch gerekli.

2. **Sonraki konu: beyin fırtınası fikirleri.** Kullanıcı "diğer beyin fırtınasındaki
   fikirlere bakalım" dedi. Fikir notu dosyası bulunamadı — kullanıcıya notun nerede
   olduğunu sor (veya roadmap tartışmasına geç). Sunum dokümanında geçen adaylar:
   CMUdict tabanlı EN vurgu sözlüğü, YIN/cspektral F0, obfuscation + lisans, WAV dışı
   formatlar, hedef telaffuz sesli referans, öğrenci ilerleme takibi.

## Tamamlananlar (referans)

- v1+v2 yayında (deploy bekleyen düzeltme hariç), portföy kartı + 3 dil blog push edildi (4b70c8c).
- Hece düzeltmesi `core/dsp.go` + regresyon testi `TestWordCountIgnoresConsonantBumps` geçti (13/13).
- Sunum dosyaları: `docs/sunum.md` (repo) + Desktop'ta `Klangauge_Sunum_Teknik.md` / `Klangauge_Sunum_Kullanici.md`.
