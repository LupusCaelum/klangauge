// Klangauge açıklama içeriği ve otomatik yorum motoru.
// Tüm metinler üç dilde (tr/de/en) — yapı dilden bağımsız sayılara,
// dilden bağımlı yorumlara ayrılır. İpucu: hece-zamanlı diller (tr)
// doğal olarak daha fazla hece/dk üretir; bantlar ona göre kayar.

export const helpSections = {
  tr: [
    {
      id: 'waveform',
      title: 'Dalga biçimi',
      visual: 'Alttaki grafik: dikey yeşil çizgiler',
      body: 'Sesinin ham görüntüsü. Çizgiler kalınlaştıkça ses güçlenir (genlik büyür), inceldikçe sessizleşir. Sessizlikte çizgiler orta çizgiye iner.',
      tip: 'Kelimelerin çizgilerde net iniş-çıkışlar yaratması doğal. Hep düz çizgi görüyorsan sesin donuk ya da çok kısık olabilir.',
    },
    {
      id: 'intensity',
      title: 'Şiddet eğrisi',
      visual: 'Dalga biçiminin üstünde uzanan cyan (turkuaz) çizgi',
      body: 'Sesin desibel (dB) cinsinden gücü. Tepe noktaları hecelere karşılık gelir: her ünlü ses ("a, e, o…") enerjinin en yüksek olduğu andır.',
      tip: 'Eğri tekdüze ve düzse vurgu yapmıyorsun demektir. Vurgulu okumada tepe noktaları belirginleşir.',
    },
    {
      id: 'spectrogram',
      title: 'Spektrogram',
      visual: 'Üstteki renkli grafik — dikey: frekans (Hz), yatay: zaman, renk: şiddet',
      body: 'Sesin renkli parmak izi. Dikey eksen sesin yüksekliğini (Hz) gösterir; renk gücü verir: koyu = zayıf, yeşil-cyan = orta, beyaz = çok güçlü. Ünlüler geniş yatay bantlar, ünsüzler dar ve dikey parlamalar yapar.',
      tip: 'Aynı kelimeyi iki kez söyleyip desenleri karşılaştır. Farklı sesler farklı desen çizer; deseni görerek sesini şekillendirmeye çalış.',
    },
    {
      id: 'pitch',
      title: 'Perde çizgisi (pembe)',
      visual: 'Spektrogramın üzerinde uzanan ince pembe çizgi',
      body: 'Sesinin temel frekansı — kulakla "yüksek / alçak" duyduğun şey. Çizgi yükseliyorsa sesin yükseliyor (soru tonu), alçalıyorsa sesin düşüyor. Çizgi yalnızca sesli konuşma algılanan anlarda çizilir: ünlülerde ve sesli ünsüzlerde görünür, sessiz harflerde (s, f, k…) ve duraklarda kaybolur — bu yüzden çizgi bölük pörçük görünebilir, bu normaldir.',
      tip: 'Soru sorarken çizginin yükselmesi gerekir. Monoton okuyorsan çizgi neredeyse düz kalır; vurgu yaparken belirgin bir tepe oluştur.',
    },
  ],
  de: [
    {
      id: 'waveform',
      title: 'Wellenform',
      visual: 'Der untere Graph: vertikale grüne Linien',
      body: 'Das Rohbild deiner Stimme. Dickere Linien = lauter (größere Amplitude), dünnere = leiser. In Pausen liegen die Linien auf der Mittellinie.',
      tip: 'Es ist normal, dass Wörter klare Aufs und Abs erzeugen. Wirkt alles ganz flach, ist dein Ton evtl. monoton oder zu leise.',
    },
    {
      id: 'intensity',
      title: 'Intensitätskurve',
      visual: 'Die cyanfarbene Linie über der Wellenform',
      body: 'Die Lautstärke in Dezibel (dB). Die Gipfel entsprechen Silben: jeder Vokal ("a, e, o…") ist ein Moment maximaler Energie.',
      tip: 'Verläuft die Kurve flach und eintönig, betonst du nicht. Beim betonten Sprechen werden die Gipfel deutlich.',
    },
    {
      id: 'spectrogram',
      title: 'Spektrogramm',
      visual: 'Das obere Farbdiagramm — vertikal: Frequenz (Hz), horizontal: Zeit, Farbe: Intensität',
      body: 'Der farbige Fingerabdruck deines Klangs. Die senkrechte Achse zeigt die Tonhöhe (Hz), die Farbe die Stärke: dunkel = schwach, grün-cyan = mittel, weiß = sehr stark. Vokale bilden breite horizontale Bänder, Konsonanten schmale vertikale Blitze.',
      tip: 'Sprich dasselbe Wort zweimal und vergleiche die Muster. Verschiedene Laute zeichnen verschiedene Muster; versuche, deinen Klang dem Muster anzupassen.',
    },
    {
      id: 'pitch',
      title: 'Tonhöhenlinie (rosa)',
      visual: 'Die dünne rosa Linie über dem Spektrogramm',
      body: 'Deine Grundfrequenz — das, was du als "hoch / tief" hörst. Steigt die Linie, wird deine Stimme höher (Fragesatz), fällt sie, tiefer. Sie erscheint nur in stimmhaften Momenten: bei Vokalen und stimmhaften Konsonanten; bei stimmlosen Lauten (s, f, k…) und Pausen verschwindet sie — deshalb wirkt sie unterbrochen, das ist normal.',
      tip: 'Bei einer Frage muss die Linie steigen. Beim monotonen Lesen bleibt sie fast flach; beim Betonen bildet sich ein deutlicher Gipfel.',
    },
  ],
  en: [
    {
      id: 'waveform',
      title: 'Waveform',
      visual: 'The lower graph: vertical green lines',
      body: 'A raw image of your voice. Thicker lines mean louder (bigger amplitude), thinner mean quieter. In silence the lines sit on the middle line.',
      tip: 'Words should produce clear up-and-down movement. If it looks flat throughout, your delivery may be dull or too quiet.',
    },
    {
      id: 'intensity',
      title: 'Intensity curve',
      visual: 'The cyan line running above the waveform',
      body: 'The power of your voice in decibels (dB). Peaks correspond to syllables: every vowel ("a, e, o…") is a moment of peak energy.',
      tip: 'If the curve is flat and uniform, you are not stressing. With stress, the peaks become clearly visible.',
    },
    {
      id: 'spectrogram',
      title: 'Spectrogram',
      visual: 'The colored graph at top — vertical: frequency (Hz), horizontal: time, color: intensity',
      body: 'The colorful fingerprint of your sound. The vertical axis shows pitch (Hz); the color shows strength: dark = weak, green-cyan = medium, white = very strong. Vowels make broad horizontal bands, consonants make narrow vertical flashes.',
      tip: 'Say the same word twice and compare patterns. Different sounds draw different patterns; try to shape your sound to match the pattern.',
    },
    {
      id: 'pitch',
      title: 'Pitch line (pink)',
      visual: 'The thin pink line over the spectrogram',
      body: 'Your fundamental frequency — what you hear as "high / low". Rising line = rising voice (question intonation), falling = lowering. It is drawn only in voiced moments: vowels and voiced consonants; on voiceless sounds (s, f, k…) and pauses it disappears — so the line looks broken, that is normal.',
      tip: 'Asking a question should make the line rise. If you read monotonously the line stays flat; stress creates a clear peak.',
    },
  ],
}

