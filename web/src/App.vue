<script setup>
import { ref, onMounted, computed } from 'vue'
import { toPng } from 'html-to-image'
import { loadWasm } from './wasm.js'
import { useRecorder } from './useRecorder.js'
import AnalysisCanvas from './components/AnalysisCanvas.vue'
import LiveCanvas from './components/LiveCanvas.vue'
import HelpPage from './components/HelpPage.vue'
import CoachView from './components/CoachView.vue'
import { explainStats } from './explainer.js'
import { lang, t, setLang } from './i18n.js'

const view = ref('analysis') // analysis | coach | help
const wasmState = ref('loading') // loading | ready | error
const wasmError = ref('')
const analysis = ref(null)
const stored = ref(null) // perde karşılaştırma için saklanan analiz
const exporting = ref(false)
const exportRef = ref(null)

const targetLang = ref('tr')
try {
  const saved = localStorage.getItem('klangauge-target')
  if (saved && ['tr', 'de', 'en'].includes(saved)) targetLang.value = saved
} catch {
  /* yoksay */
}

function setTarget(l) {
  targetLang.value = l
  try {
    localStorage.setItem('klangauge-target', l)
  } catch {
    /* yoksay */
  }
}

const { recState, recError, recSeconds, analyser, toggleRecord, loadFile } = useRecorder(
  (result) => (analysis.value = result),
)

const autoNotes = computed(() =>
  analysis.value && !analysis.value.error
    ? explainStats(analysis.value.stats, lang.value, targetLang.value)
    : [],
)

const status = computed(() => {
  if (wasmState.value === 'loading') return t('stWasmLoading')
  if (wasmState.value === 'error') return `${t('stWasmError')}: ${wasmError.value}`
  if (recState.value === 'recording') return t('stRecording').replace('{s}', recSeconds.value.toFixed(1))
  if (recState.value === 'analyzing') return t('stAnalyzing')
  if (recState.value === 'error') return `${t('stError')}: ${recError.value}`
  if (analysis.value) return t('stDone')
  return t('stIdle')
})

const statCards = computed(() => {
  if (!analysis.value) return []
  const a = analysis.value
  return [
    { label: t('lblDuration'), value: fmt(a.duration), unit: t('unitSec'), cls: 'text-cyan glow-cyan' },
    { label: t('lblSyllables'), value: String(a.stats.syllableCount), unit: '', cls: 'text-neon glow-green' },
    { label: t('lblSpeechRate'), value: fmt(a.stats.speechRate), unit: t('unitSylMin'), cls: 'text-neon glow-green' },
    { label: t('lblArtic'), value: fmt(a.stats.articulationRate, 2), unit: t('unitSylSec'), cls: 'text-cyan glow-cyan' },
    { label: t('lblVoicedRatio'), value: Math.round(a.stats.voicedRatio * 100) + '%', unit: '', cls: 'text-txt' },
    { label: t('lblMeanPitch'), value: fmt(a.stats.meanPitch), unit: 'Hz', cls: 'text-mag glow-mag' },
    { label: t('lblPitchRange'), value: fmt(a.stats.pitchRange), unit: 'Hz', cls: 'text-txt' },
    { label: t('lblMeanInt'), value: fmt(a.stats.meanIntensity), unit: 'dB', cls: 'text-txt' },
  ]
})

const fmt = (n, d = 1) => (Number.isFinite(n) ? n.toFixed(d) : '—')

onMounted(() => {
  loadWasm()
    .then(() => (wasmState.value = 'ready'))
    .catch((e) => {
      wasmState.value = 'error'
      wasmError.value = e.message
    })
})

async function onFile(event) {
  const file = event.target.files?.[0]
  if (!file) return
  try {
    await loadFile(file)
  } catch (e) {
    recError.value = e.message
  }
  event.target.value = ''
}

async function savePng() {
  if (!exportRef.value || exporting.value) return
  exporting.value = true
  try {
    const dataUrl = await toPng(exportRef.value, {
      backgroundColor: '#0c1310',
      pixelRatio: 2,
      cacheBust: true,
    })
    const a = document.createElement('a')
    a.href = dataUrl
    a.download = `klangauge-${new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')}.png`
    a.click()
  } catch (e) {
    recState.value = 'error'
    recError.value = t('pngError') + ': ' + e.message
  } finally {
    exporting.value = false
  }
}
</script>

