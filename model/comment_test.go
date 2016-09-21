package model

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	section := "Args:\n" +
		" a descA\n" +
		" b descB\n" +
		" c descC\n"
	expected := map[string]string{
		"a": "descA",
		"b": "descB",
		"c": "descC",
	}

	args := make(map[string]string)
	parseArgs(section, args)

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("args %q not expected value (%q)", args, expected)
	}

}

func TestParseComment(t *testing.T) {
	commentText := "" +
		"Sets the test priority.\n" +
		"\n" +
		"Args:\n" +
		" t         Test to operate on.\n" +
		" priority  Priority of the test.\n" +
		"\n" +
		"Sets the priority of <t> to <priority>.\n"

	expectedTitle := "Sets the test priority."
	expectedBody := "Sets the priority of <t> to <priority>.\n"
	expectedArgs := map[string]string{
		"t":        "Test to operate on.",
		"priority": "Priority of the test.",
	}

	comment := ParseComment(commentText)

	if comment.Title != expectedTitle {
		t.Error("Expected", expectedTitle, "got", comment.Title)
	}
	if comment.Body != expectedBody {
		t.Error("Expected", expectedBody, "got", comment.Body)
	}
	if !reflect.DeepEqual(comment.Args, expectedArgs) {
		t.Error("Expected", expectedArgs, "got", comment.Args)
	}
}