export const helpMetrics = {
  tr: [
    {
      id: 'duration',
      name: 'Süre',
      body: 'Kaydın toplam uzunluğu. Ne kadar uzun konuştuğunu gösterir.',
      good: 'Analiz için 1-3 saniyelik tek bir kelime/cümle veya 10-30 saniyelik bir okuma uygundur.',
    },
    {
      id: 'syllables',
      name: 'Hece sayısı',
      body: 'Şiddet eğrisindeki enerji tepelerinden tahmin edilir. Söylediğin metindeki hece sayısına yakın olması beklenir.',
      good: 'Kelimeleri net söylüyorsan sayı, beklenen hece sayısına yakın çıkar; çok düşükse heceler yutuluyor olabilir.',
    },
    {
      id: 'speechRate',
      name: 'Konuşma hızı',
      body: 'Dakikadaki hece sayısı. Yavaş < 105, dengeli 105-185, hızlı 185-250, çok hızlı > 250 hc/dk (hece-zamanlı diller, yaklaşık).',
      good: 'Yeni bir dil öğrenirken dengeli tempo idealdir. Çok hızlı okumak yutmaya, çok yavaş kalmak takılmaya işaret edebilir.',
    },
    {
      id: 'articulation',
      name: 'Artikülasyon hızı',
      body: 'Sadece konuşulan anlarda saniyede söylenen hece. Duraklamalar hesaba katılmaz — saf artikülasyon hızıdır.',
      good: 'Artikülasyon yüksek ama konuşma hızı düşükse çok duraklıyorsun demektir; akıcılık için durakları kısalt.',
    },
    {
      id: 'voicedRatio',
      name: 'Sesli oranı',
      body: 'Konuşma geçen sürenin toplam süreye oranı. Düşük değer = uzun sessizlikler, çok yüksek = hiç durmuyorsun.',
      good: 'Doğal konuşmada %60-85 arası tipiktir. Çok düşükse durakları azalt, çok yüksekse cümle aralarına küçük nefes molaları koy.',
    },
    {
      id: 'meanPitch',
      name: 'Ortalama perde',
      body: 'Sesinin ortalama temel frekansı. Erkek sesi genelde 85-180 Hz, kadın sesi 165-255 Hz civarındadır; kişisel fark normaldir.',
      good: '"Doğru" perde yoktur — bu senin doğal ses yüksekliğindir. Önemli olan sabit değil, kullanılan aralıktır (perde aralığı).',
    },
    {
      id: 'pitchRange',
      name: 'Perde aralığı',
      body: 'Konuşmadaki en düşük ve en yüksek perde farkı. Melodi ve vurgunun canlılığını gösterir.',
      good: 'Monoton okumada < 30 Hz kalır. Vurgulu, anlamlı okumada ses doğal olarak 60 Hz üzeri hareket eder.',
    },
    {
      id: 'meanIntensity',
      name: 'Ortalama şiddet',
      body: 'Konuştuğun anların ortalama ses seviyesi (dB) — sessizlikler hesaba katılmaz. Çok düşük = kısık, çok yüksek = mikrofona çok yakın ya da bağırarak.',
      good: 'Net ve rahat bir seviye hedeftir. Normal ev/oda sesi genelde -35 dB üzerindedir; kısık konuşuyorsan mikrofona yaklaşmayı dene.',
    },
  ],
  de: [
    {
      id: 'duration',
      name: 'Dauer',
      body: 'Gesamtlänge der Aufnahme. Zeigt, wie lange du gesprochen hast.',
      good: 'Für die Analyse eignen sich 1-3 Sekunden für ein einzelnes Wort/einen Satz oder 10-30 Sekunden für eine Leseprobe.',
    },
    {
      id: 'syllables',
      name: 'Silbenzahl',
      body: 'Wird aus den Energie-Gipfeln der Intensitätskurve geschätzt. Sie sollte nahe der Silbenzahl deines Texts liegen.',
      good: 'Sprichst du die Wörter klar, liegt die Zahl nahe der erwarteten Silbenzahl; ist sie viel niedriger, könntest du Silben verschlucken.',
    },
    {
      id: 'speechRate',
      name: 'Sprechtempo',
      body: 'Silben pro Minute. Langsam < 90, ausgewogen 90-160, schnell 160-220, sehr schnell > 220 Sil./min (vokalbetonte... für betontes Sprechen, ca.).',
      good: 'Beim Sprachenlernen ist ein ausgewogenes Tempo ideal. Zu schnelles Lesen führt zu verschluckten Silben, zu langsames zum Stocken.',
    },
    {
      id: 'articulation',
      name: 'Artikulation',
      body: 'Silben pro Sekunde, nur während des Sprechens gezählt — Pausen zählen nicht. Es misst, wie schnell deine Zunge arbeitet.',
      good: 'Ist die Artikulation hoch, das Sprechtempo aber niedrig, machst du viele Pausen — verkürze sie für mehr Fluss.',
    },
    {
      id: 'voicedRatio',
      name: 'Stimmanteil',
      body: 'Anteil der Sprechzeit an der Gesamtdauer. Niedrig = lange Pausen, sehr hoch = du hältst nie inne.',
      good: '60-85 % ist bei natürlicher Rede typisch. Zu niedrig: Pausen verkürzen; zu hoch: kleine Atempausen zwischen Sätzen einlegen.',
    },
    {
      id: 'meanPitch',
      name: 'Mittlere Tonhöhe',
      body: 'Deine durchschnittliche Grundfrequenz. Männer meist 85-180 Hz, Frauen 165-255 Hz; individuelle Unterschiede sind normal.',
      good: 'Es gibt keine "richtige" Tonhöhe — sie ist deine natürliche Stimmlage. Wichtig ist nicht der Wert, sondern die genutzte Spannweite (Tonhöhenumfang).',
    },
    {
      id: 'pitchRange',
      name: 'Tonhöhenumfang',
      body: 'Differenz zwischen tiefster und höchster Tonhöhe. Zeigt, wie lebendig Melodie und Betonung sind.',
      good: 'Monotones Lesen bleibt unter 30 Hz. Bei betontem, ausdrucksvollem Sprechen bewegt sich die Stimme natürlich über 60 Hz.',
    },
    {
      id: 'meanIntensity',
      name: 'Mittlere Lautstärke',
      body: 'Durchschnittliche Lautstärke (dB) nur in den Sprechmomenten — Pausen zählen nicht. Sehr niedrig = leise, sehr hoch = sehr nah am Mikrofon oder laut.',
      good: 'Ziel ist ein klares, bequemes Niveau. Normale Zimmerlautstärke liegt meist über -35 dB; bist du leise, nähere dich dem Mikrofon.',
    },
  ],
  en: [
    {
      id: 'duration',
      name: 'Duration',
      body: 'Total length of the recording — how long you spoke.',
      good: 'For analysis, 1-3 seconds of a single word/sentence, or 10-30 seconds of reading, works well.',
    },
    {
      id: 'syllables',
      name: 'Syllable count',
      body: 'Estimated from energy peaks in the intensity curve. It should be close to the syllable count of what you said.',
      good: 'If you articulate clearly, the count is close to the expected syllables; much lower, and syllables may be swallowed.',
    },
    {
      id: 'speechRate',
      name: 'Speech rate',
      body: 'Syllables per minute. Slow < 90, balanced 90-160, fast 160-220, very fast > 220 syl/min (for stress-timed languages, approx).',
      good: 'For language learners a balanced tempo is ideal. Reading too fast causes swallowing, too slow causes stumbling.',
    },
    {
      id: 'articulation',
      name: 'Articulation rate',
      body: 'Syllables per second, counted only during speech — pauses are excluded. It measures how fast your articulators move.',
      good: 'If articulation is high but speech rate low, you pause a lot — shorten the pauses for fluency.',
    },
    {
      id: 'voicedRatio',
      name: 'Voiced ratio',
      body: 'Share of speaking time within the total duration. Low = long silences, very high = you never pause.',
      good: '60-85% is typical for natural speech. Too low: reduce pauses; too high: breathe briefly between sentences.',
    },
    {
      id: 'meanPitch',
      name: 'Mean pitch',
      body: 'Your average fundamental frequency. Male voices are typically 85-180 Hz, female 165-255 Hz; individual differences are normal.',
      good: 'There is no "right" pitch — it is your natural voice level. What matters is the range you use (pitch range).',
    },
    {
      id: 'pitchRange',
      name: 'Pitch range',
      body: 'Difference between the lowest and highest pitch. Shows how lively your melody and emphasis are.',
      good: 'Monotonous reading stays under 30 Hz. With expressive, stressed speech the voice naturally moves beyond 60 Hz.',
    },
    {
      id: 'meanIntensity',
      name: 'Mean intensity',
      body: 'Average loudness (dB) during speech moments only — silence is excluded. Very low = whispering, very high = too close to the mic or shouting.',
      good: 'A clear, comfortable level is the goal. Normal room voice is usually above -35 dB; if quiet, move closer to the mic.',
    },
  ],
}

