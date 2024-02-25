// SHP - Simple Html Parser
// author: Balazs Nyiro, balazs.nyiro.ca@gmail.com
package shp

import (
	"fmt"
	"strings"
)

var QUOTE_D = `"`
var QUOTE_S = "'"

var typeNotDetected = "typeNotDetected"
var tagSelfClose = "selfClose"
var tagOpen = "open"
var tagClose = "close"

type HtmlTags []HtmlTag

type HtmlTag struct {
	Type string //  "div", "a", "p", "ul"...
	//               innerText type is a non-official tag type, I use it because
	//               it helps the parsing and object structuring
	//
	//               "typeNotDetected" - another un-official type
	//
	//              "innerText": in <p>txt <b>important</b> word</p>
	//               in this case "txt " and " word" are two innerText elem
	//               or <li>something</li> or <h1>Head</h1>
	rawSrc                 []rune // the raw source code from html
	Attributes             map[string]string
	Children               HtmlTags
	OpenCloseStatus        string // open, close, selfClose
	TagPairedWithCloseElem bool
	TagCounter             int

	ExternalProgramMapStr  map[string]string // if this tool is used as a lib,
	ExternalProgramMapInt  map[string]int    // an external program can place his data here.
}

func (tag HtmlTag) src_str() string {
	return string(tag.rawSrc)
}

func (tag HtmlTag) IsOpening() bool {
	return tag.OpenCloseStatus == tagOpen
}
func (tag HtmlTag) IsClosing() bool {
	return tag.OpenCloseStatus == tagClose
}

func (tag HtmlTag) HasChild() bool {
	return len(tag.Children) > 0
}

func (tag HtmlTag) LenSourceCode() int {
	return len(tag.rawSrc)
}

func (tag HtmlTag) hasSourceCode() bool {
	// orig: return len(tag.rawSrc) > 0
	// the parser detects the \n \t ' ' whitespaces, and it creates
	// tags for them. if a raw source code has nothing but whitespaces,
	// this is not a real source code.
	src := strings.Replace(string(tag.rawSrc), "\n", "", -1)
	src = strings.Replace(src, " ", "", -1)
	src = strings.Replace(src, "\t", "", -1)
	return len(src) > 0 // has not-whitespace char
}

// be careful: the fun creates a copy about
// the original tag and returns with a new one
// this is why I return the actual tag and this is the usage;
// tagRaw = tagRaw.addSourceCode(runeNow)
func (tag HtmlTag) addSourceCode(runeNext rune) HtmlTag {
	tag.rawSrc = append(tag.rawSrc, runeNext)
	return tag
}

func (tag HtmlTag) print(indent_level int) {
	prefix := indentation(indent_level)
	fmt.Println("tag print, child:", len(tag.Children), "counterId:", tag.TagCounter, prefix, string(tag.rawSrc))
}

func (tag HtmlTag) printWithChildren(indent_level int) {
	tag.print(indent_level)
	for _, child := range tag.Children {
		child.printWithChildren(indent_level + 1)
	}
}

//////////////////////////////////////////////////
func unknownNewTag(tagCounter int) (HtmlTag, int) {
	return HtmlTag{Type: typeNotDetected, TagCounter: tagCounter}, tagCounter + 1
}

