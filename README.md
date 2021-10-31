# Semver Bumper

Semver Bumper is a tool
to bump [semantic versions](https://semver.org/)
stored in a git repository as tags
by evaluating commit messages
between releases.

The default configuration of Semver Bumper 
is based on [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
It can be changed using command line arguments or a configuration file.

Semver bumper itself has [a configuration file](.semver-bumper.conf.yaml)
to make bump decisions based on emojis.

## Help

```
$ semver-bumper -h

Usage:
  bumper [OPTIONS]

Application Options:
  -C, --config-file=               load parameters from a JSON or YAML file
  -p, --pre=                       bump prerelease with given keyword, eg "rc" for "1.2.3-rc.4"'
  -t, --tag-prefix=                only detect tags matching the expression, eg "v" for "v1.2.3"
  -n, --no-match-bump=[none|patch] bump patch or nothing when no commits match
  -o, --output=                    write result into file, defaults to stdout
  -c, --commits=                   write commit messages into file
  -i, --path-include=              only detect commits at the given path, can be supplied multiple times
  -x, --path-exclude=              ignore commits at the given path, can be supplied multiple times
  -0, --initial-version=           release version if there are no tags yet, defaults to "1.0.0"
  -1, --major=                     commit message keywords justifying a major version bump, can be supplied multiple times
  -2, --minor=                     commit message keywords justifying a minor version bump, can be supplied multiple times
  -3, --patch=                     commit message keywords justifying a patch version bump, can be supplied multiple times
  -k, --print-keywords             print the configured version bump keywords and exit
  -W, --write-config=              write the given parameters into a JSON or YAML config file and exit

Help Options:
  -h, --help                       Show this help message
```
