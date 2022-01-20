# azctx: Azure context

`azctx` is a command line tool whose aim is to ease the use of the
official [az cli](https://docs.microsoft.com/en-us/cli/azure/). It is heavily inspired
by [kubectx](https://github.com/ahmetb/kubectx).

It provides an easier and more intuitive way to set your current default azure subscription/account and resource group.

## Requirements

- [az cli](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
- [fzf](https://github.com/junegunn/fzf)

## Usage

- Run `azctx`
- Select a Subscription/Account, (Arrow keys, Search, Enter, Double click)
- Select resource group, or quit (CTRL + C) if you don't want to/clear the default resource group