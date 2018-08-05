# hue-cli

A minimalistic commanline interface tool for managing Philips Hue products.
This README is a work in progress, and is not complete. The `hue-cli` program
supports many different commands and options, use `hue-cli --help` and `hue-cli
<command> --help` to get more details.

`hue-cli` builds on [GoHue project](https://github.com/Collinux/gohue) with
additional features that have not been included yet. The new features and
functionalities are gathered in [a special `new-features`
branch](https://github.com/nixpanic/gohue).


## Find bridges

`hue-cli discover`

Use the [N-UPnP mechanism](https://discovery.meethue.com/) to find bridges in
the local network. This prints information about the bridges that have been
reported by the Hue portal.


## Create a new user

`hue-cli [--bridge=<ip-address>] create-user [--device=<name>] [<new-config,yaml>]`

Creates a new user on the bridge. It is required to press the "link button" and
create the new user within 30 seconds.

The `<new-config.yaml>` is the optional filename where the details of the user
and bridge will be stored, it is written to the console if omitted.
