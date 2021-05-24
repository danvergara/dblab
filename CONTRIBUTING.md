# How to contribute to dblab

Thank you for considering contributing to dblab!

## First time setup

- Fork dblab to your GitHub account by clicking the [Fork](https://github.com/danvergara/dblab/fork) button.
- [Clone](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo#step-2-create-a-local-clone-of-your-fork) the main repository locally.

```bash
git clone https://github.com/danvergara/dblab.git
cd dblab
```

- Add your fork as a remote to push your work to. Replace ``{username}`` with your username. This names the remote "fork", the  default dblab remote is "origin".

```bash
git remote add fork https://github.com/{username}/dblab
```

## Start coding

- Create a branch to identify the issue, feature addition or change you would like to work on.

```bash
git fetch origin
git checkout -b your-branch-name origin/main
```

- Using your favorite editor, make your changes.
- Include tests that cover any code changes you make. Make sure the
  test fails without your patch. Run the tests as described below.
- Push your commits to your fork on GitHub and
  [create a pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request). Link to the issue being addressed with
  ``fixes #123`` in the pull request in the case that you are working on an issue.

```bash
git push --set-upstream fork your-branch-name
```

## Running the tests

Run the test suite with make.

```bash
make test
```

This runs the tests. You can check all the options with `help` command.

```bash
Usage:
test           Runs the tests
unit-test      Runs the tests with the short flag
int-test       Runs the integration tests
linter         Runs the colangci-lint command
test-all       Runs the integration testing bash script with different database docker image versions
docker-build   Builds de Docker image
build          Builds the Go program
run            Runs the application
up             Runs all the containers listed in the docker-compose.yml file
down           Shut down all the containers listed in the docker-compose.yml file
help           Prints this help message
```