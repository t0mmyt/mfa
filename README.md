# MFA (TOTP) Command Line Tool

## Installation (Nix/Mac)

    go build -o ~/bin/mfa cmd/mfa/main.go

## Examples
    # Getting an MFA
    $ mfa add github 'otpauth://totp/GitHub:..?secret=....'
    $ mfa do github
    123456
    
    # Syncing between devices
    $ mfa export | ssh mfa import

```
$ mfa help 
usage: mfa [<flags>] <command> [<args> ...]

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  do <name>
    do an mfa

  add <name> <url>
    add an mfa

  list
    show all keys

  export
    export all keys as JSON

  import
    import all keys as JSON
```