<template>
  <div class="mx-auto flex min-h-screen max-w-5xl flex-col gap-6 px-4 py-6">
    <header class="flex items-end justify-between border-b border-line pb-3">
      <div>
        <h1 class="glow-green cursor text-2xl font-bold tracking-widest text-neon">KLANGAUGE</h1>
        <p class="mt-1 text-xs text-dim">{{ t('subtitle') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <div class="flex gap-1">
          <button
            v-for="l in ['tr', 'de', 'en']"
            :key="l"
            class="rounded border px-2 py-1 text-[10px] uppercase tracking-widest transition"
            :class="lang === l ? 'border-cyan/60 text-cyan glow-box' : 'border-line text-dim hover:text-txt'"
            @click="setLang(l)"
          >
            {{ l }}
          </button>
        </div>
        <span
          class="rounded border px-2 py-1 text-[10px] uppercase tracking-widest"
          :class="
            wasmState === 'ready'
              ? 'border-neon/50 text-neon glow-box'
              : wasmState === 'error'
                ? 'border-mag/60 text-mag'
                : 'border-line text-dim'
          "
        >{{ wasmState }}</span>
      </div>
    </header>

    <nav class="flex gap-2 border-b border-line pb-3">
      <button
        class="rounded border px-3 py-1 text-xs uppercase tracking-widest transition"
        :class="view === 'analysis' ? 'border-neon/60 text-neon' : 'border-line text-dim hover:text-txt'"
        @click="view = 'analysis'"
      >
        {{ t('tabAnalysis') }}
      </button>
      <button
        class="rounded border px-3 py-1 text-xs uppercase tracking-widest transition"
        :class="view === 'coach' ? 'border-neon/60 text-neon' : 'border-line text-dim hover:text-txt'"
        @click="view = 'coach'"
      >
        {{ t('tabCoach') }}
      </button>
      <button
        class="rounded border px-3 py-1 text-xs uppercase tracking-widest transition"
        :class="view === 'help' ? 'border-cyan/60 text-cyan' : 'border-line text-dim hover:text-txt'"
        @click="view = 'help'"
      >
        {{ t('tabHelp') }}
      </button>
    </nav>

    <CoachView v-if="view === 'coach'" />
    <HelpPage v-else-if="view === 'help'" />

    <template v-else>
      <section class="flex flex-wrap items-center gap-3">
        <button
          class="rounded border px-4 py-2 font-bold tracking-wider transition"
          :class="
            recState === 'recording'
              ? 'border-mag text-mag glow-box'
              : 'border-neon text-neon hover:bg-neon/10 glow-box'
          "
          :disabled="wasmState !== 'ready' || recState === 'analyzing'"
          @click="toggleRecord"
        >
          {{ recState === 'recording' ? t('stop') : t('record') }}
        </button>
        <label
          class="cursor-pointer rounded border border-cyan px-4 py-2 font-bold tracking-wider text-cyan transition hover:bg-cyan/10 glow-box"
          :class="{ 'opacity-40 pointer-events-none': wasmState !== 'ready' }"
        >
          {{ t('file') }}
          <input type="file" accept="audio/*" class="hidden" @change="onFile" />
        </label>
        <button
          v-if="analysis && !analysis.error"
          class="rounded border border-white/30 px-4 py-2 font-bold tracking-wider text-txt transition hover:bg-white/10 glow-box"
          :disabled="exporting"
          @click="savePng"
        >
          {{ exporting ? t('preparing') : t('png') }}
        </button>
        <span class="text-xs text-dim">&gt;_ {{ status }}</span>
      </section>

      <section class="flex flex-wrap items-center gap-2 rounded border border-line bg-panel p-3">
        <span class="text-[10px] uppercase tracking-widest text-dim">{{ t('analysisLang') }}</span>
        <div class="flex gap-1">
          <button
            v-for="l in ['tr', 'de', 'en']"
            :key="l"
            class="rounded border px-2 py-1 text-[10px] uppercase tracking-widest transition"
            :class="targetLang === l ? 'border-neon/60 text-neon glow-box' : 'border-line text-dim hover:text-txt'"
            @click="setTarget(l)"
          >
            {{ l }}
          </button>
        </div>
        <span class="text-[11px] text-dim">→ {{ t('analysisLangHint') }}</span>
      </section>

      <section v-if="recState === 'recording'" class="rounded border border-line bg-panel p-4">
        <div class="text-[10px] uppercase tracking-widest text-dim">{{ t('liveLabel') }} · {{ recSeconds.toFixed(1) }}s</div>
        <LiveCanvas :analyser="analyser" />
      </section>

      <section v-else-if="analysis && !analysis.error" class="flex flex-col gap-6">
        <div ref="exportRef" class="flex flex-col gap-6 p-1">
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
            <div v-for="c in statCards" :key="c.label" class="rounded border border-line bg-panel p-3">
              <div class="text-[10px] uppercase tracking-widest text-dim">{{ c.label }}</div>
              <div class="mt-1 text-lg" :class="c.cls">{{ c.value }}<span v-if="c.unit" class="text-xs"> {{ c.unit }}</span></div>
            </div>
          </div>

          <section v-if="autoNotes.length" class="rounded border border-cyan/40 bg-panel p-4">
            <div class="text-[10px] uppercase tracking-widest text-cyan">{{ t('autoComment') }}</div>
            <ul class="mt-2 space-y-1 text-sm">
              <li v-for="n in autoNotes" :key="n.text" :class="n.level === 'warn' ? 'text-cyan' : 'text-neon'">
                ▸ {{ n.text }}
              </li>
            </ul>
            <p class="mt-2 text-[10px] text-dim">{{ t('autoNote') }}</p>
          </section>

          <section class="rounded border border-line bg-panel p-3">
            <div class="flex flex-wrap items-center gap-3">
              <span class="text-[10px] uppercase tracking-widest text-dim">{{ t('compareTitle') }}</span>
              <button
                class="rounded border border-gold/60 px-3 py-1 text-xs font-bold tracking-wider text-gold transition hover:bg-gold/10 glow-box"
                @click="stored = analysis"
              >
                {{ t('compareStore') }}
              </button>
              <button
                v-if="stored"
                class="rounded border border-line px-3 py-1 text-xs tracking-wider text-dim transition hover:text-txt"
                @click="stored = null"
              >
                {{ t('compareClear') }}
              </button>
              <span v-if="stored" class="text-[11px] text-dim">↻ {{ t('compareHint') }}</span>
            </div>
          </section>

          <AnalysisCanvas :analysis="analysis" :compare="stored" />
        </div>

        <p class="text-center text-[11px] text-dim">
          {{ t('footer') }}
        </p>
      </section>

      <section v-else class="rounded border border-line bg-panel p-6 text-center text-dim">
        <div class="text-4xl">▚▞</div>
        <p class="mt-2 text-sm">{{ t('emptyTitle') }}</p>
      </section>
    </template>
  </div>
</template>
