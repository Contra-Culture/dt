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

			styling, err := s.S("cardHeader")
			Expect(err).NotTo(HaveOccurred())
			styling.RT(
				B(
					D("border-bottom", "1px solid #e5e6e7"),
					D("padding", "1rem"),
					D("margin-bottom", "1rem")),
				T(I(SELF)))
			styling.RT(
				B(
					D("font-size", "1.618rem"),
					D("font-weight", "600"),
					D("font-family", "Calibri", "Heletica Neue", "Helvetica", "Arial", "sans-serif"),
					D("color", "#223344")),
				T(I(SELF), " > h2"))

			_styling, err := s.S("cardHeader")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("styling \"cardHeader\" already specified"))
			Expect(_styling).To(BeNil())
			Expect(s.SC("Publication-Header", "cardHeader", nil)).To(Equal("Publication-Header"))
			Expect(s.Compile()).To(Equal("html,\nbody {\n\tmargin: 0;\npadding: 0\n}\n\nbody {\n\tfont-family: Helvetica, Arial, sans-serif\n}\n\np {\n\tfont-family: Times New Roman, serif;\n}\n\n/* styling: cardHeader */\n.Publication-Header {\n\tborder-bottom: 1px solid #e5e6e7;\n\tpadding: 1rem;\n\tmargin-bottom: 1rem;\n}\n\n.Publication-Header > h2 {\n\tfont-size: 1.618rem;\n\tfont-weight: 600;\n\tfont-family: Calibri, Heletica Neue, Helvetica, Arial, sans-serif;\n\tcolor: #223344;\n}\n\n"))
		})
	})
	Describe("templates", func() {
		It("creates templates and renders views", func() {
			protected := Protected("<html><head><title>title</title></head><body>text</body></html>")
			Expect(protected).To(Equal("&lt;html&gt;&lt;head&gt;&lt;title&gt;title&lt;/title&gt;&lt;/head&gt;&lt;body&gt;text&lt;/body&gt;&lt;/html&gt;"))

			t := T(Protected("<html><head><title>title</title></head><body>text</body></html>"))
			s, err := t.Render(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("&lt;html&gt;&lt;head&gt;&lt;title&gt;title&lt;/title&gt;&lt;/head&gt;&lt;body&gt;text&lt;/body&gt;&lt;/html&gt;"))

			t = T("<html><head><title>title</title></head><body>text</body></html>")
			s, err = t.Render(map[string]string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>title</title></head><body>text</body></html>"))

			t = T("<html><head><title>", Protected("<meta/>"), "</title></head><body>text</body></html>")
			s, err = t.Render(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>&lt;meta/&gt;</title></head><body>text</body></html>"))

			t = T("<html><head><title>", I("title"), "</title></head><body>text</body></html>")
			s, err = t.Render(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("injection \"title\" not provided"))
			Expect(s).To(BeEmpty())

			t = T("<html><head><title>", I("title"), "</title></head><body>text</body></html>")
			s, err = t.Render(map[string]string{"title": "<test title>"})
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>&lt;test title&gt;</title></head><body>text</body></html>"))

			t = T("<html><head><title>", UI("title"), "</title></head><body>text</body></html>")
			s, err = t.Render(map[string]string{"title": "<test title>"})
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title><test title></title></head><body>text</body></html>"))

			t = T(
				"<html><head><title>",
				I("title"),
				"</title></head><body>",
				UI("body"),
				"</body></html>")
			s, err = t.Render(map[string]string{
				"title": "<test title>",
				"body":  "test <strong>body</strong>"})
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(Equal("<html><head><title>&lt;test title&gt;</title></head><body>test <strong>body</strong></body></html>"))
		})
	})
})
