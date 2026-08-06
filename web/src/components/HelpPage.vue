<script setup>
import { computed } from 'vue'
import { helpSections, helpMetrics } from '../explainer.js'
import { lang, t, legend } from '../i18n.js'

const sections = computed(() => helpSections[lang.value] ?? helpSections.tr)
const metrics = computed(() => helpMetrics[lang.value] ?? helpMetrics.tr)
const lg = computed(() => legend[lang.value] ?? legend.tr)
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Efsane -->
    <section class="rounded border border-line bg-panel p-4">
      <h2 class="glow-green text-sm font-bold uppercase tracking-widest text-neon">{{ t('legendTitle') }}</h2>
      <p class="mt-1 text-xs text-dim">{{ t('legendIntro') }}</p>
      <ul class="mt-3 space-y-2">
        <li v-for="l in lg" :key="l.item" class="flex flex-col gap-1 text-sm sm:flex-row sm:items-center sm:gap-3">
          <span class="inline-block h-3 w-6 shrink-0 rounded" :class="l.swatch" />
          <span class="w-44 shrink-0 font-bold text-txt">{{ l.item }}</span>
          <span class="text-dim">{{ l.where }} · {{ l.what }}</span>
        </li>
      </ul>
    </section>

    <!-- Grafik bölümleri -->
    <section v-for="s in sections" :key="s.id" class="rounded border border-line bg-panel p-4">
      <h3 class="glow-cyan text-sm font-bold uppercase tracking-widest text-cyan">{{ s.title }}</h3>
      <p class="mt-1 text-xs text-dim">▸ {{ s.visual }}</p>
      <p class="mt-2 text-sm leading-relaxed">{{ s.body }}</p>
      <p class="mt-2 text-sm leading-relaxed text-neon">{{ t('sectionsTip') }} {{ s.tip }}</p>
    </section>

    <!-- Metrikler -->
    <section class="rounded border border-line bg-panel p-4">
      <h2 class="glow-green text-sm font-bold uppercase tracking-widest text-neon">{{ t('metricsTitle') }}</h2>
      <p class="mt-1 text-xs text-dim">{{ t('metricsIntro') }}</p>
      <div class="mt-3 grid gap-3 sm:grid-cols-2">
        <div v-for="m in metrics" :key="m.id" class="rounded border border-line/60 p-3">
          <div class="text-xs font-bold uppercase tracking-wider text-txt">{{ m.name }}</div>
          <p class="mt-1 text-sm leading-relaxed">{{ m.body }}</p>
          <p class="mt-2 text-xs leading-relaxed text-cyan">{{ t('metricsTarget') }} {{ m.good }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.bg-spectro {
  background: linear-gradient(90deg, #0c1310, #006e5b, #00e5ff, #ffffff);
}
</style>
