// Ortak kayıt durumu composable'ı.
// Hem analiz görünümü (App.vue) hem koç görünümü (CoachView.vue) kullanır.
// onAnalysis(result, source): kayıt/dosya başarıyla çözülünce çağrılır.
import { ref } from 'vue'
import { analyze } from './wasm.js'
import { startRecording, decodeFile } from './recorder.js'

export function useRecorder(onAnalysis) {
  const recState = ref('idle') // idle | recording | analyzing | error
  const recError = ref('')
  const recSeconds = ref(0)
  const analyser = ref(null) // kayıt sırasında canlı seviye/gösterge için

  let recorder = null
  let timer = null

  function clearTimer() {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  async function run(samples, sampleRate) {
    const result = analyze(samples, sampleRate)
    if (result.error) throw new Error(result.error)
    onAnalysis(result)
  }

  async function toggleRecord() {
    if (recState.value === 'recording') {
      recorder?.stop()
      return
    }
    recError.value = ''
    try {
      recorder = await startRecording()
      analyser.value = recorder.analyser ?? null
      recState.value = 'recording'
      recSeconds.value = 0
      timer = setInterval(() => recSeconds.value++, 100)
      const { samples, sampleRate } = await recorder.done
      clearTimer()
      analyser.value = null
      recState.value = 'analyzing'
      await run(samples, sampleRate)
      recState.value = 'idle'
    } catch (e) {
      clearTimer()
      analyser.value = null
      recState.value = 'error'
      recError.value = e.message
    }
  }

  async function loadFile(file) {
    recError.value = ''
    try {
      recState.value = 'analyzing'
      const { samples, sampleRate } = await decodeFile(file)
      await run(samples, sampleRate)
      recState.value = 'idle'
    } catch (e) {
      recState.value = 'error'
      recError.value = e.message
    }
  }

  return { recState, recError, recSeconds, analyser, toggleRecord, loadFile }
}
