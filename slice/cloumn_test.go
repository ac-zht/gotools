package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumn(t *testing.T) {
	testCase := []struct {
		name   string
		src    []map[string]any
		colKey string
		want   []any
	}{
		{
			name: "continuous",
			src: []map[string]any{
				{
					"id":         1,
					"first_name": "John",
				},
				{
					"id":         2,
					"first_name": "Sally",
				},
				{
					"id":         3,
					"first_name": "Jane",
				},
			},
			colKey: "first_name",
			want:   []any{"John", "Sally", "Jane"},
		},
		{
			name: "with intervals",
			src: []map[string]any{
				{
					"id":         1,
					"first_name": "John",
				},
				{
					"id": 2,
				},
				{
					"id":         3,
					"first_name": "Jane",
				},
			},
			colKey: "first_name",
			want:   []any{"John", nil, "Jane"},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			col := Column[string, any](tc.src, tc.colKey)
			assert.Equal(t, tc.want, col)
		})
	}
}

func TestColumnWithFilter(t *testing.T) {
	testCase := []struct {
		name   string
		src    []map[string]any
		colKey string
		want   []any
	}{
		{
			name: "with intervals",
			src: []map[string]any{
				{
					"id":         1,
					"first_name": "John",
				},
				{
					"id": 2,
				},
				{
					"id":         3,
					"first_name": "Jane",
				},
			},
			colKey: "first_name",
			want:   []any{"John", "Jane"},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			col := ColumnWithFilterNotExist[string, any](tc.src, tc.colKey)
			assert.Equal(t, tc.want, col)
		})
	}
}
