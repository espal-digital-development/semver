package semver

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

var _ Versioning = &Semver{}

const (
	expectedChunksWithTag      = 2
	exptectedPartsWithRevision = 3
)

// Versioning represents an object that provides validation tools to check a versioning system's versions.
type Versioning interface {
	Valid(version string) bool
	InRange(version string, start string, end string) (bool, error)
	GreaterThanOrEqual(version string, compare string) (bool, error)
	SmallerThanOrEqual(version string, compare string) (bool, error)
}

// Semver validator to do checks based on the semver 2.0.0 spec.
type Semver struct {
	reValid *regexp.Regexp
}

// IsValid checks if the given version is a valid semver format.
func (s *Semver) Valid(version string) bool {
	return s.reValid.MatchString(version)
}

type semVersion struct {
	major    int
	minor    int
	revision int
	tag      string
}

// InRange checks if the version is between the given start and end versions.
func (s *Semver) InRange(version string, start string, end string) (bool, error) {
	greaterThanOrEqual, err := s.GreaterThanOrEqual(version, start)
	if err != nil {
		return false, errors.Trace(err)
	}
	smallerThanOrEqual := true
	if end != "" {
		smallerThanOrEqual, err = s.SmallerThanOrEqual(version, end)
		if err != nil {
			return false, errors.Trace(err)
		}
	}
	return greaterThanOrEqual && smallerThanOrEqual, nil
}

// GreaterThanOrEqual checks if the given version is greater than or equal to the compare version.
func (s *Semver) GreaterThanOrEqual(version string, compare string) (bool, error) {
	if !s.Valid(version) {
		return false, errors.Errorf("version `%s` is invalid", version)
	}
	if !s.Valid(compare) {
		return false, errors.Errorf("compare `%s` is invalid", compare)
	}
	semVersion, err := s.buildVersion(version)
	if err != nil {
		return false, errors.Trace(err)
	}
	semCompare, err := s.buildVersion(compare)
	if err != nil {
		return false, errors.Trace(err)
	}
	if semVersion.major < semCompare.major {
		return false, nil
	}
	if semVersion.major > semCompare.major {
		return true, nil
	}
	if semVersion.minor < semCompare.minor {
		return false, nil
	}
	if semVersion.minor > semCompare.minor {
		return true, nil
	}
	return semVersion.revision >= semCompare.revision, nil
}

// SmallerThanOrEqual checks if the given version is smaller than or equal to the compare version.
func (s *Semver) SmallerThanOrEqual(version string, compare string) (bool, error) {
	if !s.Valid(version) {
		return false, errors.Errorf("version `%s` is invalid", version)
	}
	if !s.Valid(compare) {
		return false, errors.Errorf("compare `%s` is invalid", compare)
	}
	semVersion, err := s.buildVersion(version)
	if err != nil {
		return false, errors.Trace(err)
	}
	semCompare, err := s.buildVersion(compare)
	if err != nil {
		return false, errors.Trace(err)
	}
	if semVersion.major > semCompare.major {
		return false, nil
	}
	if semVersion.major < semCompare.major {
		return true, nil
	}
	if semVersion.minor > semCompare.minor {
		return false, nil
	}
	if semVersion.minor < semCompare.minor {
		return true, nil
	}
	return semVersion.revision <= semCompare.revision, nil
}

func (s *Semver) buildVersion(version string) (*semVersion, error) {
	semVersion := &semVersion{}
	if strings.Contains(version, "-") {
		chunks := strings.SplitN(version, "-", 2)
		if len(chunks) != expectedChunksWithTag {
			return nil, errors.Errorf("versions with a tag should be 2 chunks. Got %d", len(chunks))
		}
		semVersion.tag = chunks[1]
		version = chunks[0]
	}
	versionParts := strings.Split(version, ".")
	versionPartsLength := len(versionParts)
	if versionPartsLength != 2 && versionPartsLength != 3 {
		return nil, errors.Errorf("versions should be 2 or 3 parts. Got %d", versionPartsLength)
	}
	var err error
	semVersion.major, err = strconv.Atoi(versionParts[0])
	if err != nil {
		return nil, errors.Trace(err)
	}
	semVersion.minor, err = strconv.Atoi(versionParts[1])
	if err != nil {
		return nil, errors.Trace(err)
	}
	if versionPartsLength == exptectedPartsWithRevision {
		semVersion.revision, err = strconv.Atoi(versionParts[2])
		if err != nil {
			return nil, errors.Trace(err)
		}
	}
	return semVersion, nil
}

// New returns a new instance ofSemver.
func New() (*Semver, error) {
	s := &Semver{}
	var err error
	s.reValid, err = regexp.Compile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-]` +
		`[0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return s, nil
}