// TESTED
func attributesDetect(tagRaw HtmlTag) HtmlTag {
	if len(tagRaw.rawSrc) == 0 {
		return tagRaw
	} // > 0
	tagRaw.Attributes = make(map[string]string)

	// len minimum 1
	if tagRaw.rawSrc[0] == '<' {
		tagRaw.OpenCloseStatus = tagOpen
	}
	if len(tagRaw.rawSrc) > 1 {
		if tagRaw.rawSrc[1] == '/' {
			tagRaw.OpenCloseStatus = tagClose
		}
		// if prev before last is /, it is a self/close
		if tagRaw.rawSrc[len(tagRaw.rawSrc)-2] == '/' {
			tagRaw.OpenCloseStatus = tagSelfClose
		}
	}

	src := string(tagRaw.rawSrc)
	if src[0] == '<' {
		src = src[1:]
	}
	if txt_char_last(src) == ">" {
		src = txt_char_remove_last(src)
	}

	words := strings.Split(src, " ")
	if len(words) > 0 {
		id := -1
		for true {
			id++; if id >= len(words) { break }
			word := words[id]
			word = strings.TrimSpace(word)
			if id == 0 { // the first elem is not key/value pair, this is the tag type
				//          if the src starts with <: <a href="">
				tagRaw.Type = "inner_text"  // <p>Text in paragraph, inner text is detected as a type of tag</p>
				if tagRaw.rawSrc[0] == '<' { //
					tagRaw.Type = strings.ToLower(word)
					if word[0] == '/' { // <p>txt</p> the close /p type is p too
						tagRaw.Type = word[1:]
					}
				}
				continue
			}
			// then find the maybe space container elems between 'xxx' or "xxx"
			// and this is not correct because in quotes we can find spaces
			if strings.ContainsRune(word, '=') {
				key, value := key_value_attrib_split(word)
				quote_pair_double := txt_char_first_and_last_is(value, QUOTE_D)
				quote_pair_single := txt_char_first_and_last_is(value, QUOTE_S)
				complete_quote_pair_around_val := quote_pair_double || quote_pair_single

				uncomplete_quote_first :=
					(txt_char_first_is(value, QUOTE_D) || txt_char_first_is(value, QUOTE_S)) &&
					!(txt_char_last_is(value, QUOTE_D) || txt_char_last_is(value, QUOTE_S) )

				uncomplete_quote_last :=
					!(txt_char_first_is(value, QUOTE_D) || txt_char_first_is(value, QUOTE_S)) &&
					(txt_char_last_is(value, QUOTE_D) || txt_char_last_is(value, QUOTE_S) )

				if !(uncomplete_quote_first || uncomplete_quote_last) {
					if complete_quote_pair_around_val {
						value = txt_char_first_and_last_removed(value)
					}
					tagRaw.Attributes[key] = value
					continue
				}

				if uncomplete_quote_first {
					quote_used := txt_char_first(value)
					value = txt_char_remove_first(value)
					tagRaw.Attributes[key] = value
					// there is an opening quote (removed) and we have to find the pair of that

					closing_quote_detected := false
					for true { // read next words until we reach the closing quote
						id++; if id >= len(words) { break }
						word = words[id]
						if txt_char_last_is(word, quote_used) || txt_char_last_is(word, quote_used) {
							word = txt_char_remove_last(word)
							closing_quote_detected = true
						}
						tagRaw.Attributes[key] = tagRaw.Attributes[key] + " " + word
						if closing_quote_detected {
							break
						}
					}
				}
			} else { // no = sign in word
				tagRaw.Attributes[word] = ""
			}
		}
	}
	return tagRaw
}

func rawTagAppendWithAttributes(tagRaw HtmlTag, tagsUnsctructured HtmlTags) HtmlTags {
	if tagRaw.hasSourceCode() {
		tagRaw = attributesDetect(tagRaw)
		tagsUnsctructured = append(tagsUnsctructured, tagRaw) // save previous tag
	}
	return tagsUnsctructured
}

func HtmlParseSrc(src string) HtmlTag {
	//// detect unstructured elems of HTML ////
	tagsUnstructured := HtmlTags{}
	// inQuote := false // in html we don't use escaped quotes

	// if there is anything (empty line, space) before the root object,
	tagCounter := 0
	tagRaw, tagCounter := unknownNewTag(tagCounter)

	whatHappened := "what happened in the for loop last time"
	// step 1: detect objects step by step
	for _, runeNow := range src {
		// charNow := string(runeNow)
		// _ = charNow // to debug

		// the inside texts in <p> for example stored in different tag, too
		if runeNow == '<' {
			tagsUnstructured = rawTagAppendWithAttributes(tagRaw, tagsUnstructured)
			if whatHappened != "tagClose" {
				// at tagClose we start a new tag object, and it is totally empty
				// and in this situation I don't want to start an empty tag
				// between A and B (after a closing we start a new,
				// and at a new start I start a new one too): <A><B>
				tagRaw, tagCounter = unknownNewTag(tagCounter)
			}
			tagRaw = tagRaw.addSourceCode(runeNow)
			whatHappened = "tagOpen"
			continue
		}
		if runeNow == '>' {
			tagRaw = tagRaw.addSourceCode(runeNow)
			tagsUnstructured = rawTagAppendWithAttributes(tagRaw, tagsUnstructured)
			tagRaw, tagCounter = unknownNewTag(tagCounter)
			whatHappened = "tagClose"
			continue
		}
		tagRaw = tagRaw.addSourceCode(runeNow)
		whatHappened = "read a character from source code"
	}
	// fmt.Println(tagsUnstructured)

	// step 2: embed structures into each other:
	tagsStructured := DomStructureBuilder(tagsUnstructured)
	root := tagsStructured[0]
	// root.printWithChildren(0)
	return root
}

// TESTED
func DomStructureBuilder(tagsStructured HtmlTags) HtmlTags {
	for true {
		numOfOpenTags, numOfCloseTags := countOpenCloseTags(tagsStructured)
		// fmt.Println("Dom Struct Builder", numOfOpenTags, numOfCloseTags)
		if numOfOpenTags == 0 || numOfCloseTags == 0 {
			break
		}
		tagsStructured = tagPairingOneOpeningAndClosing__InsertChildren__onePairingStepUntilAllProcessed(tagsStructured)
	}
	return tagsStructured
}


// ####################### PAIRING function blocks. TESTED from DomStructureBuilder###################################
// tags have minimum 1 elem!
func _lastTag(tags HtmlTags) HtmlTag {
	return tags[len(tags)-1]
}

func _tagsAppendAll(acc, tags HtmlTags) HtmlTags {
	for _, tag := range tags{
		acc = append(acc, tag)
	}
	return acc
}

