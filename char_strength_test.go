package spg

import (
	"math"
	"testing"
)

// expectation is data for a single entropy test
type expectation struct {
	Length int

	Allow   CTFlag
	Require CTFlag
	Exclude CTFlag

	AllowChars   string
	RequireSets  []string
	ExcludeChars string

	N int64
}

const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lower = "abcdefghijklmnopqrstuvwxyz"
const digits = "0123456789"

var expectations = []expectation{
	{RequireSets: []string{""}, Length: 1, N: 0},
	{RequireSets: []string{"a"}, Length: 0, N: 0},
	{RequireSets: []string{"a"}, Length: 1, N: 1},
	{RequireSets: []string{"a"}, Length: 5, N: 1},
	{RequireSets: []string{"abcde"}, Length: 1, N: 5},
	{RequireSets: []string{"abcde"}, Length: 2, N: 25},
	{RequireSets: []string{"a", "1"}, Length: 2, N: 2},
	{RequireSets: []string{"ab", "123"}, Length: 2, N: 12},

	{Require: Letters | Digits, Length: 0, N: 0},
	{Require: Letters | Digits, Length: 1, N: 0},
	{Require: Letters | Digits, Length: 2, N: 0},
	{Require: Letters | Digits, Length: 3, N: 40560},
	{RequireSets: []string{upper + lower, digits}, Length: 3, N: 96720},

	{Require: Uppers, Length: 3, N: 17576},
	{Allow: Uppers, Length: 3, N: 17576},
	{RequireSets: []string{upper + lower + digits}, Length: 3, N: 238328},
	{Allow: Letters | Digits, Length: 3, N: 238328},

	{RequireSets: []string{"a", "1"}, Length: 1, N: 0},
	{RequireSets: []string{"a"}, AllowChars: "1", Length: 1, N: 1},
	{RequireSets: []string{"a"}, AllowChars: "A1", Length: 2, N: 5},
	{RequireSets: []string{"a", "A"}, AllowChars: "1", Length: 2, N: 2},
	{RequireSets: []string{"a", "A"}, AllowChars: "1", Length: 3, N: 12},

	// Test excluded characters
	{Allow: Uppers, ExcludeChars: "ABC", Length: 1, N: 23},
	{Require: Uppers, ExcludeChars: "ABC", Length: 1, N: 23},
	{RequireSets: []string{"a", "AB"}, AllowChars: "12", ExcludeChars: "B2", Length: 3, N: 12},
}

func TestN(t *testing.T) {
	for i, exp := range expectations {
		recipe := &CharRecipe{
			Length:       exp.Length,
			Allow:        exp.Allow,
			Require:      exp.Require,
			Exclude:      exp.Exclude,
			AllowChars:   exp.AllowChars,
			RequireSets:  exp.RequireSets,
			ExcludeChars: exp.ExcludeChars,
		}
		recipe.buildCharacterList()
		intResult := recipe.n().Int64()

		if intResult != exp.N {
			t.Errorf("%d: result should be %d, was %d", i, exp.N, intResult)
		}
	}
}

func TestEntropy(t *testing.T) {
	recip := &CharRecipe{
		Length:      2,
		AllowChars:  "",
		RequireSets: []string{upper, lower, digits},
	}
	e := recip.Entropy()
	if !math.IsInf(float64(e), -1) {
		t.Errorf("Length is less than number of required sets: entropy should be -Inf, was %f", e)
	}
}

/**
 ** Copyright 2018 AgileBits, Inc.
 ** Licensed under the Apache License, Version 2.0 (the "License").
 **/
