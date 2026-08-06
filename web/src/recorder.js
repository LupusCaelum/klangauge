// Mikrofon kaydı ve dosya çözme yardımcıları.
// Çıktı: { samples: Float32Array, sampleRate: number }
// Kayıt sırasında canlı seviye/grafik için bir AnalyserNode da döner.

export function startRecording() {
  return navigator.mediaDevices.getUserMedia({ audio: true }).then((stream) => {
    const ac = new AudioContext()
    const source = ac.createMediaStreamSource(stream)
    const analyser = ac.createAnalyser()
    analyser.fftSize = 2048
    source.connect(analyser)
    ac.resume()

    const mediaRecorder = new MediaRecorder(stream)
    const chunks = []
    mediaRecorder.ondataavailable = (e) => {
      if (e.data.size > 0) chunks.push(e.data)
    }
    const done = new Promise((resolve, reject) => {
      mediaRecorder.onstop = async () => {
        stream.getTracks().forEach((t) => t.stop())
        ac.close()
        try {
          resolve(await decodeBlob(new Blob(chunks, { type: mediaRecorder.mimeType })))
        } catch (err) {
          reject(err)
        }
      }
      mediaRecorder.onerror = () => reject(new Error('kayıt hatası'))
    })
    mediaRecorder.start()
    return { stop: () => mediaRecorder.stop(), done, analyser }
  })
}

// Analyser'dan anlık seviyeyi dB olarak okur (canlı gösterge için).
export function readLevel(analyser) {
  const buf = new Float32Array(analyser.fftSize)
  analyser.getFloatTimeDomainData(buf)
  let sum = 0
  for (let i = 0; i < buf.length; i++) sum += buf[i] * buf[i]
  const rms = Math.sqrt(sum / buf.length)
  if (rms <= 1e-9) return -90
  return 20 * Math.log10(rms)
}

export function decodeFile(file) {
  return file.arrayBuffer().then((buf) => decodeBlob(new Blob([buf], { type: file.type })))
}

function decodeBlob(blob) {
  return blob.arrayBuffer().then((buf) => {
    const ac = new AudioContext()
    return ac.decodeAudioData(buf).then(
      (audioBuf) => {
        ac.close()
        return { samples: audioBuf.getChannelData(0), sampleRate: audioBuf.sampleRate }
      },
      (err) => {
        ac.close()
        throw new Error('ses çözülemedi: ' + err.message)
      },
    )
  })
}
