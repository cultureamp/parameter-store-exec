# Maintenance

## Cutting a New Release

Releases are automated via GitHub Actions using [GoReleaser](https://goreleaser.com/). A release is triggered by pushing a Git tag that follows semantic versioning with a `v` prefix.

### Steps

1. **Ensure `master` is up to date:**

   ```bash
   git checkout master
   git pull
   ```

2. **Run tests locally (optional but recommended):**

   ```bash
   make test
   ```

3. **Create a version tag:**

   ```bash
   git tag v<MAJOR>.<MINOR>.<PATCH>
   ```

   For example: `git tag v1.3.0`

4. **Push the tag to GitHub:**

   ```bash
   git push origin v<MAJOR>.<MINOR>.<PATCH>
   ```

5. **Verify the release:**

   The [release workflow](/.github/workflows/release.yaml) will run automatically. It builds binaries for all configured platforms and publishes them as a GitHub Release. Monitor progress in the [Actions tab](../../actions/workflows/release.yaml).

### Versioning Guidelines

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** — breaking changes (e.g. CLI flag removals, incompatible behaviour changes)
- **MINOR** — new features or capabilities, backwards-compatible
- **PATCH** — bug fixes and dependency updates
