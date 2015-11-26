# reportcopter 
[![wercker status](https://app.wercker.com/status/535b33345ba6504c8ae4701aa354db9d/s "wercker status")](https://app.wercker.com/project/bykey/535b33345ba6504c8ae4701aa354db9d)

Why create release changes manualy when you have a well documented control version?
This tool generate changelog reports based on ```git-log``` and the difference between release tags.
If you want to see how the output looks, check the changelog from this project [here](./changelog.md)

## Features

* Regular expresion input for release tags
* Multiple output format
  * JSON
  * Markdown
  * Html

## Install

Download the binaries or install via go

    go get github.com/3zcurdia/reporcopter

## Usage

    $ reporcopter >> changelog.md

### Options

    reporcopter [global options] command [command options] [arguments...]

    COMMANDS:
       help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
        --pattern, -p "v[\d{1,4}\.]{1,}"	Regular expresion for release tags
        --format,  -f "markdown"	      	Output format for report


## Contribue to the project

To contribuer Just follow the next stepts:

* Check out the latest master to make sure the feature hasn't been implemented or the bug hasn't been fixed yet
* Check out the issue tracker to make sure someone already hasn't requested it and/or contributed it
* Fork the project
* Start a feature/bugfix branch
* Commit and push until you are happy with your contribution
* It is desired to add some tests for it.

## License

The MIT License (MIT)
