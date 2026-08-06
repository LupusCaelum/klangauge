//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/LupusCaelum/klangauge/core"
)

// registerAnalyzer, tarayıcıya "klangaugeAnalyze" adında bir fonksiyon verir.
// JS tarafı şöyle çağırır:
//
//	var samples = new Float32Array([...])  // -1..1 arası örnekler
//	var result = klangaugeAnalyze(samples, 44100)  // JSON string
func registerAnalyzer() {
	js.Global().Set("klangaugeAnalyze", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) < 2 {
			return `{"error":"klangaugeAnalyze(samples, sampleRate) gerekli"}`
		}
		samples := args[0]
		sampleRate := args[1].Int()
		if samples.Type() != js.TypeObject {
			return `{"error":"samples bir Float32Array olmalı"}`
		}

		n := samples.Get("length").Int()
		floats := make([]float64, n)
		for i := 0; i < n; i++ {
			floats[i] = float64(samples.Index(i).Float())
		}

		result := core.Analyze(floats, sampleRate)
		out, err := json.Marshal(result)
		if err != nil {
			return `{"error":"analiz sonucu serileştirilemedi"}`
		}
		return string(out)
	}))
}

// registerSyllabifier, tarayıcıya "klangaugeSyllabify" adında bir fonksiyon verir.
// JS tarafı şöyle çağırır:
//
//	var result = klangaugeSyllabify("verstehen", "de")  // JSON string
func registerSyllabifier() {
	js.Global().Set("klangaugeSyllabify", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) < 2 {
			return `{"error":"klangaugeSyllabify(text, lang) gerekli"}`
		}
		text := args[0].String()
		lang := args[1].String()

		result := core.SyllabifyWord(lang, text)
		if result == nil {
			return `{"error":"metin boş"}`
		}
		out, err := json.Marshal(result)
		if err != nil {
			return `{"error":"heceleme sonucu serileştirilemedi"}`
		}
		return string(out)
	}))
}

func main() {
	c := make(chan struct{}, 0)
	registerAnalyzer()
	registerSyllabifier()
	// Ana fonksiyon dönmemeli; aksi halde WASM runtime kapanır.
	<-c
}
