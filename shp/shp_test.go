// SHP - Simple Html Parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package shp

import (
	"fmt"
	"testing"
)

func Test_tagraw(t *testing.T) {
	tagRaw := HtmlTag{rawSrc: runes_from_str(`<img src_str="url" width=45  alt='space in text' alt2="Peter's book  two_space" >`)}
	tagWithAttributes := attributesDetect(tagRaw)
	str_same(tagWithAttributes.Type, "img", t)
	str_same(tagWithAttributes.Attributes["src_str"], "url", t)
	str_same(tagWithAttributes.Attributes["width"], "45", t)
	str_same(tagWithAttributes.Attributes["alt"], "space in text", t)
	str_same(tagWithAttributes.Attributes["alt2"], "Peter's book  two_space", t)
}

func Test_char_funcs(t *testing.T){
	str_same("", txt_char_first(""), t)
	str_same("", txt_char_last(""), t)
	str_same("a", txt_char_first("apple"), t)
	str_same("e", txt_char_last("apple"), t)
	str_same("e", txt_char_last("apple"), t)

	bool_same(true, txt_char_first_is("apple", "a"), t)
	bool_same(true, txt_char_last_is("apple", "e"), t)
	bool_same(false, txt_char_first_is("apple", "x"), t)
	bool_same(false, txt_char_last_is("apple", "x"), t)

	// empty text with pattern?
	bool_same(false, txt_char_first_is("", "x"), t)
	bool_same(false, txt_char_last_is("", "x"), t)

	bool_same(true, txt_char_first_and_last_is("abba", "a"), t)
	bool_same(false, txt_char_first_and_last_is("abba", "X"), t)
	bool_same(false, txt_char_first_and_last_is("a", "X"), t) // too short str

	str_same("pple", txt_char_remove_first("apple"), t)
	str_same("appl", txt_char_remove_last("apple"), t)
	str_same("ppl", txt_char_first_and_last_removed("apple"), t)
	// empty basic string:
	str_same("", txt_char_remove_first(""), t)
	str_same("", txt_char_remove_last(""), t)

	// too short string, what happens if I remove first and last?
	str_same("", txt_char_first_and_last_removed("a"), t)

	key, val := key_value_attrib_split("a=b")
	str_same(key, "a", t)
	str_same(val, "b", t)

	key2, val2 := key_value_attrib_split("a")
	str_same(key2, "a", t)
	str_same(val2, "", t)

	runes := runes_from_str("ab")
	rune_same(runes[0], 'a', t)
	rune_same(runes[1], 'b', t)
}

func Test_insert_children_simple(t *testing.T) {
	tagH1 := HtmlTag{Type: "html",   OpenCloseStatus: tagOpen}
	tagP1 := HtmlTag{Type: "p",   OpenCloseStatus: tagOpen}
	tagB1 := HtmlTag{Type: "img", OpenCloseStatus: tagSelfClose}
	tagB2 := HtmlTag{Type: "img", OpenCloseStatus: tagSelfClose}
	tagB3 := HtmlTag{Type: "img", OpenCloseStatus: tagSelfClose}
	tagP2 := HtmlTag{Type: "p",   OpenCloseStatus: tagClose}
	tagH2 := HtmlTag{Type: "html",   OpenCloseStatus: tagClose}
	tagsSimpleList := HtmlTags{ tagH1, tagP1, tagB1, tagB2, tagB3, tagP2, tagH2}
	tagsStructured := DomStructureBuilder(tagsSimpleList)

	tagHtml := tagsStructured[0]
	int_same_msg("simple testRoot", len(tagHtml.Children), 1, t)

	tagP := tagHtml.Children[0]
	int_same_msg("simple test tag P children", len(tagP.Children), 3, t)
}

