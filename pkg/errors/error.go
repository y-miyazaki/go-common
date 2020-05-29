package errors

import "golang.org/x/xerrors"

// New uses xerrors.New.
func New(format string) error {
	return xerrors.New(format)
}

// Wrap uses xerrors.Errorf  need to add ": %w".
func Wrap(format string, a ...interface{}) error {
	return xerrors.Errorf(format+": %w", a)
}

// UnWrap uses xerrors.Unwrap.
func Unwrap(err error) error {
	return xerrors.Unwrap(err)
}

// Is uses xerrors.Is().
func Is(err, target error) bool {
	return xerrors.Is(err, target)
}

// As uses xerrors.As().
func As(err error, target interface{}) bool {
	return xerrors.As(err, target)
}
