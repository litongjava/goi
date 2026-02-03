# goi
goi=go install
A concise command-line tool: clone a project from a repository URL into `~/code/<repo>`, optionally switch branches, then pull the latest code and run `build`.

## Features

* Automatically parses the repository name and places it under `~/code`
* Automatically runs `git clone` if the repository does not exist locally
* If the repository already exists, runs `git reset --hard` to keep the working tree clean
* Optional branch switching
* Pulls the latest code and then runs `build`

## Installation

```bash
go install github.com/litongjava/goi@latest
```

## Usage

```bash
goi <repo_url> [-b branch]
```

Examples:

```bash
goi https://github.com/user/project.git
goi git@github.com:user/project.git -b main
```

## Behavior Details

* The code is placed in `~/code/<repo>` by default.
* If the directory already exists, `git reset --hard` is executed, which discards any local uncommitted changes.
* If `-b` is specified, `git checkout <branch>` is executed.
* `git pull` is run to fetch the latest code.
* Finally, the `build` command is executed. Make sure a usable `build` command is available in your environment. For details about the `build` command, refer to the documentation at [https://github.com/litongjava/go-build/](https://github.com/litongjava/go-build/).

## Exit Codes

If any step fails, the program exits immediately and prints an error message.
