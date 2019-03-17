# git-walk

Command line to iterate through a git project's history.

# Examples

```
git-walk to start
```
Goes to the first commit, by commit time, in the current history. The current reference is saved into _../.git-walk_ file.

```
git-walk to next
```
Goes to the next commit, chronologically, in the current history.

```
git-walk to end
```
Goes to the saved reference, which is the one saved in _.git-walk_ the first time `git-walk to start` is run.

# How it works

Whenever `git-walk to start` is run, the current reference is saved into the _../.git-walk_ file. The contents of this file allow for checking out commits in the future of the target commit. (_The file is saved in the parent directory due to an issue on [go-git](https://github.com/src-d/go-git/issues/1026) implementation, which removes any untracked files upon checking out a commit._)

`git-walk`, when used with `start` or `next` checks out a commit, so git HEAD becomes detached.

> Note that, one cannot run `git-walk to start` for the first time while HEAD is detached, as a non-detached reference needs to be saved.

# To do

- [ ] Tests
- [ ] Add `--reverse` flag to [go-git](https://github.com/src-d/go-git), so **git-walk** can be done in a smarter way. It would also make it simple to have a `git-walk to next -n 10` to move multiple steps.
- [ ] Improve how reference state is saved. The _.git-walk_ file should not be saved outside of the repository (may require [issue 1026 of go-git](https://github.com/src-d/go-git/issues/1026) to be fixed first). Additionaly, a proper check should be made so it is saved in the root of the repository, not the current directory in which **git-walk** is run.
