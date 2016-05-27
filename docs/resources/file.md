## file

The `file` resource is used for managing files and directories.

## Embeds

* [base](base.md)

## Parameters

### path

Path to the target file to manage.

* Type: string
* Required: no
* Default: defaults to the resource title

### mode

Permission bits to set on the file.

* Type: integer
* Required: no
* Default: 0644

### owner

Owner to set on the target file.

* Type: string
* Required: no
* Default: the username of the user currently executing the resource

### group

Group to set on the target file.

* Type: string
* Required: no
* Default: the primary group of the user currently executing the resource

### source

Source file to use for setting the file content.

The source file is expected to be found in the site repo of the minion.

* Type: string
* Required: no
* Default: none

### filetype

The file type of the managed file.

* Type: string
* Required: no
* Default: regular
* Values: `regular`, `directory`

### recursive

A flag used to indicate that the target directory content should be
managed recursively.

Applies only to resources, which have the `filetype`
parameter set to `directory`.

* Type: bool
* Required: no
* Default: false

### purge

A flag used to indicate that extra files present in the target
directory, but not present in the source should be purged.

* Type: bool
* Required: no
* Default: false

## Example usage

The following resource ensures that `/tmp/foo-dir` is a directory.

```hcl
file "/tmp/foo-dir" {
  state = "present"
  filetype = "directory"
}
```

The following resource ensures that `/tmp/foo-file` is a regular
file and it's content are set from the provided source file in the
site repo.

```hcl
file "/tmp/foo-file" {
  state = "present"
  mode = 0644
  source = "data/foo-file"
}
```
