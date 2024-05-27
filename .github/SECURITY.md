# Security Policy

## Reporting a Security Vulnerability

This project has a small surface-area for general vulnerability exploitation,
since it does not directly process any user data in any way, nor does it work
with encrypted

As such, if you discover a security vulnerability within this project, please
just open a [Github Issue] using the [Security Vulnerability] template.
I appreciate your responsible disclosure and will make every effort to quickly
address the issue.

When doing so, please provide the following details in your email:

* Your affiliation (if applicable).
* A detailed description of the vulnerability, including information on how to
  reproduce it.
* Any potential impact of the vulnerability.

Once we receive your report, we will acknowledge the receipt of the
vulnerability within **1 week**, and strive to provide regular updates on our
progress.

[Github Issue]: https://github.com/bitwizeshift/protobuild/issues/new
[Security Vulnerability]:

## Security Practices

### Static CodeQL Analysis

We use CodeQL to ensure the overall quality code integrity within the
codebase. CodeQL provides a powerful static analysis tool that aids in finding
security-related issues early in the development process.

### License Scanning

We use [`go-licenses`] to ensure that all 3rd-party dependencies use an
OSI-approved license that is compatible with this project, and to generate a
manifest of third-party dependencies.

**Note:** If you believe there may be an issue with attribution in `protobuild`,
please open a [Github Issue] and tell us what is missing! We make no effort to
misrepresent the origin of software, and are thankful for all projects which we
rely on -- and want to be responsible about attributing everyone accordingly.

[`go-licenses`]: https://github.com/google/go-licenses

## Updates and Notifications

Security updates and notifications will be provided through the project's
GitHub repository. Users are encouraged to watch the repository for any
announcements regarding security releases or patches.

## Supported Versions

We prioritize addressing security vulnerabilities in the latest stable release
of the project. Users are strongly encouraged to keep their installations
up-to-date with the latest releases.

## License

By participating in this responsible disclosure process, you agree that your
actions comply with applicable laws and regulations. We appreciate your
contribution to the security of this project.

Thank you for helping to keep this project secure!
