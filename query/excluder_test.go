package query

import "testing"

func TestShouldExclude_ExpectAllExcluded(t *testing.T) {
	exclusions := []string{".git", ".gitignore"}
	excluder := RegexpExclude{exclusions: exclusions}

	b := excluder.ShouldExclude(".git")
	if b == false {
		t.Fail()
	}

	b = excluder.ShouldExclude(".git/")
	if b == false {
		t.Fail()
	}

	b = excluder.ShouldExclude(".git/some/other/file")
	if b == false {
		t.Fail()
	}

	b = excluder.ShouldExclude(".gitignore")
	if b == false {
		t.Fail()
	}
}

func TestShouldExclude_ExpectNotExcluded(t *testing.T) {
	exclusions := []string{".git"}
	excluder := RegexpExclude{exclusions: exclusions}

	b := excluder.ShouldExclude(".gitignore")

	if b == true {
		t.Fail()
	}
}
