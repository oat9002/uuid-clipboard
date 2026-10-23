# uuid-clipboard

Generate uuid and automatically save in clipboard

You can download file from `release` or clone this project and do `go run main.go`

![image](https://user-images.githubusercontent.com/11442238/172682857-4cb455c8-781c-41f2-9db0-b1da66e024a9.png)

## Requirement from library

- [clipboard copy](https://github.com/golang-design/clipboard?tab=readme-ov-file#dependency)

## Release

Create and push a version tag, then open **Actions → Release → Run workflow**
and enter the existing tag (for example, `v1.0.0`). The workflow tests and builds
that tag, uploads the binaries with `actions/upload-artifact`, downloads them in
the release job, and attaches them to the GitHub Release for that tag.

Releases run manually; pushing a tag does not publish a release.
