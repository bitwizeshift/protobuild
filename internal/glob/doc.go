/*
Package glob provides a better globbing mechanism than the builtin filepath.Glob
functionality. This extends it by supporting both `**` substitutions for
any number of directories, as well as `!` negation solutions that can be
used to invert a selection and deny values.
*/
package glob
