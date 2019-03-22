# git-walk

Command line to iterate through a git project's history, going forward in time (from oldest to newest commit). Requires [git](https://git-scm.com/) to be installed.

# Installation

- Download latest binary from [releases page](https://github.com/MarcPer/git-walk/releases), choosing the right platform.
- Place binary into a `$PATH` folder.

Below we assume the binary is renamed to `git-walk`.

# Examples

- Move to the first commit, by commit time, in the current history:
  ```
  git-walk to start
  ```

- Checkout a given commit, while allowing navigation to commits created after it:
  ```
  git-walk to start <commit>
  ```

- Go to the commit created after the current one:
  ```
  git-walk to next
  ```

- Go to last known [non-detached HEAD](https://git-scm.com/docs/git-checkout#_detached_head), saved into _.git-walk_ file when `git-walk to start` is run.
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
