package helpers

import (
	"strconv"
	"testing"
)

type addIsValidReferrerTest struct {
	ref, websiteDomain string
	result             bool
}

var isValidReferrerTest = []addIsValidReferrerTest{
	{"", "oldavista.com", false},
	{"null", "oldavista.com", false},
	{"#blah=1&blah2=2", "oldavista.com", false},
	{"/search.php", "oldavista.com", false},
	{"http://www.oldavista.com", "oldavista.com", false},

	{"http://www.ericexperiment.com", "oldavista.com", true},
	{"https://www.ericexperiment.com", "oldavista.com", true},
	{"ftp://www.ericexperiment.com", "oldavista.com", true},

	{"http://www.ericexperiment.com/windows31", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31/", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31?", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31/downloads.php", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31/downloads.php/", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31/downloads.php?", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31/downloads.php?part=1", "oldavista.com", true},
	{"http://www.ericexperiment.com/windows31/downloads.php?part=1&other=2", "oldavista.com", true},

	{"    http://www.ericexperiment.com/windows31    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31/    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31?    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31/downloads.php    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31/downloads.php/    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31/downloads.php?    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31/downloads.php?part=1    ", "oldavista.com", true},
	{"    http://www.ericexperiment.com/windows31/downloads.php?part=1&other=2    ", "oldavista.com", true},

	{"https://www.ericexperiment.com/windows31", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31/", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31?", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31/downloads.php", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31/downloads.php/", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31/downloads.php?", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31/downloads.php?part=1", "oldavista.com", true},
	{"https://www.ericexperiment.com/windows31/downloads.php?part=1&other=2", "oldavista.com", true},

	{"    https://www.ericexperiment.com/windows31    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31/    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31?    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31/downloads.php    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31/downloads.php/    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31/downloads.php?    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31/downloads.php?part=1    ", "oldavista.com", true},
	{"    https://www.ericexperiment.com/windows31/downloads.php?part=1&other=2    ", "oldavista.com", true},
}

func TestIsValidReferrer(t *testing.T) {
	for _, test := range isValidReferrerTest {
		result := isValidReferrer(test.ref, test.websiteDomain)
		if result != test.result {
			t.Errorf("For %s Domain was incorrect, got: %s, want: %s.", test.ref, strconv.FormatBool(result), strconv.FormatBool(test.result))
		}
	}
}

type addFilterReferrerTest struct {
	ref, websiteDomain, domain, fullUrl string
}

