package sqlite

import "strings"

// searchFolder maps letters with diacritics to their plain base letters, so a search
// typed on the POS letter keyboard (A–Z only) finds "Müller" with "muller" and
// "Straße" with "strasse". The input is lower-cased first.
var searchFolder = strings.NewReplacer(
	"ä", "a", "á", "a", "à", "a", "â", "a", "ã", "a", "å", "a",
	"æ", "ae",
	"ç", "c", "č", "c", "ć", "c",
	"é", "e", "è", "e", "ê", "e", "ë", "e", "ě", "e",
	"í", "i", "ì", "i", "î", "i", "ï", "i",
	"ñ", "n", "ń", "n", "ň", "n",
	"ö", "o", "ó", "o", "ò", "o", "ô", "o", "õ", "o", "ø", "o",
	"œ", "oe",
	"ř", "r",
	"ß", "ss", "š", "s", "ś", "s", "ş", "s",
	"ü", "u", "ú", "u", "ù", "u", "û", "u", "ů", "u",
	"ý", "y", "ÿ", "y",
	"ž", "z", "ź", "z", "ż", "z",
	"ł", "l",
)

// FoldForSearch normalises a text for a search that ignores case and diacritics.
func FoldForSearch(text string) string {
	return searchFolder.Replace(strings.ToLower(text))
}

// guestMatchesQuery reports whether a guest matches a POS search: the name contains
// the query (ignoring case and diacritics), or the code starts with it (ignoring case).
func guestMatchesQuery(name string, code *string, query string) bool {
	query = strings.TrimSpace(query)
	if query == "" {
		return true
	}

	if strings.Contains(FoldForSearch(name), FoldForSearch(query)) {
		return true
	}

	return code != nil && strings.HasPrefix(strings.ToLower(*code), strings.ToLower(query))
}
