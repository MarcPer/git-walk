# git-walk

Command line to iterate through a git project's history. Requires [git](https://git-scm.com/) to be installed.

# Examples

- Move to the first commit, by commit time, in the current history:
  ```
  git-walk to start
  ```

- Goes to the next commit, chronologically, in the current history:
  ```
  git-walk to next
  ```

- Goes to the saved reference, which is the one saved in _.git-walk_ the first time `git-walk to start` is run.
  ```
  git-walk to end
  ```

# How it works

**git-walk**, when used with `start` or `next` checks out a commit, so git HEAD becomes [detached](https://git-scm.com/docs/git-checkout#_detached_head).

Whenever `git-walk to start` is run and git HEAD is not detached, the current reference is saved into the _.git-walk_ file. **git-walk** uses it to go back to the reference, even if git HEAD is in detached state.

> Note that one should not run `git-walk to start` for the first time while HEAD is detached, as a non-detached reference needs to be saved. Otherwise, `next` and `end` don't work.

# To do

- [ ] Better error handling. Stderr and exit codes are currently ignored.
- [ ] Allow for moving multiple commits with something like `git-walk to next 10`.
- [ ] Possibly include _.git-walk_ file into _.git/info/exclude_ automatically, so it is not tracked by git.
