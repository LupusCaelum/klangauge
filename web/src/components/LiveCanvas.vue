<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { readLevel } from '../recorder.js'

const props = defineProps({
  analyser: { type: Object, default: null },
})

const waveRef = ref(null)
const seconds = ref(0)
const level = ref(-60)
const db = ref('-∞')

const BUCKETS = 600 // 600 × ~46ms ≈ 27s görünür pencere
const ring = new Float32Array(BUCKETS * 2)
let ringLen = 0
let raf = null
let start = 0
let buf = null

function drawWave() {
  const canvas = waveRef.value
  if (!canvas) return
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
  const n = ringLen
  ctx.strokeStyle = 'rgba(0,255,159,0.85)'
  ctx.lineWidth = 1
  ctx.shadowColor = '#00ff9f'
  ctx.shadowBlur = 3
  ctx.beginPath()
  for (let i = 0; i < n; i++) {
    const x = (i / BUCKETS) * cssW
    const lo = mid + ring[i * 2] * mid * 0.9
    const hi = mid + ring[i * 2 + 1] * mid * 0.9
    ctx.moveTo(x, lo)
    ctx.lineTo(x, hi)
  }
  ctx.stroke()
  ctx.shadowBlur = 0

  // sağ kenarda "canlı" imleci
  ctx.fillStyle = '#00e5ff'
  ctx.beginPath()
  ctx.arc(cssW - 3, mid, 3, 0, Math.PI * 2)
  ctx.fill()
}

function tick() {
  const a = props.analyser
  if (a) {
    level.value = readLevel(a)
    db.value = level.value <= -60 ? '-∞' : level.value.toFixed(1)
    if (!buf) buf = new Float32Array(a.fftSize)
    a.getFloatTimeDomainData(buf)
    let lo = Infinity
    let hi = -Infinity
    for (let i = 0; i < buf.length; i++) {
      if (buf[i] < lo) lo = buf[i]
      if (buf[i] > hi) hi = buf[i]
    }
    if (ringLen < BUCKETS) {
      ring[ringLen * 2] = lo === Infinity ? 0 : lo
      ring[ringLen * 2 + 1] = hi === -Infinity ? 0 : hi
      ringLen++
    } else {
      ring.copyWithin(0, 2)
      ring[ringLen * 2 - 2] = lo === Infinity ? 0 : lo
      ring[ringLen * 2 - 1] = hi === -Infinity ? 0 : hi
    }
    seconds.value = (performance.now() - start) / 1000
    drawWave()
  }
  raf = requestAnimationFrame(tick)
}

onMounted(() => {
  start = performance.now()
  raf = requestAnimationFrame(tick)
})
onBeforeUnmount(() => cancelAnimationFrame(raf))

const levelPct = () => Math.max(0, Math.min(100, ((level.value + 60) / 50) * 100))
const levelCls = () => {
  if (level.value < -35) return 'bg-line'
  if (level.value < -20) return 'bg-neon'
  return 'bg-cyan'
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="rounded border border-line bg-panel p-3">
      <div class="flex items-center justify-between text-[10px] uppercase tracking-widest text-dim">
        <span>level</span>
        <span :class="level.value < -35 ? 'text-dim' : level.value < -20 ? 'text-neon' : 'text-cyan'">{{ db }} dB</span>
      </div>
      <div class="mt-2 h-2 w-full overflow-hidden rounded bg-black/40">
        <div class="h-full transition-[width] duration-75" :class="levelCls()" :style="{ width: levelPct() + '%' }" />
      </div>
      <div class="mt-1 text-right text-[10px] tabular-nums text-dim">{{ seconds.toFixed(1) }}s</div>
    </div>
    <div class="relative overflow-hidden rounded border border-line bg-panel">
      <canvas ref="waveRef" class="block h-28 w-full" />
      <span class="pointer-events-none absolute right-2 top-1 text-[10px] uppercase tracking-widest text-cyan">canlı</span>
    </div>
  </div>
</template>
