package epub

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
)

const (
	navBodyTemplate = `
    <nav epub:type="toc">
      <h1>目次</h1>
      <ol>
				<span>サブタイトル</span>
      </ol>
    </nav>
`
	navFilename       = "navigation.xhtml"
	navItemID         = "nav"
	navItemProperties = "nav"
	navEpubType       = "toc"

	navXmlnsEpub = "http://www.idpf.org/2007/ops"
)

// nav implements the EPUB table of contents
type nav struct {
	// This holds the body XML for the EPUB v3 TOC file (nav.xhtml). Since this is
	// an XHTML file, the rest of the structure is handled by the xhtml type
	//
	// Sample: https://github.com/SaitoJP/epub-samples/blob/master/minimal-v3plus2/EPUB/nav.xhtml
	// Spec: http://www.idpf.org/epub/301/spec/epub-contentdocs.html#sec-xhtml-nav
	navXML *navBody

	title   string // EPUB title
	cssPath string
}

type navBody struct {
	XMLName  xml.Name  `xml:"nav"`
	EpubType string    `xml:"epub:type,attr"`
	H1       string    `xml:"h1"`
	Links    []navItem `xml:"ol>li"`
}

type navItem struct {
	A navLink `xml:"a"`
}

type navLink struct {
	XMLName xml.Name `xml:"a"`
	Class   string   `xml:"class,attr"`
	Href    string   `xml:"href,attr"`
	Data    string   `xml:",chardata"`
}

// Constructor for nav
func newNav() *nav {
	t := &nav{}

	t.navXML = newNavXML()

	return t
}

// Constructor for navBody
func newNavXML() *navBody {
	b := &navBody{
		EpubType: navEpubType,
	}
	err := xml.Unmarshal([]byte(navBodyTemplate), &b)
	if err != nil {
		panic(fmt.Sprintf(
			"Error unmarshalling navBody: %s\n"+
				"\tnavBody=%#v\n"+
				"\tnavBodyTemplate=%s",
			err,
			*b,
			navBodyTemplate))
	}

	return b
}

// Add a section to the Navigation (navXML)
func (t *nav) addSection(index int, title string, relativePath string, isNavigationPage bool) {
	var className string
	if isNavigationPage {
		className = "heading1"
	} else {
		className = "heading2"
	}
	relativePath = filepath.ToSlash(relativePath)
	l := &navItem{
		A: navLink{
			Class: className,
			Href:  relativePath,
			Data:  title,
		},
	}

	t.navXML.Links = append(t.navXML.Links, *l)
}

func (t *nav) setTitle(title string) {
	t.title = title
}

func (t *nav) setCSS(path string) {
	t.cssPath = path
}

// Write the Navigation files
func (t *nav) write(tempDir string) {
	t.writeNavDoc(tempDir)
}

// Write the the Navigation file (nav.xhtml) to the temporary directory
func (t *nav) writeNavDoc(tempDir string) {
	navBodyContent, err := xml.MarshalIndent(t.navXML, "    ", "  ")
	if err != nil {
		panic(fmt.Sprintf(
			"Error marshalling XML for Navigation file: %s\n"+
				"\tXML=%#v",
			err,
			t.navXML))
	}

	n := newXhtml(string(navBodyContent))
	n.setXmlnsEpub(navXmlnsEpub)
	n.setTitle(t.title)
	if t.cssPath != "" {
		n.setCSS(t.cssPath)
	}

	navFilePath := filepath.Join(tempDir, contentFolderName, xhtmlFolderName, navFilename)
	n.write(navFilePath)
}
