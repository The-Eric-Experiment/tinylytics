package helpers

import "testing"

type addTest struct {
	input, expected string
}

type addTestDoubleExpected struct {
	input, expected1, expected2 string
}

var removeWWWTests = []addTest{
	{"www.oldavista.com", "oldavista.com"},
	{"oldavista.com", "oldavista.com"},
	{"dash.oldavista.com", "dash.oldavista.com"},
	{"http://www.ericexperiment.com", "http://www.ericexperiment.com"},
}

func TestRemoveWWW(t *testing.T) {

	for _, test := range removeWWWTests {
		result := RemoveWWW(test.input)
		if result != test.expected {
			t.Errorf("Result was incorrect, got: %s, want: %s.", result, test.expected)
		}
	}
}

var removeTrailingSlashTests = []addTest{
	{"/search.php?hello=1/", "/search.php?hello=1"},
	{"search.php?hello=1/", "search.php?hello=1"},
	{"www.google.com/", "www.google.com"},
	{"www.google.com/////", "www.google.com"},
	{"www.google.com", "www.google.com"},
}

func TestRemoveTrailingSlash(t *testing.T) {
	for _, test := range removeTrailingSlashTests {
		result := RemoveTrailingSlash(test.input)
		if result != test.expected {
			t.Errorf("Result was incorrect, got: %s, want: %s.", result, test.expected)
		}
	}
}

var removePrecedingSlashTests = []addTest{
	{"/search.php?hello=1/", "search.php?hello=1/"},
	{"/search.php?hello=1", "search.php?hello=1"},
	{"/www.google.com", "www.google.com"},
	{"//////www.google.com", "www.google.com"},
	{"www.google.com", "www.google.com"},
}

func TestRemovePrecedingSlash(t *testing.T) {
	for _, test := range removePrecedingSlashTests {
		result := RemovePrecedingSlash(test.input)
		if result != test.expected {
			t.Errorf("Result was incorrect, got: %s, want: %s.", result, test.expected)
		}
	}
}

//////

var withQueryPrefixTests = []addTest{
	{"", ""},
	{"?hello=1", "?hello=1"},
	{"hello=1", "?hello=1"},
}

func TestWithQueryPrefix(t *testing.T) {
	for _, test := range withQueryPrefixTests {
		result := WithQueryPrefix(test.input)
		if result != test.expected {
			t.Errorf("Result was incorrect, got: %s, want: %s.", result, test.expected)
		}
	}
}

var cleanupUrlTests = []addTestDoubleExpected{
	{"blaurus", "(none)", ""},
	{"search.php?hello=1", "(none)", ""},
	{"www.oldavista.com", "(none)", ""},
	{"ftp://www.oldavista.com", "(none)", ""},
	{"http://www.oldavista.com", "oldavista.com", ""},
	{"https://www.oldavista.com", "oldavista.com", ""},

	{"http://www.oldavista.com/search.php", "oldavista.com", "oldavista.com/search.php"},
	{"https://www.oldavista.com/search.php", "oldavista.com", "oldavista.com/search.php"},
	{"http://www.oldavista.com/sub/search.php", "oldavista.com", "oldavista.com/sub/search.php"},
	{"https://www.oldavista.com/sub/search.php", "oldavista.com", "oldavista.com/sub/search.php"},

	{"http://www.oldavista.com/search.php?s=Potato&search=Search", "oldavista.com", "oldavista.com/search.php?s=Potato&search=Search"},
	{"https://www.oldavista.com/search.php?s=Potato&search=Search", "oldavista.com", "oldavista.com/search.php?s=Potato&search=Search"},

	{"http://www.oldavista.com/sub/search.php?s=Potato&search=Search", "oldavista.com", "oldavista.com/sub/search.php?s=Potato&search=Search"},
	{"https://www.oldavista.com/sub/search.php?s=Potato&search=Search", "oldavista.com", "oldavista.com/sub/search.php?s=Potato&search=Search"},

	{"http://www.oldavista.com/search.php?", "oldavista.com", "oldavista.com/search.php"},
	{"https://www.oldavista.com/search.php?", "oldavista.com", "oldavista.com/search.php"},

	{"http://www.oldavista.com/sub/search.php?", "oldavista.com", "oldavista.com/sub/search.php"},
	{"https://www.oldavista.com/sub/search.php?", "oldavista.com", "oldavista.com/sub/search.php"},
}

func TestCleanupUrl(t *testing.T) {
	for _, test := range cleanupUrlTests {
		result1, result2 := CleanupUrl(test.input)
		if result1 != test.expected1 {
			t.Errorf("Result 1 was incorrect, got: %s, want: %s.", result1, test.expected1)
		}
		if result2 != test.expected2 {
			t.Errorf("Result 2 was incorrect, got: %s, want: %s.", result2, test.expected2)
		}
	}
}
