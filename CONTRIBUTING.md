# Contributing to Purser

Welcome! We gladly accept contributions from the community. If you wish
to contribute code and you have not signed our contributor license
agreement (CLA), our bot will update the issue when you open a pull
request. For any questions about the CLA process, please refer to our
[FAQ](https://cla.vmware.com/faq).

## Logging Bugs

Anyone can log a bug using the GitHub 'New Issue' button.  Please use
a short title and give as much information as you can about what the
problem is, relevant software versions, and how to reproduce it.  If you
know the fix or a workaround include that too.

## Install dependencies
- Install [git](https://git-scm.com/downloads)
- Install [Go](https://golang.org/dl/) version at least 1.7
- Set GOPATH environment variable. [https://github.com/golang/go/wiki/SettingGOPATH](https://github.com/golang/go/wiki/SettingGOPATH)

## Code Contribution Flow

We use GitHub pull requests to incorporate code changes from external
contributors.  Typical contribution flow steps are:

- Fork the Purser repo into a new repo on GitHub
- Clone the forked repo locally and set the original Purser repo as the upstream repo
- Make changes in a topic branch and commit
- Fetch changes from upstream and resolve any merge conflicts so that your topic branch is up-to-date
- Push all commits to the topic branch in your forked repo
- Submit a pull request to merge topic branch commits to upstream master

If this process sounds unfamiliar have a look at the
excellent [overview of collaboration via pull requests on
GitHub](https://help.github.com/categories/collaborating-with-issues-and-pull-requests) for more information. 

## Coding Style

Our standard for Golang contributions is to match the format of the [standard
Go package library](https://golang.org/pkg).  

- Run `go fmt` on all code.
- All public interfaces, functions, and structs must have complete, grammatically correct Godoc comments that explain their purpose and proper usage.
- Use self-explanatory names for all variables, functions, and interfaces.
- Add comments for non-obvious features of internal implementations but otherwise let the code explain itself.
- Include unit tests for new features and update tests for old ones.

Go is pretty readable so if you follow these rules most functions
will not need additional comments.

### Commit Message Format

We follow the conventions on [How to Write a Git Commit Message](http://chris.beams.io/posts/git-commit/).

Be sure to include any related GitHub
issue references in the commit message.  See [GFM
syntax](https://guides.github.com/features/mastering-markdown/#GitHub-flavored-markdown)
for referencing issues.

### Sign the Contributor License Agreement (CLA)

VMware Apache-licensed projects require all contributors to sign a CLA. 
Visit https://cla.vmware.com and follow steps presented there. 

### Fork the Repo

Navigate to the [Purser repo on
GitHub](https://github.com/vmware/purser) and use the 'Fork' button to
create a forked repository under your GitHub account.  This gives you a copy 
of the repo for pull requests back to purser in https://github.com/<your-username>/purser

### Clone and Set Upstream Remote

Make a local clone of the forked repo and add the base purser
repo as the upstream remote repository.

``` shell
# (go to directory $GOPATH/src/github.com/vmware)
cd $GOPATH/github.com/vmware
# (clone the forked repository)
git clone https://github.com/<your-username>/purser.git
# (go to purser directory)
cd $GOPATH/src/github.com/vmware/purser
# add upstream repository as the original purser repo)
git remote add upstream https://github.com/vmware/purser.git
```

The last git command prepares your clone to pull changes from the
upstream repo and push them into the fork, which enables you to keep
the fork up to date. More on that shortly.

### Make Changes and Commit

Start a new topic branch from the current HEAD position on master and
commit your feature changes into that branch.  

``` shell
git checkout -b foo-api-fix-22 master
# (Make feature changes)
git commit -a --signoff
git push origin foo-api-fix-22
```

The --signoff puts your signature in the commit.  It's required by our CLA
bot. 

It is a git best practice to put work for each new feature in a separate
topic branch and use git checkout commands to jump between them.  This
makes it possible to have multiple active pull requests.  We can accept
pull requests from any branch, so it's up to you how to manage them.

### Stay in Sync with Upstream

From time to time you'll need to merge changes from the upstream
repo so your topic branch stays in sync with other checkins.  To
do so switch to your topic branch, pull from the upstream repo, and
push into the fork.  If there are conflicts you'll need to [merge
them now](https://stackoverflow.com/questions/161813/how-to-resolve-merge-conflicts-in-git).

``` shell
git checkout foo-api-fix-22
git fetch -a
git pull --rebase upstream master --tags
git push --force-with-lease origin foo-api-fix-22
```

The git pull and push options are important.  Here are some details if you 
need deeper understanding. 

- 'pull --rebase' eliminates unnecessary merges
by replaying your commit(s) into the log as if they happened
after the upstream changes.  Check out [What is a "merge
bubble"?](https://stackoverflow.com/questions/26239379/what-is-a-merge-bubble)
for why this is important.  
- --tags ensures that object tags are also pulled
- Depending on your git configuration push --force-with-lease is required to make git update your fork with commits from the upstream repo.

### Create a Pull Request

To contribute your feature, create a pull request by going to the [purser upstream repo on GitHub](https://github.com/vmware/purser) and pressing the 'New pull request' button. 

Select 'compare across forks' and select <your-username>/purser as 'head fork'
and foo-api-fix-22 as the 'compare' branch.  Leave the base fork as 
vmware/purser and master. 

### Wait...

A committer will look the request over and do one of three things: 

- accept it
- send back comments about things you need to fix
- or close the request without merging if we don't think it's a good addition.

### Updating Pull Requests with New Changes

If your pull request needs changes based on code review, 
you'll most likely want to squash the fixes into existing commits.

If your pull request contains a single commit or your changes are related
to the most recent commit, you can simply amend the commit.

``` shell
git add .
git commit --amend
git push --force-with-lease origin foo-api-fix-22
```

If you need to squash changes into an earlier commit, you can use:

``` shell
git add .
git commit --fixup <commit>
git rebase -i --autosquash master
git push --force-with-lease origin foo-api-fix-22
```

Be sure to add a comment to the pull request indicating your new changes
are ready to review, as GitHub does not generate a notification when
you git push.

## Final Words

Thanks for helping us make the project better!
