## Dependencies ##

1. Install [Go](https://golang.org/dl/), [git](https://git-scm.com/downloads),
   and [Docker](https://www.docker.com/) . You may use the official binaries or 
   your usual package manager, whatever you prefer is fine.

1. Verify that they were properly installed.

        go version, should be at least 1.7
    
        git version
    
        docker version

## Compiling the binaries locally ##

1. If GOPATH environment variable isn't set then set it by following [https://github.com/golang/go/wiki/SettingGOPATH](https://github.com/golang/go/wiki/SettingGOPATH)

1. Add GOPATH/bin to your PATH environment variable by running following commands

        export PATH=$PATH:$GOPATH/bin
        	
   Optionally add the above exports to your `.bash_profile` to persist across console sessions.

1. Fetch Purser project

        go get github.com/vmware/purser

1. Run the following command to create a purser plugin binary in 
   `GOPATH/bin` directory

        go install github.com/vmware/purser/cmd/purser_plugin

## Running locally ##

1. Copy the [plugin.yaml](../plugin.yaml) into one of the paths specified under 
   section [Installing kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)

1. Run the following command to check working of purser plugin locally.

        kubectl --kubeconfig=<absolute path to kubeconfig file> plugin purser help