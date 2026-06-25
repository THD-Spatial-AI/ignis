# Open Source Readiness Checklist

Use this checklist before making a repository public under the **THD-Spatial** group.

This checklist is designed to be practical and easy to review. It separates **required**, **conditional**, and **recommended** items so you can focus on what matters first.

---

## Review Summary

| Section | Requirement Level | Status | Notes |
| ------- | ----------------- | ------ | ----- |
| [Essential Requirements](#essential-requirements) | Required | Complete | |
| [Data Management (Git LFS)](#data-management-required-if-applicable) | Required if applicable | Complete | `data/tabula-calculator.xlsx` tracked via LFS |
| [Attribution](#attribution-required-if-applicable) | Required if applicable | Complete | TABULA/IWU attributed in `ATTRIBUTIONS.md` |
| [Citation](#citation-recommended-for-research-projects) | Recommended | Complete | `CITATION.cff` present |
| [Optional but Recommended Files](#optional-but-recommended-files-recommended) | Recommended | Partial | Issue templates present; CHANGELOG, CODEOWNERS, SECURITY, SUPPORT pending |
| [Quality Checks](#quality-checks-required) | Required | Partial | Repository settings to be configured after GitHub push |
| [Final Steps](#final-steps-required) | Required | Pending | Verify after push |

---

## Essential Requirements

Make sure the following essential requirements are completed before making your repository public.

Quick links: [LICENSE](#license) · [README](#readme) · [CONTRIBUTING](#contributing) · [CODE_OF_CONDUCT](#code_of_conduct)

### LICENSE

- [x] `LICENSE` file present in repository root
- [x] Appropriate license chosen — MIT
- [x] License committed to repository

### README

- [x] `README.md` file present in repository root

#### README content

- [x] Project title and short description included
- [x] Project purpose clearly explained
- [x] Key features listed
- [x] Installation/setup instructions provided
- [x] Usage examples included
- [x] Contribution guidance referenced (`CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`)
- [x] README uses clear Markdown formatting

### CONTRIBUTING

- [x] `CONTRIBUTING.md` file present in repository root
- [x] Issue reporting process documented
- [x] Pull request submission process documented
- [x] Coding standards / best practices outlined (`go build`, `go test`, validation regression)
- [x] Commit message guidance included
- [x] `CONTRIBUTING.md` linked from `README.md`

### CODE_OF_CONDUCT

- [x] `CODE_OF_CONDUCT.md` file present in repository root
- [x] Contributor Covenant 3.0
- [x] Contact method for reporting issues included
- [x] `CODE_OF_CONDUCT.md` linked from `README.md`
---

## Data Management (Required if applicable)

If your repository contains large files (e.g. datasets, binaries, media, generated assets), you must use Git LFS to manage them.

### Git LFS (Large File Storage)

- [x] Git LFS installed locally
- [x] Git LFS initialised in repository (`git lfs install`)
- [x] Large file types tracked — `*.xlsx`, `*.csv`, archives, media (`data/tabula-calculator.xlsx`)
- [x] `.gitattributes` committed to repository
- [x] Large files added and committed properly
- [ ] Release strategy for LFS files documented (if applicable)

> **Note:** Files tracked by Git LFS are not included in release assets by default.
> If needed, enable **Include Git LFS objects in archives** in repository settings.

---

## Attribution (Required if applicable)

If your project uses third-party components, assets, or generated code that require attribution (e.g. UI components from Figma, design systems, libraries under Apache 2.0, BSD, or CC BY licenses), include an `ATTRIBUTIONS.md` file.

### ATTRIBUTION

- [x] `ATTRIBUTIONS.md` file created in repository root
- [x] TABULA building typology data (IWU) listed
- [x] Source platform or author identified
- [x] License type specified — CC BY 4.0
- [x] Links to original sources included
- [x] Go dependencies noted

!!! note

    If no third-party assets requiring attribution are used, this section can be skipped.

---

## Citation (Recommended for research projects)

If your project is related to research or uses open-source tools and data that should be cited, include a `CITATION.cff` file. GitHub automatically detects this file and shows a "Cite this repository" button on the repository page.

For full specification details, see [Citation File Format](https://citation-file-format.github.io/).

### CITATION.cff

- [x] `CITATION.cff` file created in repository root
- [x] `cff-version` set to `1.2.0`
- [x] Project title specified — ignis
- [x] Author listed — THD-Spatial-AI Research Group, Technische Hochschule Deggendorf
- [ ] ORCID included for each author (add when individual authors are listed)
- [x] License field matches `LICENSE` — MIT
- [x] Repository URL included
- [x] Version and release date set
- [ ] File validates against the CFF schema (use [cff-validator](https://github.com/citation-file-format/cff-initializer-javascript))

!!! tip

    You can generate a `CITATION.cff` file interactively using the [CFF Initializer](https://citation-file-format.github.io/cff-initializer-javascript/).

---

## Optional but Recommended Files (Recommended)

These files are not always required, but they improve project quality, collaboration, and maintainability.

### CHANGELOG

- [ ] `CHANGELOG.md` file created (if using versioned releases)
- [ ] Version history documented
- [ ] Changes formatted consistently (e.g. [Keep a Changelog](https://keepachangelog.com/))

### Issue and PR Templates

- [x] `.github/ISSUE_TEMPLATE/` directory created (if using issue templates)
- [x] Issue template(s) added
- [ ] `.github/pull_request_template.md` file added
- [ ] PR template includes a contributor checklist (recommended)

### Code Owners

- [ ] `CODEOWNERS` file created (root or `.github/`)
- [ ] Owners assigned for relevant parts of the repository
- [ ] Teams/individuals correctly specified

### Security Policy

- [ ] `SECURITY.md` file created
- [ ] Vulnerability reporting process documented
- [ ] Contact method for security issues provided
- [ ] Supported versions documented (if applicable)

### Support Guidelines

- [ ] `SUPPORT.md` file created
- [ ] Support channels documented (e.g. email, discussions, issue tracker)
- [ ] Response expectations documented (if applicable)

---

## Quality Checks (Required)

### Documentation Quality

- [x] Documentation reviewed for clarity
- [x] Spelling and grammar checked
- [x] Links tested and working
- [x] Code examples/commands tested (if applicable)

### Repository Settings

- [ ] Repository description added on GitHub
- [ ] Topics/tags added for discoverability
- [ ] Repository visibility set to **Public**
- [ ] Default branch configured (typically `main`)
- [ ] Branch protection rules configured (if needed)

### Code Quality / Safety

- [x] Code reviewed and cleaned up
- [x] Sensitive information removed (API keys, passwords, credentials)
- [x] `.gitignore` configured properly
- [x] Dependencies documented
- [x] Build/test instructions included (if applicable)

---

## Final Steps (Required)

- [ ] All essential requirements completed
- [ ] Repository tested by cloning a fresh copy
- [ ] Links and references verified
- [ ] Project builds/runs successfully from scratch (if applicable)
- [ ] Team members notified about repository availability (if applicable)

---

## Repository Ready for Open Source

Once all required items are completed, the repository meets the THD-Spatial-AI minimum readiness standard for public release.

- **Date completed:** `14.04.2026`
- **Reviewed by:** [`Jay Ravani`](https://github.com/jravani)
