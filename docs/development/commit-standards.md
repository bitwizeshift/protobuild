# Commit Standards

This document outlines the basic commit standards that the `protobuild` project
uses.

This document adheres to [RFC-2119]' language such as MUST/SHALL or MAY/SHOULD
to convey the level of requirement.

[RFC-2119]: https://datatracker.ietf.org/doc/html/rfc2119

## TL;DR

I get it, most people don't want to read the long-form. This is the short form.
See the [Rules](#rules) section for the long-form.

* Keep commits small
* Follow [50/72] for commits
* Use imperative present-tense language
* Provide details in messages about the change
* Include a `Change-Category` and optional `Change-Type` trailer in messages

## Rules

### Commits _SHALL_ follow the Git 50/72 rule

All commit messages shall, in general, adhere to the Git [50/72] rule.

The rule, in its simplicity, is that the "title" message must be no longer than
50 characters, and lines in the body should be no longer than 72 characters.
This does not include the (invisible) trailing newline character.

This convention is not an _explicit_ standard, but it helps `git` viewers, IDEs,
`blame` and `diff` tools, and even websites such as [GitHub] itself view the
commit messages without any truncation.

There are some exceptions to this rule that only apply to the _body_ of the
commit message:

* code-snippets that may run longer than 72 characters are not required to
  follow this convention, since it breaks the coherence of the code. However,
  in general, ideally the code is not getting this long to begin with.
* URLs, Links, or other text-content that cannot be broken into multiple lines
  is not required to follow this convention.

### Commits _SHALL_ answer the "What", "Where", "When", and "Why"

This _should_ be a no-brainer since it helps with traceability and understanding
the content long into the future -- but experience shows that this is not often
done by the _"average"_ developer.

Commit messages ultimately are useful for going back-in-time to understand
context of a change. It's _not_ meant to exactly re-state what the change is,
but it _is_ meant to provide the surrounding context that only we as humans know
at the time of implementing the change. As such, it's _extremely useful_ to
provide this context.

In general, this context typically should answer the following 4 basic
questions:

1. **What is this change ultimately doing?** Is there anything we should know
   about what is being done?
2. **Where is this change being made?** Is this a new subsystem, a new concept,
   or a modification of an old one?
3. **When is this change relevant?** Are there conditions that impact or affect
   this change, or that motivated this change? This is often relevant for
   bugs.
4. **Why is this change being made?** What was the motivation for it?

It is not require that every commit answer _all_ of these as independent
sentences. Commits can still be terse; but the context should be clear from the
reader.

**Imagine your reader is a new developer to this project with only a basic
understanding of the structure**.

### Commit titles _SHALL_ be in Imperitive Present-tense

This is the standard for Git, even the `git` utility does this by default.
It's `Merge branch ... into ...` not `Merging`, `Revert` not
`Reverting`, etc.

All commit titles must follow this convention. Aside from the obvious
consistency reasons, because `git` does this implicitly, it means that it will
work better with forming correct git commit message when performing `revert`s,
`cherry-pick`s, `merge`s, etc.

### Commit titles _SHALL_ not end in punctation

This is just for consistency. Don't end the commit title with a punctation.
Commit bodies are fair-game though.

### Commits _SHALL_ indicate a `Change-Category` trailer

All commits are required to include a [Git Trailer] of `Change-Category` which
indicates the category of change being performed. They must be one or more of
the following separated by spaces, which is case-sensitive:

* `infrastructure` - This applies to any tooling or infrastructure that is used
  for the repository to produce artifacts and continue functioning. This may
  include `.github` content such as workflows, actions -- as long as it applies
  towards continuous integration/deployment.
* `community` - This applies to any community-level content, such as dev docs,
  licensing, or guides. It may also apply to community-facing tools, such as
  bots or verification workflows for following standards.
* `library` - This applies to an exported package/Go-library that can be used
  by consumers that are outside of the `protobuild` project.
* `cli` - This applies to the primary `protobuild` tool itself. Any changes that
  impact the public UX, commands, output, etc are all tracked as this.

### Commits _MAY_ indicate a `Change-Type` trailer

Similar to the above, `Change-Type` is a [Git Trailer] which may be used to
indicate specifically what kind of change is being performed. It may be one or
more of the following, separated by spaces:

* `fix` - This is a bugfix of some kind, but it does not have an associated
  [GitHub Issue] that tracked it.
* `feature` - This is the default for any commit that is not labeled with
  anything else.
* `breaking` - This indicates that a change is a BREAKING change, which will
  ultimately correspond to a major version bump. Use of this type should not be
  taken lightly once reaching version 1.0.

If nothing is specified, the change is always assumed to be a `feature` change.

[GitHub Issue]: https://github.com/bitwizeshift/protobuild/issues
[Git Trailer]: https://git-scm.com/docs/git-interpret-trailers
[GitHub]: https://github.com
[50/72]: https://dev.to/noelworden/improving-your-commit-message-with-the-50-72-rule-3g79

### Commits _SHOULD_ be atomic with respect to their change

All commits are expected to be small, and atomic for the change they are making.
This means that commits should not be doing a lot of things at once, but rather
it should be doing small, reversable things that follow separation-of-concerns
principles.

This means not updating the tool's build-engine at the same time as updating
several workflows -- since these can be done piecewise and tackle different
things. Using a `Change-Category` and `Change-Type` will help for this, since if
it's hard to articulate a category/type, your commit may be doing too much.

## Good Example

Below is an example of a simple but useful commit:

```text
Fix panic in protobuild init command

The protobuild init command is panicking due to an off-by-one error in a
slice indexing operation. This is only triggered when the --verbose flag is
provided, since the slice is conditionally logged.

This addresses the problem by guarding the index for accurate output.

Change-Category: cli
Change-Type: fix
```

This tells us several things:

1. The issue being resolved is a `panic`
2. The `panic` is caused by an off-by-one indexing error
3. The scope is limited to when `--verbose` is provided, since it's from a
   conditional log
4. The fix is to guard the index
