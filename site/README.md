## Site Repo

This directory contains an example of a site repo for Gru.

A site repo in Gru is essentially a Git repository with branches,
where each branch maps to an environment, which can be used by
remote minions.

In order to use this site repo, simply copy the contents of this
directory and add them to a Git repository, which you can use by
your minions.

```bash
$ cp -a site ~/gru-site
$ cd ~/gru-site
$ git init
$ git add
$ git commit -m 'Initial commit of site repo'
$ git checkout -b production
```

Once you've got the site repo in Git you can start you minions by
pointing them to your site repo, e.g.

```bash
$ gructl serve --siterepo https://github.com/you/gru-site
```
