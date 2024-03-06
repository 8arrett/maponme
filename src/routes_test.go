package main

import (
	"testing"
)

// --------------------------------------------------------
//   Validates username
// --------------------------------------------------------

func TestSmallUsername(t *testing.T) {
	names := []string{"", "a", "ab", "abc", "abcd", "abcde", "abcdef", "abcdefg"}
	want := false

	for _, name := range names {
		msg := validUserFormat(name)
		if want != msg {
			t.Fatalf(`validUserFormat("%s") = %t, want match for %t`, name, msg, want)
		}
	}
}

func TestLargeUsername(t *testing.T) {
	names := []string{"012345678", "0123456789", "01234567890", "012345678901", "0123456789012", "01234567890123", "012345678901234", "0123456789012345", "20202020202020202020"}
	want := false

	for _, name := range names {
		msg := validUserFormat(name)
		if want != msg {
			t.Fatalf(`validUserFormat("%s") = %t, want match for %t`, name, msg, want)
		}
	}
}

func TestGoodUsername(t *testing.T) {
	names := []string{"abcdefgh", "ABCDEFGH", "01234567", "--------", "========", "ijklmnop", "qrstuvwx", "yz--==zy", "IJKLMNOP", "QRSTUVWX", "YZ==--ZY"}
	want := true

	for _, name := range names {
		msg := validUserFormat(name)
		if want != msg {
			t.Fatalf(`validUserFormat("%s") = %t, want match for %t`, name, msg, want)
		}
	}
}

func TestBadUsername(t *testing.T) {
	names := []string{
		"áááááááá",
		"........",
		"////////",
		"世界世界世界世界",
		"\x00\x00\x00\x00\x00\x00\x00\x00",
		"\xff\xff\xff\xff\xff\xff\xff\xff",
	}
	want := false

	for _, name := range names {
		msg := validUserFormat(name)
		if want != msg {
			t.Fatalf(`validUserFormat("%s") = %t, want match for %t`, name, msg, want)
		}
	}
}

// --------------------------------------------------------
//   Validates key input
// --------------------------------------------------------

func TestGoodPasskey(t *testing.T) {
	inputs := []string{
		"aaaaaaaaaabcdefghijklmno",
		"pqrstuvwxyzaaaaaaaaaaaaa",
		"AAAAAAAAAABCDEFGHIJKLMNO",
		"PQRSTUVWXYZAAAAAAAAAAAAA",
		"000000000123456789000000",
		"+++++++++/////////++++++",
	}
	want := true

	for _, input := range inputs {
		msg := validUserKey(input)
		if want != msg {
			t.Fatalf(`validUserKey("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}

func TestPasskeySize(t *testing.T) {
	inputs := []string{
		"aaaaaaaaaabcdefghijklmnoB",
		"pqrstuvwxyzaaaaaaaaaaaaaBC",
		"AAAAAAAAAABCDEFGHIJKLMNOBCD",
		"PQRSTUVWXYZAAAAAAAAAAAAABCDE",
		"00000000012345678900000",
		"+++++++++/////////++++",
		"",
	}
	want := false

	for _, input := range inputs {
		msg := validUserKey(input)
		if want != msg {
			t.Fatalf(`validUserKey("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}

func TestBadPasskey(t *testing.T) {
	inputs := []string{
		"áááááááááááááááááááááááá",
		"************************",
		"------------------------",
		"!!!!!!!!!!!!!!!!!!!!!!!!",
		"........................",
		"世界世界世界世界世界世界世界世界世界世界世界世界",
		"\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\",
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff",
	}
	want := false

	for _, input := range inputs {
		msg := validUserKey(input)
		if want != msg {
			t.Fatalf(`validUserKey("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}

// --------------------------------------------------------
//   Validates location input
// --------------------------------------------------------

func TestGoodLocations(t *testing.T) {
	inputs := []string{
		"1.12345678901234",
		"1.",
		"12.1234567",
		"123.1",
		"123.123456789012345",
		"-1.12345678901234",
		"-1.",
		"-12.1234567",
		"-123.1",
		"-123.123456789012345",
	}
	want := true

	for _, input := range inputs {
		msg := validLocation(input)
		if want != msg {
			t.Fatalf(`validLocation("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}

func TestBadLocations(t *testing.T) {
	inputs := []string{
		"1.1234567890123456",
		"1.12345678901234567",
		"1234.",
		"1234.12345",
		"1-23456789012345",
		"1.23456789012345-",
		"1.23456-789012345",
		"+1.23456789012345",
		"1..2345678901234",
		"1",
		"",
	}
	want := false

	for _, input := range inputs {
		msg := validLocation(input)
		if want != msg {
			t.Fatalf(`validLocation("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}

func TestGoodNumbers(t *testing.T) {
	inputs := []string{
		"0",
		"000",
		"0001",
		"123456789012345",
		"1234567890123456",
	}
	want := true

	for _, input := range inputs {
		msg := validNumber(input)
		if want != msg {
			t.Fatalf(`validNumber("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}

func TestBadNumbers(t *testing.T) {
	inputs := []string{
		"",
		"A",
		"a",
		"Z",
		"z",
		".",
		"+",
		"-",
		"12345678901234567",
		"123456789012345678",
	}
	want := false

	for _, input := range inputs {
		msg := validNumber(input)
		if want != msg {
			t.Fatalf(`validNumber("%s") = %t, want match for %t`, input, msg, want)
		}
	}
}
