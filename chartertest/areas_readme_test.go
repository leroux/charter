package chartertest

import (
	"strings"
	"testing"

	"github.com/leroux/charter/claims"
)

// README.md beside the area documents is documentation: its headings and
// its example entries are not read.
func TestReadmeIsNotAnAreaDocument(t *testing.T) {
	Verifies(t, claims.CH_010, "readme-is-not-a-document")
	root := copyFixture(t)
	writeFile(t, root, "claims/README.md", "# Charter\n\n## SW-001 An example that would be declared twice\n\nStatus: active\n\nText.\n")
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("the README was read as a document:\n%s", strings.Join(problems, "\n"))
	}
}

// audit.md beside the area documents is the audit's report, read by
// people: its headings, which name entries, are not read as entries.
func TestTheAuditReportIsNotAnAreaDocument(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "claims/audit.md", "# Audit\n\n## SW-001 The decision the audit read\n\nIts scenario's test checks less than the sentence says.\n\n### Findings\n")
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("the audit's report was read as a document:\n%s", strings.Join(problems, "\n"))
	}
	ids, err := NextIDs(fixtureOptions(root), "SW", 1)
	if err != nil || len(ids) != 1 {
		t.Fatalf("next with an audit report: %v, %v", ids, err)
	}
}

// A spec file of the earlier form may link to an entry that has moved into
// an area document; the link's file name is then not checked.
func TestALinkIntoAnAreaDocumentIsCheckedByIdAlone(t *testing.T) {
	Verifies(t, claims.CH_014, "link-into-a-document")
	root := copyFixture(t)
	writeFile(t, root, "claims/refill.md", sampleArea)
	writeFile(t, root, "pkg/a/refill.go", refillGo)
	writeFile(t, root, "pkg/a/refill_test.go", refillTest)
	appendTo(t, root, "claims/specs/CON-001-lock-file-removed-on-close.md", "\n**Related**\n\n- [SW-010](SW-010-the-bucket-refills.md)\n")
	if problems := check(t, root); len(problems) != 0 {
		t.Fatalf("a link into a document reported:\n%s", strings.Join(problems, "\n"))
	}
	appendTo(t, root, "claims/specs/CON-001-lock-file-removed-on-close.md", "- [SW-099](SW-099-nowhere.md)\n")
	reports(t, root, "rule 23: CON-001 is related to SW-099, which is not a spec")
}
