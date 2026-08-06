<script setup>
import { onMounted, onBeforeUnmount, ref, watch } from 'vue'

const props = defineProps({
  analysis: { type: Object, default: null },
  compare: { type: Object, default: null }, // karşılaştırma analizi (perde overlay)
})

const specRef = ref(null)
const waveRef = ref(null)
let ro = null

const clamp = (v, lo, hi) => Math.max(lo, Math.min(hi, v))

// dB (-50..0) → neon renk gradyanı: saydam → yeşil → cyan → beyaz.
// Zemin -60 yerine -50 dB'de başlar → düşük enerjili kısımlar da görünür.
function specColor(dB) {
  const t = clamp((dB + 50) / 50, 0, 1)
  if (t <= 0.02) return [0, 0, 0, 0]
  if (t < 0.4) {
    const s = t / 0.4
    return [0, Math.round(150 * s), Math.round(95 * s), 255]
  }
  if (t < 0.78) {
    const s = (t - 0.4) / 0.38
    return [0, Math.round(150 + 105 * s), Math.round(95 + 160 * s), 255]
  }
  const s = (t - 0.78) / 0.22
  return [Math.round(140 * s), 255, 255, 255]
}

function draw() {
  const a = props.analysis
  drawSpectrogram(a)
  drawWaveform(a)
}

function drawSpectrogram(a) {
  const canvas = specRef.value
  if (!canvas || !a?.spectrogram) return
  const cssW = canvas.clientWidth
  const cssH = canvas.clientHeight
  const dpr = window.devicePixelRatio || 1
  if (canvas.width !== cssW * dpr || canvas.height !== cssH * dpr) {
    canvas.width = cssW * dpr
    canvas.height = cssH * dpr
  }
  const ctx = canvas.getContext('2d')
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, cssW, cssH)

  const spec = a.spectrogram
  const nTimes = spec.times.length
  const nFreqs = spec.freqs.length
  const maxFreq = spec.freqs[nFreqs - 1]

  const img = ctx.createImageData(cssW, cssH)
  for (let x = 0; x < cssW; x++) {
    const fi = Math.min(nTimes - 1, Math.floor((x / cssW) * nTimes))
    for (let y = 0; y < cssH; y++) {
      const rel = 1 - y / cssH
      const bi = Math.min(nFreqs - 1, Math.round(rel * (nFreqs - 1)))
      const [r, g, b, al] = specColor(spec.db[fi][bi])
      const idx = (y * cssW + x) * 4
      img.data[idx] = r
      img.data[idx + 1] = g
      img.data[idx + 2] = b
      img.data[idx + 3] = al
    }
  }
  ctx.putImageData(img, 0, 0)

  // perde çizgisi (pembe) — sesli frame'lerde
  if (a.pitch?.length) {
    ctx.beginPath()
    ctx.strokeStyle = '#ff2ec4'
    ctx.lineWidth = 1.5
    ctx.shadowColor = '#ff2ec4'
    ctx.shadowBlur = 6
    for (let i = 0; i < a.pitch.length; i++) {
      const p = a.pitch[i]
      if (!p.voiced) continue
      const x = (p.t / a.duration) * cssW
      const y = (1 - p.hz / maxFreq) * cssH
      if (i === 0 || !a.pitch[i - 1].voiced) ctx.moveTo(x, y)
      else ctx.lineTo(x, y)
    }
    ctx.stroke()
    ctx.shadowBlur = 0
  }

  // karşılaştırma: ikinci kaydın perde çizgisi (altın) — zamanı 0..1'e normalle.
  if (props.compare?.pitch?.length) {
    const c = props.compare
    ctx.beginPath()
    ctx.strokeStyle = '#ffd166'
    ctx.lineWidth = 1.5
    ctx.shadowColor = '#ffd166'
    ctx.shadowBlur = 6
    let started = false
    for (let i = 0; i < c.pitch.length; i++) {
      const p = c.pitch[i]
      if (!p.voiced) {
        started = false
        continue
      }
      const x = (p.t / c.duration) * cssW
      const y = (1 - p.hz / maxFreq) * cssH
      if (!started) {
        ctx.moveTo(x, y)
        started = true
      } else {
        ctx.lineTo(x, y)
      }
    }
    ctx.stroke()
    ctx.shadowBlur = 0
  }

  // frekans ekseni etiketleri
  ctx.fillStyle = 'rgba(114,143,128,0.95)'
  ctx.font = '10px "JetBrains Mono", monospace'
  ctx.textAlign = 'left'
  for (let f = 0; f <= maxFreq; f += 1000) {
    const y = (1 - f / maxFreq) * cssH
    ctx.fillText(f >= 1000 ? `${f / 1000}k` : `${f}`, 4, y + 10)
  }
}

