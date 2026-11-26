package toml

import (
	"testing"
)

func TestARefusalNamesLineAndColumn(t *testing.T) {
	pe := refused(t, "a = 1\nb = 2\na = 3\n", "duplicate key")
	if pe.Line != 3 || pe.Column != 1 {
		t.Errorf("refused at line %d, column %d, want line 3, column 1", pe.Line, pe.Column)
	}
	pe = refused(t, "x = 1\ny = tru\n")
	if pe.Line != 2 || pe.Column != 5 {
		t.Errorf("refused at line %d, column %d, want line 2, column 5", pe.Line, pe.Column)
	}
}

func TestColumnsCountCharacters(t *testing.T) {
	pe := refused(t, "'é' = 1 x\n")
	if pe.Line != 1 || pe.Column != 9 {
		t.Errorf("refused at line %d, column %d, want line 1, column 9", pe.Line, pe.Column)
	}
}

func TestTheErrorsText(t *testing.T) {
	_, err := Decode([]byte("a = 1\na = 2\n"))
	if want := `toml: line 2, column 1: duplicate key "a"`; err == nil || err.Error() != want {
		t.Errorf("error %v, want %q", err, want)
	}
}

func TestARefusedDocumentGivesNothing(t *testing.T) {
	doc, err := Decode([]byte("a = 1\n[t]\nb = 'x'\nc = \n"))
	if err == nil || doc != nil {
		t.Errorf("Decode = %#v, %v; want a nil table and an error", doc, err)
	}
}
