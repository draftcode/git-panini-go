# git-panini

`git-panini` is a little tool to treat multiple repositories.

## Problem

- Want to see the statuses of the repositories I checked out.
- Want to keep track of the multiple upstreams.

## Usage

Create `~/.git-panini` like this:

```yaml
# vim: ft=yaml
---
repositories:
  - ~/src/foo
  - ~/src/bar
```

### git-panini fetch

`git-panini fetch` fetches all of the remotes of them:

```
$ git panini fetch
Execute git fetch in all repositories...
/home/draftcode/src/foo
         origin
/home/draftcode/src/bar
         forked
         origin
Complete!
```

### git-panini status

`git-panini status` executes `git status` and shows the differences of the
remote branches and the local branches.

```
$ git panini status
/home/draftcode/src/foo
         M README.md
        master [origin +1]
/home/draftcode/src/bar
         M spec/spec_helper.rb
        add-awesome-feature [origin +2] [forked nobranch]
        master [origin +0] [forked -3]
```

`git-panini status` compares a local branch and a remote branch which has the
same name as the local one. In this example, `bar` repository has
`add-awesome-feature` branch both in local and in `origin`. It seems the local
branch is 2 commits ahead compared to the remote one, so it shows `[origin +2]`.
This shows you have some commits which haven't pushed yet. Since there is no
such branch in `forked`, it shows `[forked nobranch]`. On the other hand,
`master` branch is 3 commits behind from the one in `forked`. Sometimes the
local branch and the remote one is not in the relation of parent and child. In
this case it shows like `[origin diverged]`.

### git-panini find-nonpanini

`git-panini find-nonpanini` finds repositories which aren't listed in
`~/.git-panini`.

```
$ git panini find-nonpanini src
/home/draftcode/src/baz
```
