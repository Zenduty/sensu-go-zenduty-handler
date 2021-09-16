[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/Zenduty/sensu-go-zenduty-handler)
![goreleaser](https://github.com/Zenduty/sensu-go-zenduty-handler/workflows/goreleaser/badge.svg)
[![Go Test](https://github.com/Zenduty/sensu-go-zenduty-handler/workflows/Go%20Test/badge.svg)](https://github.com/Zenduty/sensu-go-zenduty-handler/actions?query=workflow%3A%22Go+Test%22)
[![goreleaser](https://github.com/Zenduty/sensu-go-zenduty-handler/workflows/goreleaser/badge.svg)](https://github.com/Zenduty/sensu-go-zenduty-handler/actions?query=workflow%3Agoreleaser)

# Sensu Go Zenduty Plugin

## Table of Contents

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
  - [Annotations](#annotations)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The Sensu Zenduty Handler is a [Sensu Handler][6] that can trigger and resolve incident in Zenduty for alerting Operators and entities.

## Usage examples

### Help

```
Usage:
  sensu-go-zenduty-handler [flags]

Flags:
  -d, --debug                    Enable debug mode, which prints JSON object which would be POSTed to the Zenduty webhook instead of actually POSTing it
  -h, --help                     help for sensu-go-zenduty-handler
  -w, --webhook string           The Zenduty Webhook URL, use default from ZENDUTY_WEBHOOK env var
  -a, --withAnnotations string   The Zenduty handler will parse check and entity annotations with these values. Use ZENDUTY_ANNOTATIONS env var with commas, like: documentation,playbook
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add Zenduty/sensu-go-zenduty-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/Zenduty/sensu-go-zenduty-handler].

### Handler definition

```yml
---
api_version: core/v2
type: Handler
metadata:
  name: sensu-go-zenduty-handler
  namespace: default
spec:
  type: pipe
  command: 'sensu-go-zenduty-handler -w "${ZENDUTY_WEBHOOK}"'
  timeout: 0
  filters:
    - is_incident
    - not_silenced
  env_vars:
    - ZENDUTY_WEBHOOK={{YOUR_ZENDUTY_WEBHOOK}}
  runtime_assets:
    - sensu-go-zenduty-handler
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-go-zenduty-handler repository:

```
go build
```

[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/handler-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/handler-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/handlers/
[7]: https://github.com/sensu-community/handler-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