// Otomatik yorum motoru (koçun ilk çekirdeği). Kural bazlıdır.
// level: 'ok' (yeşil) | 'warn' (cyan) — sonraki fazda "yapıcı geri bildirim"e evrilecek.
const rateBands = {
  // hece-zamanlı (syllable-timed): doğal tempo daha yüksek
  tr: { slow: 105, balanced: 185, fast: 250 },
  // vurgu-zamanlı (stress-timed): doğal tempo daha düşük
  de: { slow: 90, balanced: 160, fast: 220 },
  en: { slow: 90, balanced: 160, fast: 220 },
}

const artBands = {
  tr: { min: 1.8, max: 5.5 },
  de: { min: 1.5, max: 5.0 },
  en: { min: 1.5, max: 5.0 },
}

const msg = {
  tr: {
    empty: 'Veri yok — önce bir kayıt yap.',
    rateSlow: 'Konuşma hızın yavaş ve sakin — yeni öğrenenler için ideal. Akıcılık için tempoyu azar azar artırabilirsin.',
    rateBalanced: 'Konuşma hızın dengeli — çoğu öğrenme senaryosu için uygun bir tempo.',
    rateFast: 'Konuşma hızın biraz yüksek. Heceleri hafifçe uzatarak anlaşılırlığı artırmayı dene.',
    rateVeryFast: 'Konuşma hızın çok yüksek — heceleri yutuyor olabilirsin. Bilinçli olarak yavaşla.',
    artSlow: 'Artikülasyonun ağır ve özenli — net ama biraz yavaş. Zamanla hızlanacak.',
    artNormal: 'Artikülasyon hızın doğal bir tempoda.',
    artFast: 'Artikülasyon çok hızlı — dil yeterince net şekillenmeden sesler peş peşe geliyor olabilir.',
    vrLow: 'Çok fazla sessizlik var (sesli oranı %25 altı). Kısa bir deneme cümlesiyle tekrar dene.',
    vrHigh: 'Neredeyse hiç duraklamıyorsun — cümle aralarına küçük nefes molaları koy.',
    vrOk: 'Sessizlik/konuşma dengesin doğal görünüyor.',
    prFlat: 'Perde aralığın çok dar — konuşman monoton. Vurgu yapacağın kelimelerde sesini yükseltip alçaltmayı dene.',
    prModerate: 'Melodik hareketin var; vurguları biraz daha büyüterek ifadeyi güçlendirebilirsin.',
    prRich: 'Zengin perde aralığı kullanıyorsun — vurgu ve ifade için çok iyi.',
    miLow: 'Konuşma anlarında ses seviyen düşük (kısık). Mikrofona yaklaşıp biraz daha gür konuşmayı dene.',
    miHigh: 'Ses seviyen çok yüksek — mikrofondan uzaklaş veya sesini alçalt.',
    allOk: 'Kayıt sağlıklı görünüyor — devam et!',
  },
  de: {
    empty: 'Keine Daten — erst eine Aufnahme machen.',
    rateSlow: 'Dein Sprechtempo ist langsam und ruhig — ideal zum Lernen. Steigere das Tempo für mehr Fluss nach und nach.',
    rateBalanced: 'Dein Sprechtempo ist ausgewogen — gut für fast alle Lern-Situationen.',
    rateFast: 'Dein Sprechtempo ist etwas hoch. Zieh die Silben leicht in die Länge, um verständlicher zu sein.',
    rateVeryFast: 'Dein Sprechtempo ist sehr hoch — du verschluckst womöglich Silben. Bremse bewusst ab.',
    artSlow: 'Deine Artikulation ist langsam und sorgfältig — deutlich, aber etwas zäh. Das wird mit der Zeit schneller.',
    artNormal: 'Deine Artikulation liegt in einem natürlichen Tempo.',
    artFast: 'Die Artikulation ist sehr schnell — die Laute kommen vielleicht schneller, als die Zunge sie formen kann.',
    vrLow: 'Sehr viel Stille (Stimmanteil unter 25 %). Versuch es mit einem kurzen Probiersatz noch einmal.',
    vrHigh: 'Du machst fast keine Pausen — setze kleine Atempausen zwischen die Sätze.',
    vrOk: 'Dein Gleichgewicht aus Pausen und Sprechen wirkt natürlich.',
    prFlat: 'Dein Tonhöhenumfang ist sehr eng — du klingst monoton. Heb und senk die Stimme bei betonten Wörtern.',
    prModerate: 'Du hast melodische Bewegung; vergrößere die Betonungen etwas, um mehr Ausdruck zu bekommen.',
    prRich: 'Du nutzt einen reichen Tonhöhenumfang — sehr gut für Betonung und Ausdruck.',
    miLow: 'In den Sprechmomenten ist deine Lautstärke niedrig (leise). Geh näher ans Mikrofon und sprich etwas kräftiger.',
    miHigh: 'Deine Lautstärke ist sehr hoch — geh vom Mikrofon weg oder sprich leiser.',
    allOk: 'Die Aufnahme sieht gesund aus — mach weiter!',
  },
  en: {
    empty: 'No data — record something first.',
    rateSlow: 'Your speech rate is slow and calm — ideal for beginners. Raise the tempo gradually for fluency.',
    rateBalanced: 'Your speech rate is balanced — a comfortable tempo for most learning scenarios.',
    rateFast: 'Your speech rate is a bit high. Try lengthening syllables slightly to stay clear.',
    rateVeryFast: 'Your speech rate is very high — you may be swallowing syllables. Slow down deliberately.',
    artSlow: 'Your articulation is slow and careful — clear but a little heavy. It will speed up with time.',
    artNormal: 'Your articulation is at a natural tempo.',
    artFast: 'Articulation is very fast — sounds may follow faster than your tongue can shape them.',
    vrLow: 'A lot of silence (voiced ratio under 25%). Try again with a short trial sentence.',
    vrHigh: 'You barely pause — take small breathing breaks between sentences.',
    vrOk: 'Your balance of silence and speech looks natural.',
    prFlat: 'Your pitch range is very narrow — you sound monotone. Try raising and lowering your voice on stressed words.',
    prModerate: 'You have melodic movement; make the stresses a bit bigger to strengthen expression.',
    prRich: 'You use a rich pitch range — great for stress and expression.',
    miLow: 'Your level is low (quiet) during speech. Move closer to the mic and speak a bit louder.',
    miHigh: 'Your level is very high — move away from the mic or lower your voice.',
    allOk: 'The recording looks healthy — keep going!',
  },
}

