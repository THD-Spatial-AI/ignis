# Contributing

Thank you for taking the time to contribute to **ignis**.

This project welcomes contributions such as bug reports, feature requests, documentation improvements, and code changes.

Please read this guide before opening an issue or submitting a pull request.

## Code of Conduct

By participating in this project, you agree to follow the rules and expectations described in the [Code of Conduct](CODE_OF_CONDUCT.md).

## Ways to Contribute

- Reporting bugs or validation discrepancies
- Requesting features or improvements
- Improving documentation
- Fixing issues or adding country support
- Reviewing pull requests

## Before You Start

- Read the `README.md` to understand the project purpose and setup
- Check existing issues and pull requests to avoid duplicates
- Make sure your idea is relevant to the project scope

## Reporting Bugs and Requesting Changes

Use the [GitHub issue tracker](https://github.com/thd-spatial-ai/ignis/issues) for bug reports, feature requests, and validation issues.

When reporting an issue, please include:

- What you expected to happen
- What actually happened
- Steps to reproduce the issue
- Relevant logs or error messages
- Country and building type (if a calculation issue)
- Go version and OS

## Development Workflow

### 1. Fork and clone

```bash
git clone https://github.com/thd-spatial-ai/ignis.git
cd ignis
```

### 2. Create a branch

```bash
git checkout -b type/short-description
```

Examples: `fix/germany-uvalue-calc`, `feat/add-portugal`, `docs/readme-update`

### 3. Make your changes

Keep changes focused and small where possible.

### 4. Test your changes

```bash
go build ./...
go test ./...
./bin/validate -country germany
```

Ensure validation accuracy is not regressed for existing countries.

### 5. Commit and push

```bash
git add .
git commit -m "Short summary of the change"
git push -u origin <your-branch-name>
```

### 6. Open a pull request

Create a pull request against `main`. Include:

- What changed and why
- Any validation results affected
- Related issue(s) (e.g. `Closes #123`)

## Pull Request Checklist

- [ ] Change is relevant and scoped appropriately
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes (if tests exist)
- [ ] Existing country validation accuracy is not regressed
- [ ] No sensitive information (credentials, API keys) included
- [ ] Documentation updated if usage or behaviour changed

## Commit Message Guidance

Keep messages clear and specific.

Good examples:
- `Fix U-value calculation for cavity walls (Germany)`
- `Add Serbia country support`
- `Update README prerequisites`

## Licensing of Contributions

By contributing to this project, you confirm that your contribution is your own work and you agree it will be licensed under the [MIT License](LICENSE).
