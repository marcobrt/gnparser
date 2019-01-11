package preprocess_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/extensions/table"
	. "gitlab.com/gogna/gnparser/preprocess"
)

var _ = Describe("Preprocess", func() {
	DescribeTable("NormalizeHybridChar",
		func(s string, expected string) {
			Expect(NormalizeHybridChar([]byte(s))).To(Equal([]byte(expected)))
		},
		Entry(
			"'×', no space at the start",
			"×Agropogon P. Fourn. 1934",
			"×Agropogon P. Fourn. 1934",
		),
		Entry(
			"'x', no space at the start",
			"xAgropogon P. Fourn. 1934",
			"×Agropogon P. Fourn. 1934",
		),
		Entry(
			"'X', no space at the start",
			"XAgropogon P. Fourn. 1934",
			"×Agropogon P. Fourn. 1934",
		),
		Entry(
			"'×', space at the start",
			"× Agropogon P. Fourn. 1934",
			"× Agropogon P. Fourn. 1934",
		),
		Entry(
			"'x', space at the start",
			"x Agropogon P. Fourn. 1934",
			"× Agropogon P. Fourn. 1934",
		),
		Entry(
			"'X', space at the start",
			"X Agropogon P. Fourn. 1934",
			"× Agropogon P. Fourn. 1934",
		),
		Entry(
			"'×', no space at species",
			"Mentha ×smithiana ",
			"Mentha ×smithiana ",
		),
		Entry(
			"'X', spaces at species",
			"Asplenium X inexpectatum",
			"Asplenium × inexpectatum",
		),
		Entry(
			"'x', spaces at species",
			"Salix x capreola Andersson",
			"Salix × capreola Andersson",
		),
		Entry(
			"'x', spaces in formula",
			"Asplenium rhizophyllum DC. x ruta-muraria E.L. Braun 1939",
			"Asplenium rhizophyllum DC. × ruta-muraria E.L. Braun 1939",
		),
		// This one is brittle!
		Entry(
			"'X', spaces in formula",
			"Arthopyrenia hyalospora Hall X Hydnellum scrobiculatum D.E. Stuntz",
			"Arthopyrenia hyalospora Hall × Hydnellum scrobiculatum D.E. Stuntz",
		),
		Entry(
			"'x', in the end",
			"Arthopyrenia hyalospora x",
			"Arthopyrenia hyalospora ×",
		),
	)

	DescribeTable("IsVirus",
		func(s string, itIs bool) {
			res := IsVirus([]byte(s))
			Expect(res).To(Equal(itIs))
		},
		Entry("No match", "Homo sapiens", false),
		Entry("Match word", "Arv1virus ", true),
		Entry("Match word", "Turtle herpesviruses", true),
		Entry("Match word", "Cre expression vector", true),
		Entry("Match word", "Abutilon mosaic vir. ICTV", true),
		Entry("Match word", "Aeromonas phage 65", true),
		Entry("Match word", "Apple scar skin viroid", true),
		Entry("Match word", "Agents of Spongiform Encephalopathies CWD prion Chronic wasting disease", true),
		Entry("Match word", "Phi h-like viruses", true),
		Entry("Match word", "Viroids", true),
		Entry("Match word", "Human rhinovirus A11", true),
		Entry("Match word", "Gossypium mustilinum symptomless alphasatellite", true),
		Entry("Match word", "Bemisia betasatellite LW-2014", true),
		Entry("Match word", "Intracisternal A-particles", true),
		Entry("Match word", "Uranotaenia sapphirina NPV", true),
		Entry("Match word", "Spodoptera frugiperda MNPV", true),
		Entry("Match word", "Mamestra configurata NPV-A", true),
		Entry("Match word", "Bacteriophage PH75", true),
	)

	DescribeTable("NoParse",
		func(s string, itIs bool) {
			res := NoParse([]byte(s))
			Expect(res).To(Equal(itIs))
		},
		Entry("No match", "Homo sapiens", false),
		Entry("No word at the start", "Not Homo sapiens", true),
		Entry("Noword at the start", "Nothomo sapiens", false),
		Entry("Not word at the start", "Not Homo sapiens", true),
		Entry("None word at the start", "None Homo sapiens", true),
		Entry("Unidentified at the start", "Unidentified species", true),
		Entry("Incertae sedis1", "incertae sedis", true),
		Entry("Incertae sedis2", "Incertae Sedis", true),
		Entry("Incertae sedis3", "Something incertae sedis", true),
		Entry("Incertae sedis4", "Homo sapiens inc.sed.", true),
		Entry("Incertae sedis5", "Incertae sedis", true),
		Entry("Phytoplasma in the middle", "Homo sapiensphytoplasmaoid", false),
		Entry("Phytoplasma in the end", "Homo sapiensphytoplasma Linn", true),
		Entry("Phytoplasma in the end", "Homo sapiensphytoplasma Linn", true),
		Entry("Plasmid1", "E. coli plasmids", true),
		Entry("Plasmid2", "E. coli plasmidia", false),
		Entry("Plasmid3", "E. coli plasmid", true),
		Entry("RNA1", "E. coli RNA", true),
		Entry("RNA2", "E. coli 32RNA", true),
		Entry("RNA3", "KURNAKOV", false),
		Entry("RNA4", "E. coli mRNA", true),
	)

	DescribeTable("Annotations",
		func(s string, body string, tail string) {
			bs := []byte(s)
			i := Annotation(bs)
			Expect(string(bs[0:i])).To(Equal(body))
			Expect(string(bs[i:])).To(Equal(tail))
		},
		Entry("No tail", "Homo sapiens", "Homo sapiens", ""),
		Entry("No tail", "Homo sapiens S. S.", "Homo sapiens S. S.", ""),
		Entry("No tail", "Homo sapiens s. s.", "Homo sapiens", " s. s."),
		Entry("No tail", "Homo sapiens sensu Linn.", "Homo sapiens", " sensu Linn."),
		Entry("No tail", "Homo sapiens nomen nudum", "Homo sapiens", " nomen nudum"),
	)

	Describe("Preprocess", func() {
		It("does not remove spaces at the start of a string", func() {
			name := "    Asplenium       × inexpectatum(E. L. Braun ex Friesner      )Morton"
			res := Preprocess([]byte(name))
			Expect(string(res.Body)).To(Equal(name))
		})
	})
})