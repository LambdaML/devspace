---
title: Command - devspace update dependencies
sidebar_label: dependencies
id: version-v4.0.1-devspace_update_dependencies
original_id: devspace_update_dependencies
---


Updates the git repositories of the dependencies defined in the devspace.yaml

## Synopsis


```
devspace update dependencies [flags]
```

```
#######################################################
############ devspace update dependencies #############
#######################################################
Updates the git repositories of the dependencies defined
in the devspace.yaml
#######################################################
```
## Options

```
      --allow-cyclic   When enabled allows cyclic dependencies
  -h, --help           help for dependencies
```

### Options inherited from parent commands

```
      --debug                 Prints the stack trace if an error occurs
      --kube-context string   The kubernetes context to use
  -n, --namespace string      The kubernetes namespace to use
      --no-warn               If true does not show any warning when deploying into a different namespace or kube-context than before
  -p, --profile string        The devspace profile to use (if there is any)
      --silent                Run in silent mode and prevents any devspace log output except panics & fatals
  -s, --switch-context        Switches and uses the last kube context and namespace that was used to deploy the DevSpace project
      --var strings           Variables to override during execution (e.g. --var=MYVAR=MYVALUE)
```

## See Also

* [devspace update](../../cli/commands/devspace_update)	 - Updates the current config
