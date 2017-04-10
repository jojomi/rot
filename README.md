# rot

[![Go Report Card](https://goreportcard.com/badge/github.com/jojomi/rot)](https://goreportcard.com/report/github.com/jojomi/rot)

You have that file on your hard disk (oh SSD these days, pardon me) that you might need in a couple hours or a few days. You know you probably won't need it again, but what if? What if you have to install your Linux server yet again tomorrow because you already screwed it up? Would you want to redownload that huge ISO file? Rather not, but you will forget about that file sitting in your downloads folder occupying all that precious expansive storage space. If you are lucky, you will find it hogging your resources unnecessarily when cleaning up in a year thinking why the hell you never got around to delete that hopelessly outdated and outrageously big file. You are human! And you are not alone. I feel with you. Let's have someone (something?) else take care.

`rot` is a tool that will take care of removing files and folders you know you will not need anymore after a certain time. You can specify either a number of hours or days or a rot date after which `rot` can safely remove it from your disk.


# Usage

```
rot empowers you to stage files and folders for rotting (later deletion).

Usage:
  rot [command]

Available Commands:
  add         Add files or folders for rotting
  clean       Clean files and folders that are rotten
  help        Help about any command
  list        List the files and folders current rotting
  stop        Stops a file or folder from rotting
```

For all usage details call `rot --help` or add `--help` to any of the subcommands listed.

You could use `rot clean --dry-run` for monitoring if there are files that should be deleted. If you are brave enough to automate it more, run the command `rot clean` periodically. If you want to monitor for files or folders that were changed since being staged for deletion, use this: `rot list --changed`.


# Downloads

Currently this project is in **alpha** status. This means even though there is unit tests for the basic features it still needs to be tested with real life use cases. This is why there is not binaries for downloads as of yet.

You can build and run this code like any other [Go](https://golang.org) code:

    go get github.com/jojomi/rot


# Development

If you encounter problems please open an issue.

## Run Tests

    go test github.com/jojomi/rot/cmd
