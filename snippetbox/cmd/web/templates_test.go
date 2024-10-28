package main

import (
	"testing"
	"time"

	"snippetbox.xscotophilic.art/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "UTC",
			input:  time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			output: "17 Mar 2022 at 10:15",
		},
		{
			name:   "Empty",
			input:  time.Time{},
			output: "",
		},
		{
			name:   "CET",
			input:  time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			output: "17 Mar 2022 at 09:15",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				hd := humanDate(test.input)

				assert.Equal(t, hd, test.output)
			},
		)
	}
}
