// klangauge.wasm yükleyici.
// Go'nun wasm_exec.js runtime'ı + derlenmiş WASM'i yükler ve
// JS tarafından çağrılabilir `klangaugeAnalyze` fonksiyonunu hazırlar.

let analyzeFn = null
let syllabifyFn = null
let loading = null

function loadScript(src) {
  return new Promise((resolve, reject) => {
    const s = document.createElement('script')
    s.src = src
    s.onload = resolve
    s.onerror = () => reject(new Error(`script yüklenemedi: ${src}`))
    document.head.appendChild(s)
  })
}

export function loadWasm() {
  if (loading) return loading
  loading = (async () => {
    await loadScript('./wasm_exec.js')
    const go = new globalThis.Go()
    const resp = await fetch('./klangauge.wasm')
    if (!resp.ok) throw new Error(`WASM yüklenemedi (${resp.status})`)
    const buf = await resp.arrayBuffer()
    const { instance } = await WebAssembly.instantiate(buf, go.importObject)
    go.run(instance) // main() bloke olur; köprü fonksiyonları kayıtlı kalır
    analyzeFn = globalThis.klangaugeAnalyze
    syllabifyFn = globalThis.klangaugeSyllabify
    if (!analyzeFn) throw new Error('klangaugeAnalyze global değil')
    if (!syllabifyFn) throw new Error('klangaugeSyllabify global değil')
    return analyzeFn
  })()
  return loading
}

// samples: Float32Array (-1..1), sampleRate: int → Analysis (JSON → obje)
export function analyze(samples, sampleRate) {
  if (!analyzeFn) throw new Error('WASM henüz yüklenmedi')
  return JSON.parse(analyzeFn(samples, sampleRate))
}

// text: kelime, lang: "tr"|"de"|"en" → ExpectedWord (JSON → obje)
export function syllabify(text, lang) {
  if (!syllabifyFn) throw new Error('WASM henüz yüklenmedi')
  return JSON.parse(syllabifyFn(text, lang))
}
