package dt_test

import (
	. "github.com/Contra-Culture/dt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("dumb templates", func() {
	Describe("styling", func() {
		It("generates CSS stylesheets", func() {
			s := S(
				"main",
				R(
					B(
						D("margin", "0"),
						"padding: 0"),
					"html", "body"))
			Expect(s).NotTo(BeNil())
			Expect(s.Compile()).To(Equal("html,\nbody {\n\tmargin: 0;\npadding: 0\n}\n\n"))

			s.Append(R("{\n\tfont-family: Helvetica, Arial, sans-serif\n}", "body"))
			Expect(s.Compile()).To(Equal("html,\nbody {\n\tmargin: 0;\npadding: 0\n}\n\nbody {\n\tfont-family: Helvetica, Arial, sans-serif\n}\n\n"))

			s.Append(
				R(
					B(
						D("font-family", "Times New Roman", "serif")),
					"p"))
			Expect(s.Compile()).To(Equal("html,\nbody {\n\tmargin: 0;\npadding: 0\n}\n\nbody {\n\tfont-family: Helvetica, Arial, sans-serif\n}\n\np {\n\tfont-family: Times New Roman, serif;\n}\n\n"))

			styling := s.S("cardHeader")
			styling.RT(
				B(
					D("border-bottom", "1px solid #e5e6e7"),
					D("padding", "1rem"),
					D("margin-bottom", "1rem")),
				T(I("self")))
			styling.RT(
				B(
					D("font-size", "1.618rem"),
					D("font-weight", "600"),
					D("font-family", "Calibri", "Heletica Neue", "Helvetica", "Arial", "sans-serif"),
					D("color", "#223344")),
				T(I("self"), " > h2"))
			s.C("the end of predefined styles")

			Expect(func() { s.S("cardHeader") }).To(Panic())
			Expect(s.SC("Publication-Header", "cardHeader")).To(Equal("Publication-Header"))
			Expect(s.Compile()).To(Equal("html,\nbody {\n\tmargin: 0;\npadding: 0\n}\n\nbody {\n\tfont-family: Helvetica, Arial, sans-serif\n}\n\np {\n\tfont-family: Times New Roman, serif;\n}\n\n\n/* the end of predefined styles */\n/* styling: cardHeader */\n.Publication-Header {\n\tborder-bottom: 1px solid #e5e6e7;\n\tpadding: 1rem;\n\tmargin-bottom: 1rem;\n}\n\n.Publication-Header > h2 {\n\tfont-size: 1.618rem;\n\tfont-weight: 600;\n\tfont-family: Calibri, Heletica Neue, Helvetica, Arial, sans-serif;\n\tcolor: #223344;\n}\n\n"))
		})
	})
	Describe("templates", func() {
		It("creates templates and renders views", func() {
			safe := Safe("<html><head><title>title</title></head><body>text</body></html>")
			Expect(safe).To(Equal("&lt;html&gt;&lt;head&gt;&lt;title&gt;title&lt;/title&gt;&lt;/head&gt;&lt;body&gt;text&lt;/body&gt;&lt;/html&gt;"))

			t := T(Safe("<html><head><title>title</title></head><body>text</body></html>"))
			s, err := t.Render()
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("&lt;html&gt;&lt;head&gt;&lt;title&gt;title&lt;/title&gt;&lt;/head&gt;&lt;body&gt;text&lt;/body&gt;&lt;/html&gt;"))

			t = T("<html><head><title>title</title></head><body>text</body></html>")
			s, err = t.Render()
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>title</title></head><body>text</body></html>"))

			t = T("<html><head><title>", Safe("<meta/>"), "</title></head><body>text</body></html>")
			s, err = t.Render()
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>&lt;meta/&gt;</title></head><body>text</body></html>"))

			t = T("<html><head><title>", I("title"), "</title></head><body>text</body></html>")
			s, err = t.Render()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("*Template.Render(): injection [0] not provided, got: \"[]string(nil)\""))
			Expect(s).To(BeEmpty())

			t = T("<html><head><title>", I("title"), "</title></head><body>text</body></html>")
			s, err = t.Render(Safe("<test title>"))
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>&lt;test title&gt;</title></head><body>text</body></html>"))

			t = T("<html><head><title>", I("title"), "</title></head><body>text</body></html>")
			s, err = t.Render("<test title>")
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title><test title></title></head><body>text</body></html>"))

			t = T(
				"<html><head><title>",
				I("title"),
				"</title></head><body>",
				I("body"),
				"</body></html>")
			s, err = t.Render(Safe("<test title>"), "test <strong>body</strong>")
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>&lt;test title&gt;</title></head><body>test <strong>body</strong></body></html>"))
		})
	})
})