func _appendOpeningPrevious_ifHeHavePreviousOpeningTag(tagsPaired, openingTagsLast HtmlTags) HtmlTags{
	if len(openingTagsLast) > 0 { // save PREV OPENING TAG if we have
		openingTagPrevious := _lastTag(openingTagsLast)
		tagsPaired = append(tagsPaired, openingTagPrevious)
	}
	return tagsPaired
}

// there is minimum 1 opening tag or 1 closing tag when it is called.
// IMPORTANT: this func does only 1 step in children inserting -
// it pairs only one opening/closing tags in one execution - you have to call it until the num of closing tags will be 0
// TESTED in DomStructureBuilder
func tagPairingOneOpeningAndClosing__InsertChildren__onePairingStepUntilAllProcessed(tagsToPair HtmlTags) HtmlTags {
	tagsPaired := HtmlTags{}
	openingsNowDetected := HtmlTags{}

	var openingTagActual HtmlTag
	children := HtmlTags{}

	idNow := 0 // this var is available in the next for loop too
	for id, tag := range tagsToPair {
		idNow = id
		// first we always find an opening tags.
		if tag.IsOpening() && !tag.TagPairedWithCloseElem { // PairedWith... this flag can be set from previous step executions

			// Now we care only about the most inner OPENING-CLOSING PAIRS.
			// SAVE PREVIOUS RESULTS without any re-organisation ========================================
			tagsPaired = _appendOpeningPrevious_ifHeHavePreviousOpeningTag(tagsPaired, openingsNowDetected)
			childrenOfPrevOpening := children                              // we entered into a newer open tag
			tagsPaired = _tagsAppendAll(tagsPaired, childrenOfPrevOpening) // save the previous children into tagsPaired
			openingsNowDetected = append(openingsNowDetected, tag)
			// ==========================================================================================

			openingTagActual = tag  // and only focus the most inner opening tag.
			children = HtmlTags{}
			continue
		}

		// if the closing tag type is same than prev opening:
		// example: "<p>txt</p>" is ok, but "<p>incorrect html</a>" is not ok, p!=a types
		if tag.IsClosing() && tag.Type == openingTagActual.Type {
			openingTagActual.TagPairedWithCloseElem = true
			openingTagActual.Children = children
			tagsPaired = append(tagsPaired, openingTagActual)
			break // one pairing jog is done!
		}

		// the most important BASIC JOB: collect children
		children = append(children, tag)
	}

	for idToCopy, tagToCopy := range tagsToPair { // copy the elems after last Processed
		if idToCopy > idNow {
			tagsPaired = append(tagsPaired, tagToCopy)
		}
	}

	return tagsPaired
}
// ####################### PAIRING ##################################################################

// TESTED
func countOpenCloseTags(tags HtmlTags) (int, int) {
	numOfOpenTags, numOfCloseTags := 0, 0
	for _, tag := range tags {
		if tag.OpenCloseStatus == tagOpen {
			numOfOpenTags++
		}
		if tag.OpenCloseStatus == tagClose {
			numOfCloseTags++
		}
	}
	return numOfOpenTags, numOfCloseTags
}

// TESTED - duplication from jyp //////////
func indentation(level int) string {
	indent := ""
	for i := 0; i < level; i++ {
		indent = indent + " "
	}
	return indent
}

////////////// small tools ////////////
// TESTED
func runes_from_str(txt string) []rune {
	return []rune(txt)
}

// TESTED
func txt_char_first_and_last_is(txt, pattern string) bool {
	return txt_char_first_is(txt, pattern) && txt_char_last_is(txt, pattern)
}

// TESTED
func txt_char_last_is(txt, pattern string) bool {
	// empty can't be equal with anything
	if len(txt) < 1 { return false }
	return txt_char_last(txt) == txt_char_first(pattern)
}

// TESTED
func txt_char_first_is(txt, pattern string) bool {
	if len(txt) < 1 { return false }
	return txt_char_first(txt) == txt_char_first(pattern)
}

// TESTED
func txt_char_first(txt string) string {
	if len(txt) == 0 {
		return ""
	}
	return string(txt[0])
}

// TESTED
func txt_char_remove_first(txt string) string {
	if len(txt) == 0 {
		return ""
	}
	return txt[1 : len(txt)]
}

// I can't return with empty rune so I return with string.
// text can be empty!
// TESTED
func txt_char_last(txt string) string {
	if len(txt) == 0 {
		return ""
	}
	return string(txt[len(txt)-1])
}

// TESTED
func txt_char_remove_last(txt string) string {
	if len(txt) == 0 {
		return ""
	}
	return txt[0 : len(txt)-1]
}

// TESTED
func txt_char_first_and_last_removed(txt string) string {
	if len(txt) < 2 { return "" }
	return txt[1 : len(txt)-1]
}

// TESTED
func key_value_attrib_split(word string) (string, string) {
	key := word
	value := ""
	if strings.ContainsRune(word, '=')	{
		splitted := strings.Split(word, "=")
		key = splitted[0]
		value = splitted[1]
	}
	return key, value
}