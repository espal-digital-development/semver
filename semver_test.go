package semver_test

import (
	"testing"

	"github.com/espal-digital-development/semver"
	"github.com/juju/errors"
)

var (
	validVersions = []string{
		"0.0.0",
		"1.2.3",
		"11.3.9",
		"192.10.999",
		"1.9.18-hotfix",
		"1.0.5-master",
		"10.5.7-with-multiple-dashes",
		"10.5.7-dashes-and.dots",
	}
	invalidVersions = []string{
		"9",
		"2.0",
		"8.",
		"1.2.",
		"9.7.0.",
		"3.8.2-",
	}
	greaterThanVersions = [][]string{
		{"1.2.3", "0.0.1"},
		{"1.2.3", "1.2.3"},
		{"1.1.0", "1.0.0"},
		{"1.0.1", "1.0.0"},
		{"2.12.13", "1.33.44"},
		{"2.12.13", "1.33.44"},
	}
	smallerThanVersions = [][]string{
		{"0.0.1", "0.0.2"},
		{"0.2.1", "0.2.2"},
		{"0.1.2", "0.2.0"},
		{"0.20.19", "1.1.2"},
	}
	inRangeVersions = [][]string{
		{"12.13.14", "11.22.33", "13.44.55"},
		{"12.13.14", "11.22.33", ""},
		{"12.13.14", "0.0.1", "12.22.33"},
		{"11.22.33", "11.22.33", "13.22.33"},
		{"12.13.14", "11.22.33", "13.22.33"},
		{"11.22.33", "11.22.33", "13.22.33-hotfix"},
	}
	outOfRangeVersions = [][]string{
		{"1.13.14", "11.22.33", "13.44.55"},
		{"1.13.14", "11.22.33", ""},
		{"1.13.14", "10.0.1", "12.22.33"},
		{"11.22.32", "11.22.33", "11.22.33"},
		{"11.22.33", "11.22.31", "11.22.32-hotfix"},
	}
)

func TestNew(t *testing.T) {
	semver, err := semver.New()
	if err != nil {
		t.Fatal(err)
	}
	if semver == nil {
		t.Fatal("expected semver to not be nil")
	}
}

func TestValid(t *testing.T) {
	for k := range validVersions {
		version := validVersions[k]
		t.Run("valid-"+version, func(t2 *testing.T) {
			semver, err := semver.New()
			if err != nil {
				t2.Fatal(err)
			}
			if !semver.Valid(version) {
				t2.Fatalf("expecting `%s` to be valid", version)
			}
		})
	}
}

func TestInvalid(t *testing.T) {
	for k := range invalidVersions {
		version := invalidVersions[k]
		t.Run("invalid-"+version, func(t2 *testing.T) {
			semver, err := semver.New()
			if err != nil {
				t2.Fatal(err)
			}
			if semver.Valid(version) {
				t2.Fatalf("expecting `%s` to be invalid", version)
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	for k := range greaterThanVersions {
		version := greaterThanVersions[k][0]
		compare := greaterThanVersions[k][1]
		t.Run("greater-than-"+version+"_"+compare, func(t2 *testing.T) {
			semver, err := semver.New()
			if err != nil {
				t2.Fatal(err)
			}
			greaterThan, err := semver.GreaterThanOrEqual(version, compare)
			if err != nil {
				t2.Fatal(err)
			}
			if !greaterThan {
				t2.Fatalf("expect `%s` to be greater than `%s`", version, compare)
			}
		})
	}
}

func TestSmallerThan(t *testing.T) {
	for k := range smallerThanVersions {
		version := smallerThanVersions[k][0]
		compare := smallerThanVersions[k][1]
		t.Run("smaller-than-"+version+"_"+compare, func(t2 *testing.T) {
			semver, err := semver.New()
			if err != nil {
				t2.Fatal(err)
			}
			smallerThan, err := semver.SmallerThanOrEqual(version, compare)
			if err != nil {
				t2.Fatal(err)
			}
			if !smallerThan {
				t2.Fatalf("expect `%s` to be smaller than `%s`", version, compare)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	for k := range inRangeVersions {
		version := inRangeVersions[k][0]
		from := inRangeVersions[k][1]
		to := inRangeVersions[k][2]
		t.Run("in-range-"+version+"_"+from+"_"+to, func(t2 *testing.T) {
			semver, err := semver.New()
			if err != nil {
				t2.Fatal(err)
			}
			inRange, err := semver.InRange(version, from, to)
			if err != nil {
				t2.Fatal(err)
			}
			if !inRange {
				t2.Fatalf("expect `%s` to be between `%s` and `%s`", version, from, to)
			}
		})
	}
}

func TestOutOfRange(t *testing.T) {
	for k := range outOfRangeVersions {
		version := outOfRangeVersions[k][0]
		from := outOfRangeVersions[k][1]
		to := outOfRangeVersions[k][2]
		t.Run("out-of-range-"+version+"_"+from+"_"+to, func(t2 *testing.T) {
			semver, err := semver.New()
			if err != nil {
				t2.Fatal(err)
			}
			inRange, err := semver.InRange(version, from, to)
			if err != nil {
				t2.Fatal(err)
			}
			if inRange {
				t2.Fatalf("expect `%s` to not be between `%s` and `%s`", version, from, to)
			}
		})
	}
}

func TestInRangeErrors(t *testing.T) {
	semver, err := semver.New()
	if err != nil {
		t.Fatal(err)
	}
	wrongVersion := invalidVersions[0]
	expectedErr := errors.Errorf("version `%s` is invalid", wrongVersion)
	_, err = semver.InRange(wrongVersion, "0.0.1", "0.0.1")
	if err == nil || err == expectedErr {
		t.Fatalf("expected error to be thrown `%s`", expectedErr.Error())
	}

	_, err = semver.InRange("0.0.1", wrongVersion, wrongVersion)
	if err == nil || err == expectedErr {
		t.Fatalf("expected error to be thrown `%s`", expectedErr.Error())
	}

	_, err = semver.InRange("0.0.1", "0.0.1", wrongVersion)
	if err == nil || err == expectedErr {
		t.Fatalf("expected error to be thrown `%s`", expectedErr.Error())
	}
}
