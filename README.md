# uuid-clipboard

Generate uuid and automatically save in clipboard

You can download file from `release` or clone this project and do `go run main.go`

![image](https://user-images.githubusercontent.com/11442238/172682857-4cb455c8-781c-41f2-9db0-b1da66e024a9.png)

## Requirement from library

- [clipboard copy](https://github.com/golang-design/clipboard?tab=readme-ov-file#dependency)

## Release

Open **Actions → Release → Run workflow**, select the branch to release
(usually `main`), and enter a version tag (for example, `v1.2.0`).

If the tag exists, the workflow builds its commit. If it is new, the workflow
builds the commit selected when the workflow was started and creates the tag
when publishing the release, after all builds succeed. Existing tags are kept
at their original commits.

The workflow tests and builds the selected commit, uploads the binaries with
`actions/upload-artifact`, downloads them in the release job, and attaches them
to the GitHub Release for that tag.

Releases run manually; pushing a tag does not publish a release.