export function explainStats(stats, uiLang = 'tr', targetLang = uiLang) {
  const M = msg[uiLang] ?? msg.tr
  if (!stats) return [{ level: 'warn', text: M.empty }]

  const rb = rateBands[targetLang] ?? rateBands.tr
  const ab = artBands[targetLang] ?? artBands.tr
  const notes = []
  const sr = stats.speechRate
  const ar = stats.articulationRate
  const vr = stats.voicedRatio
  const pr = stats.pitchRange
  const mi = stats.meanIntensity

  if (sr > 0) {
    if (sr < rb.slow)
      notes.push({ level: 'ok', text: M.rateSlow })
    else if (sr < rb.balanced)
      notes.push({ level: 'ok', text: M.rateBalanced })
    else if (sr < rb.fast)
      notes.push({ level: 'warn', text: M.rateFast })
    else
      notes.push({ level: 'warn', text: M.rateVeryFast })
  }

  if (ar > 0) {
    if (ar < ab.min)
      notes.push({ level: 'ok', text: M.artSlow })
    else if (ar <= ab.max)
      notes.push({ level: 'ok', text: M.artNormal })
    else
      notes.push({ level: 'warn', text: M.artFast })
  }

  if (vr > 0) {
    if (vr < 0.25)
      notes.push({ level: 'warn', text: M.vrLow })
    else if (vr > 0.85)
      notes.push({ level: 'warn', text: M.vrHigh })
    else
      notes.push({ level: 'ok', text: M.vrOk })
  }

  if (pr > 0) {
    if (pr < 30)
      notes.push({ level: 'warn', text: M.prFlat })
    else if (pr < 70)
      notes.push({ level: 'ok', text: M.prModerate })
    else
      notes.push({ level: 'ok', text: M.prRich })
  }

  if (mi < -35)
    notes.push({ level: 'warn', text: M.miLow })
  else if (mi > -10)
    notes.push({ level: 'warn', text: M.miHigh })

  if (notes.length === 0) notes.push({ level: 'ok', text: M.allOk })
  return notes
}
