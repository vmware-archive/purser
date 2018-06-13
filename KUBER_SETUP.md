# Kuber Setup #

## Dependencies ##

1. Install [Go](https://golang.org/dl/), [git](https://git-scm.com/downloads),
   and [Docker](https://www.docker.com/) . You may use the official binaries or 
   your usual package manager, whatever you prefer is fine.

1. Verify that they were properly installed.

        `go version`, should be at least 1.7
    
        `git version`
    
        `docker version`

## Compiling the binaries locally ##

1. Create the directory `kuber-plugin` in your home directory for `kuber-plugin` 
   development.
   
1. Inside `kuber-plugin` directory run the following to get the `kuber-plugin`
    code.

        `git clone git@gitlab.eng.vmware.com:kuber/kuber-plugin.git .`

1. Run the following commands to set the go environment variables.

        `export GOPATH=$HOME/kuber-plugin`
    
        `export GOBIN=$HOME/kuber-plugin/bin`
    
        `export PATH=$PATH:GOBIN`
	
	Optionally add the above exports to your `.bash_profile` to persist across 
	console sessions.

1. Run the following command to create a kuber plugin binary in 
   `kuber-plugin/bin` directory

        `go install kuber`


## Running locally ##

1. Copy the [plugin.yaml](plugin.yaml) into one of the paths specified under 
   section [Installing kubectl plugins]
   (https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)

1. Run the following command to check working of kuber plugin locally.

        `kubectl --kubeconfig=<absolute path to kubeconfig file> plugin kuber help`