function drawWaveform(a) {
  const canvas = waveRef.value
  if (!canvas || !a?.waveform) return
  const cssW = canvas.clientWidth
  const cssH = canvas.clientHeight
  const dpr = window.devicePixelRatio || 1
  if (canvas.width !== cssW * dpr || canvas.height !== cssH * dpr) {
    canvas.width = cssW * dpr
    canvas.height = cssH * dpr
  }
  const ctx = canvas.getContext('2d')
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, cssW, cssH)

  const mid = cssH / 2
  const pts = a.waveform
  const n = a.waveformPts || pts.length / 2

  // dalga biçimi: her çift (min,max) dikey çizgi
  ctx.strokeStyle = 'rgba(0,255,159,0.85)'
  ctx.lineWidth = 1
  ctx.shadowColor = '#00ff9f'
  ctx.shadowBlur = 3
  ctx.beginPath()
  for (let i = 0; i < n; i++) {
    const x = (i / n) * cssW
    const lo = mid + pts[i * 2] * mid * 0.9
    const hi = mid + pts[i * 2 + 1] * mid * 0.9
    ctx.moveTo(x, lo)
    ctx.lineTo(x, hi)
  }
  ctx.stroke()
  ctx.shadowBlur = 0

  // şiddet eğrisi (cyan): dB -70..-10 aralığını dikeye eşle
  if (a.intensity?.length) {
    ctx.beginPath()
    ctx.strokeStyle = '#00e5ff'
    ctx.lineWidth = 1.5
    ctx.shadowColor = '#00e5ff'
    ctx.shadowBlur = 4
    for (let i = 0; i < a.intensity.length; i++) {
      const x = (a.intensity[i].t / a.duration) * cssW
      const rel = clamp((a.intensity[i].value + 70) / 60, 0, 1)
      const y = cssH - 6 - rel * (cssH - 12)
      if (i === 0) ctx.moveTo(x, y)
      else ctx.lineTo(x, y)
    }
    ctx.stroke()
    ctx.shadowBlur = 0
  }

  // orta çizgi
  ctx.strokeStyle = 'rgba(95,122,111,0.4)'
  ctx.setLineDash([4, 6])
  ctx.beginPath()
  ctx.moveTo(0, mid)
  ctx.lineTo(cssW, mid)
  ctx.stroke()
  ctx.setLineDash([])

  // hece tepeleri: her hecede dikey çizgi, vurgulu hecede altın nokta
  if (a.syllables?.length) {
    for (let i = 0; i < a.syllables.length; i++) {
      const s = a.syllables[i]
      const x = (s.t / a.duration) * cssW
      const isStress = i === a.stressIndex
      ctx.strokeStyle = isStress ? 'rgba(255,209,102,0.9)' : 'rgba(0,255,159,0.35)'
      ctx.lineWidth = isStress ? 2 : 1
      ctx.beginPath()
      ctx.moveTo(x, 4)
      ctx.lineTo(x, cssH - 4)
      ctx.stroke()
      if (isStress) {
        ctx.fillStyle = '#ffd166'
        ctx.shadowColor = '#ffd166'
        ctx.shadowBlur = 6
        ctx.beginPath()
        ctx.arc(x, 8, 3, 0, Math.PI * 2)
        ctx.fill()
        ctx.shadowBlur = 0
      }
    }
  }
}

function resize() {
  draw()
}

onMounted(() => {
  ro = new ResizeObserver(resize)
  if (specRef.value) ro.observe(specRef.value)
  if (waveRef.value) ro.observe(waveRef.value)
  draw()
})
onBeforeUnmount(() => ro?.disconnect())

watch(() => props.analysis, () => draw(), { deep: false })
</script>

<template>
  <div class="flex flex-col gap-3">
    <div class="relative scanlines overflow-hidden rounded border border-line bg-panel">
      <canvas ref="specRef" class="block h-56 w-full" />
      <span class="pointer-events-none absolute right-2 top-1 text-[10px] uppercase tracking-widest text-dim">spektrogram · perde</span>
    </div>
    <div class="relative overflow-hidden rounded border border-line bg-panel">
      <canvas ref="waveRef" class="block h-28 w-full" />
      <span class="pointer-events-none absolute right-2 top-1 text-[10px] uppercase tracking-widest text-dim">dalga · şiddet</span>
    </div>
  </div>
</template>