var filterReferrerTests = []addFilterReferrerTest{
	{"", "oldavista.com", "(none)", ""},
	{"    ", "oldavista.com", "(none)", ""},
	{"null", "oldavista.com", "(none)", ""},
	{"#blah=1&blah2=2", "oldavista.com", "(none)", ""},
	{"/search.php", "oldavista.com", "(none)", ""},
	{"  null  ", "oldavista.com", "(none)", ""},
	{"   #blah=1&blah2=2  ", "oldavista.com", "(none)", ""},
	{"  /search.php  ", "oldavista.com", "(none)", ""},
	{"http://www.oldavista.com", "oldavista.com", "(none)", ""},
	{"https://www.OLDAVISTA.com", "oldavista.com", "(none)", ""},
	{"  http://www.oldavista.com  ", "oldavista.com", "(none)", ""},
	{"  https://www.OLDAVISTA.com  ", "oldavista.com", "(none)", ""},
	{"http://www.ericexperiment.com", "oldavista.com", "ericexperiment.com", ""},

	{"http://www.ericexperiment.com", "oldavista.com", "ericexperiment.com", ""},
	{"https://www.ericexperiment.com", "oldavista.com", "ericexperiment.com", ""},
	{"ftp://www.ericexperiment.com", "oldavista.com", "(none)", ""},

	{"http://www.ericexperiment.com/windows31", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"http://www.ericexperiment.com/windows31/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"http://www.ericexperiment.com/windows31?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"http://www.ericexperiment.com/windows31/downloads.php", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"http://www.ericexperiment.com/windows31/downloads.php/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"http://www.ericexperiment.com/windows31/downloads.php?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"http://www.ericexperiment.com/windows31/downloads.php?part=1", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"http://www.ericexperiment.com/windows31/downloads.php?part=1&other=2", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"    http://www.ericexperiment.com/windows31    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    http://www.ericexperiment.com/windows31/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    http://www.ericexperiment.com/windows31?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    http://www.ericexperiment.com/windows31/downloads.php    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    http://www.ericexperiment.com/windows31/downloads.php/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    http://www.ericexperiment.com/windows31/downloads.php?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    http://www.ericexperiment.com/windows31/downloads.php?part=1    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"    http://www.ericexperiment.com/windows31/downloads.php?part=1&other=2    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"https://www.ericexperiment.com/windows31", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"https://www.ericexperiment.com/windows31/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"https://www.ericexperiment.com/windows31?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"https://www.ericexperiment.com/windows31/downloads.php", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"https://www.ericexperiment.com/windows31/downloads.php/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"https://www.ericexperiment.com/windows31/downloads.php?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"https://www.ericexperiment.com/windows31/downloads.php?part=1", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"https://www.ericexperiment.com/windows31/downloads.php?part=1&other=2", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"    https://www.ericexperiment.com/windows31    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    https://www.ericexperiment.com/windows31/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    https://www.ericexperiment.com/windows31?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    https://www.ericexperiment.com/windows31/downloads.php    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    https://www.ericexperiment.com/windows31/downloads.php/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    https://www.ericexperiment.com/windows31/downloads.php?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    https://www.ericexperiment.com/windows31/downloads.php?part=1    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"    https://www.ericexperiment.com/windows31/downloads.php?part=1&other=2    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"http://www.ERICEXPERIMENT.com", "oldavista.com", "ericexperiment.com", ""},
	{"https://www.ERICEXPERIMENT.com", "oldavista.com", "ericexperiment.com", ""},
	{"ftp://www.ERICEXPERIMENT.com", "oldavista.com", "(none)", ""},

	{"http://www.ERICEXPERIMENT.com/windows31", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"http://www.ERICEXPERIMENT.com/windows31/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"http://www.ERICEXPERIMENT.com/windows31?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"http://www.ERICEXPERIMENT.com/windows31/downloads.php", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"http://www.ERICEXPERIMENT.com/windows31/downloads.php/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"http://www.ERICEXPERIMENT.com/windows31/downloads.php?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"http://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"http://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1&other=2", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"    http://www.ERICEXPERIMENT.com/windows31    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    http://www.ERICEXPERIMENT.com/windows31/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    http://www.ERICEXPERIMENT.com/windows31?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    http://www.ERICEXPERIMENT.com/windows31/downloads.php    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    http://www.ERICEXPERIMENT.com/windows31/downloads.php/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    http://www.ERICEXPERIMENT.com/windows31/downloads.php?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    http://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"    http://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1&other=2    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"https://www.ERICEXPERIMENT.com/windows31", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"https://www.ERICEXPERIMENT.com/windows31/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"https://www.ERICEXPERIMENT.com/windows31?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"https://www.ERICEXPERIMENT.com/windows31/downloads.php", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"https://www.ERICEXPERIMENT.com/windows31/downloads.php/", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"https://www.ERICEXPERIMENT.com/windows31/downloads.php?", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"https://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"https://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1&other=2", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},

	{"    https://www.ERICEXPERIMENT.com/windows31    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    https://www.ERICEXPERIMENT.com/windows31/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    https://www.ERICEXPERIMENT.com/windows31?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31"},
	{"    https://www.ERICEXPERIMENT.com/windows31/downloads.php    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    https://www.ERICEXPERIMENT.com/windows31/downloads.php/    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    https://www.ERICEXPERIMENT.com/windows31/downloads.php?    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php"},
	{"    https://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1"},
	{"    https://www.ERICEXPERIMENT.com/windows31/downloads.php?part=1&other=2    ", "oldavista.com", "ericexperiment.com", "ericexperiment.com/windows31/downloads.php?part=1&other=2"},
}

func TestFilterReferrer(t *testing.T) {
	for _, test := range filterReferrerTests {
		result1, result2 := FilterReferrer(test.ref, test.websiteDomain)
		if result1 != test.domain {
			t.Errorf("For %s Domain was incorrect, got: %s, want: %s.", test.ref, result1, test.domain)
		}
		if result2 != test.fullUrl {
			t.Errorf("For %s FullUrl was incorrect, got: %s, want: %s.", test.ref, result2, test.fullUrl)
		}
	}
}