func Test_domStructureBuilder(t *testing.T) {
	tagH1 := HtmlTag{Type: "html",   OpenCloseStatus: tagOpen}
		tagP1 := HtmlTag{Type: "p",   OpenCloseStatus: tagOpen}
			tagImg := HtmlTag{Type: "img", OpenCloseStatus: tagSelfClose}
			tagA1 := HtmlTag{Type: "a", OpenCloseStatus: tagOpen}
				tagUl1 := HtmlTag{Type: "ul", OpenCloseStatus: tagOpen}
					tagLi1 := HtmlTag{Type: "li", OpenCloseStatus: tagOpen}
					tagLi2 := HtmlTag{Type: "li", OpenCloseStatus: tagClose}
					tagLi3 := HtmlTag{Type: "li", OpenCloseStatus: tagOpen}
					tagLi4 := HtmlTag{Type: "li", OpenCloseStatus: tagClose}
					tagLi5 := HtmlTag{Type: "li", OpenCloseStatus: tagOpen}
					tagLi6 := HtmlTag{Type: "li", OpenCloseStatus: tagClose}
					tagLi7 := HtmlTag{Type: "li", OpenCloseStatus: tagOpen}
					tagLi8 := HtmlTag{Type: "li", OpenCloseStatus: tagClose}
				tagUl2 := HtmlTag{Type: "ul", OpenCloseStatus: tagClose}
			tagA2 := HtmlTag{Type: "a", OpenCloseStatus: tagClose}
		tagP2 := HtmlTag{Type: "p",   OpenCloseStatus: tagClose}
	tagH2 := HtmlTag{Type: "html",   OpenCloseStatus: tagClose}
	tagsSimpleList := HtmlTags{
		tagH1,
		tagP1,
		tagImg,
		tagA1,
		tagUl1,
		tagLi1,
		tagLi2,
		tagLi3,
		tagLi4,
		tagLi5,
		tagLi6,
		tagLi7,
		tagLi8,
		tagUl2,
		tagA2,
		tagP2,
		tagH2,
	}
 	tagsStructured := DomStructureBuilder(tagsSimpleList)
	tagHtml := tagsStructured[0]
	fmt.Println("test, tagHtml type:", tagHtml.Type)
	int_same_msg("test tagHtml children", len(tagHtml.Children), 1, t)

	tagP_children := tagHtml.Children[0].Children
	int_same_msg("test tagP children", len(tagP_children), 2, t)

	tagUl := tagP_children[1].Children[0]
	int_same_msg("testLi", len(tagUl.Children), 4, t)
}

func Test_open_close_tag_count(t *testing.T) {
	tagA := HtmlTag{OpenCloseStatus: tagClose}
	tagB := HtmlTag{OpenCloseStatus: tagClose}
	tagC := HtmlTag{OpenCloseStatus: tagOpen}
	tags := HtmlTags{tagA, tagB, tagC}
	count_open, count_close := countOpenCloseTags(tags)
	int_same(count_open, 1, t)
	int_same(count_close, 2, t)
}

func Test_indentation(t *testing.T) {
	indent := indentation(5)
	str_same(indent, "     ", t)
}

func Test_basic_html(t *testing.T) {
	domRoot := HtmlParseSrc(`<html><body><p>basic <br />test</p></body></html>`)
	fmt.Println("TEST Html Parse Src")
	tagP := domRoot.Children[0].Children[0]
	fmt.Println("tagp", tagP.Type)
	int_same_msg("test basic, tag p:", len(tagP.Children), 3, t)
	str_same_msg("test basic, tag p child 3:", tagP.Children[2].src_str(), "test", t)
}

// TODO: build up this section with more complex checks (attributes with numbers without quotes...
func Test_complex_html(t *testing.T) {
	domRoot := HtmlParseSrc(`<html>
                                 <body>
                                     <h1 id="headerId">Text</h1>
                                     <div id="divId">
                                         <p>This is <br />a <b>paragraph</b> in paragraph</p>
                                     </div>
                                 </body>
                             </html>`)

	tagP := domRoot.Children[0].Children[1].Children[0]
	int_same_msg("test complex, tag p:", len(tagP.Children), 5, t)
	str_same_msg("test complex, inner text:", tagP.Children[0].src_str(), "This is ", t)
	str_same_msg("test complex, inner text:", tagP.Children[4].src_str(), " in paragraph", t)
}

func rune_same(received, wanted rune, t *testing.T) {
	if received != wanted {
		t.Fatalf("\rune received: %v\n  wanted: %v, error", received, wanted)
	}
}

func str_same(received, wanted string, t *testing.T) {
	str_same_msg("str same", received, wanted, t)
}

func str_same_msg(msg, received, wanted string, t *testing.T) {
	if received != wanted {
		t.Fatalf("\n%s string received: '%v'\n  wanted: '%v', error", msg, received, wanted)
	}
}

func bool_same(received, wanted bool, t *testing.T) {
	if received != wanted {
		t.Fatalf("\nbool received: %v\n  wanted: %v, error", received, wanted)
	}
}

func int_same(received, wanted int, t *testing.T) {
	int_same_msg("int same", received, wanted, t)
}
func int_same_msg(msg string, received, wanted int, t *testing.T) {
	if received != wanted {
		t.Fatalf("\n%s  int received: %v\n  wanted: %v, error", msg, received, wanted)
	}
}
