package main

import (
	"context"
	"testing"

	"github.com/stainour/test12/internal/counter"
	"github.com/stainour/test12/internal/reader"
)

func Test(t *testing.T) {
	r, _ := reader.NewFileReader("../../test_files")
	_, _ = counter.NewASCIICounter(r).Count(context.Background())
}
