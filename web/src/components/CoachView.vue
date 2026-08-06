<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { loadWasm, syllabify } from '../wasm.js'
import { useRecorder } from '../useRecorder.js'
import { evaluateCoach, formatExpected } from '../coach.js'
import { t, lang } from '../i18n.js'
import AnalysisCanvas from './AnalysisCanvas.vue'
import LiveCanvas from './LiveCanvas.vue'

const targetLang = ref('tr')
const word = ref('')
const expected = ref(null)
const analysis = ref(null)
const wasmReady = ref(false)
const wasmErr = ref('')

const { recState, recError, recSeconds, analyser, toggleRecord, loadFile } = useRecorder((a) => {
  analysis.value = a
})

function onFile(event) {
  const file = event.target.files?.[0]
  if (!file) return
  loadFile(file).catch(() => {})
  event.target.value = ''
}

function updateExpected() {
  const w = word.value.trim()
  if (!w || !wasmReady.value) {
    expected.value = null
    return
  }
  try {
    expected.value = syllabify(w, targetLang.value)
  } catch {
    expected.value = null
  }
}

watch([word, targetLang], updateExpected)

onMounted(() => {
  loadWasm()
    .then(() => {
      wasmReady.value = true
      updateExpected()
    })
    .catch((e) => (wasmErr.value = e.message))
})

const expectedParts = computed(() => formatExpected(expected.value))

const stressLabel = computed(() => {
  const e = expected.value
  if (!e || e.syllables.length === 0) return ''
  let base
  if (e.stress < 0) base = t('coachStressUnknown')
  else if (e.lang === 'tr') base = t('coachStressFinal')
  else base = t('coachStressIndex').replace('{n}', e.stress + 1)
  return e.approx ? base + ' ' + t('coachApprox') : base
})

const coachNotes = computed(() => evaluateCoach(expected.value, analysis.value, lang.value))

const status = computed(() => {
  if (recState.value === 'recording') return t('stRecording').replace('{s}', recSeconds.value.toFixed(1))
  if (recState.value === 'analyzing') return t('stAnalyzing')
  if (recState.value === 'error') return `${t('stError')}: ${recError.value}`
  if (wasmErr.value) return `${t('stWasmError')}: ${wasmErr.value}`
  return ''
})
</script>

<template>
  <div class="flex flex-col gap-5">
    <section class="rounded border border-line bg-panel p-4">
      <h2 class="glow-green text-sm font-bold uppercase tracking-widest text-neon">{{ t('coachTitle') }}</h2>
      <p class="mt-1 text-xs text-dim">{{ t('coachIntro') }}</p>

      <div class="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center">
        <div class="flex items-center gap-2">
          <span class="text-[10px] uppercase tracking-widest text-dim">{{ t('coachTargetLang') }}</span>
          <div class="flex gap-1">
            <button
              v-for="l in ['tr', 'de', 'en']"
              :key="l"
              class="rounded border px-2 py-1 text-[10px] uppercase tracking-widest transition"
              :class="targetLang === l ? 'border-cyan/60 text-cyan glow-box' : 'border-line text-dim hover:text-txt'"
              @click="targetLang = l"
            >
              {{ l }}
            </button>
          </div>
        </div>
        <input
          v-model="word"
          type="text"
          class="flex-1 rounded border border-line bg-black/40 px-3 py-2 font-mono text-sm text-txt outline-none transition focus:border-cyan/60"
          :placeholder="t('coachWordPlaceholder')"
          @keyup.enter="toggleRecord"
        />
      </div>

      <div v-if="expected" class="mt-3 rounded border border-line/60 p-3">
        <div class="text-[10px] uppercase tracking-widest text-dim">{{ t('coachExpected') }}</div>
        <div class="mt-1 font-mono text-lg">
          <span v-for="(p, i) in expectedParts" :key="i">
            <span :class="p.stressed ? 'font-bold text-cyan glow-cyan' : 'text-txt'">{{ p.text }}</span>
            <span v-if="i < expectedParts.length - 1" class="text-dim">-</span>
          </span>
        </div>
        <div class="mt-1 text-xs text-dim">
          {{ t('coachSyllables').replace('{n}', expected.syllables.length) }} · {{ stressLabel }}
        </div>
      </div>
      <p v-else class="mt-3 text-xs text-dim">{{ t('coachWordPlaceholder') }}</p>
    </section>

    <section class="flex flex-wrap items-center gap-3">
      <button
        class="rounded border px-4 py-2 font-bold tracking-wider transition"
        :class="
          recState === 'recording'
            ? 'border-mag text-mag glow-box'
            : 'border-neon text-neon hover:bg-neon/10 glow-box'
        "
        :disabled="!wasmReady || !word.trim() || recState === 'analyzing'"
        @click="toggleRecord"
      >
        {{ recState === 'recording' ? t('stop') : t('record') }}
      </button>
      <label
        class="cursor-pointer rounded border border-cyan px-4 py-2 font-bold tracking-wider text-cyan transition hover:bg-cyan/10 glow-box"
        :class="{ 'opacity-40 pointer-events-none': !wasmReady }"
      >
        {{ t('file') }}
        <input type="file" accept="audio/*" class="hidden" @change="onFile" />
      </label>
      <span v-if="status" class="text-xs text-dim">&gt;_ {{ status }}</span>
    </section>

    <section v-if="analysis && !analysis.error" class="flex flex-col gap-5">
      <div class="rounded border border-cyan/40 bg-panel p-4">
        <div class="text-[10px] uppercase tracking-widest text-cyan">{{ t('coachDetected') }}</div>
        <ul class="mt-2 space-y-1 text-sm">
          <li v-for="n in coachNotes" :key="n.text" :class="n.level === 'warn' ? 'text-cyan' : 'text-neon'">
            ▸ {{ n.text }}
          </li>
        </ul>
      </div>
      <AnalysisCanvas :analysis="analysis" />
    </section>
    <section v-else-if="recState === 'recording'" class="rounded border border-line bg-panel p-4">
      <div class="text-[10px] uppercase tracking-widest text-dim">{{ t('liveLabel') }} · {{ recSeconds.toFixed(1) }}s</div>
      <LiveCanvas :analyser="analyser" />
    </section>
    <section v-else class="rounded border border-line bg-panel p-6 text-center text-xs text-dim">
      {{ t('coachEmpty') }}
    </section>
  </div>
</template>
