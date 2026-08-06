module github.com/LupusCaelum/klangauge

go 1.26.5

require (
	github.com/FreibergVlad/go-yinfft v0.0.0-20251124124437-4470139c2a32
	github.com/LupusCaelum/syllabifier v0.0.0-00010101000000-000000000000
	github.com/mjibson/go-dsp v0.0.0-20260128111154-6db759bd4208
)

require github.com/madelynnblue/go-dsp v1.0.0 // indirect

// syllabifier henüz yayınlanmadı; yerel kopyadan bağlanıyoruz.
replace github.com/LupusCaelum/syllabifier => /home/umut/Projects/dev/Frontend/syllabifier
