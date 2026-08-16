package dto

import "testing"

func TestListQueryNormalizeDefaults(t *testing.T) {
	q := &ListQuery{Page: 0, PageSize: 0}
	q.Normalize()
	if q.Page != 1 {
		t.Fatalf("Page = %d, want 1", q.Page)
	}
	if q.PageSize != 20 {
		t.Fatalf("PageSize = %d, want 20", q.PageSize)
	}
}

