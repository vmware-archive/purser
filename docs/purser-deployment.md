# Purser Deployment

In order to deploy the Purser UI and DGraph database service, follow the below listed steps:

1. Switch the current context to point to the desired cluster.

    ``` bash
    kubectl config use-context <context>
    ```

    Read more about configuring and setting the `KUBECONFIG` and kubernetes context [here](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/).

2. If the cluster does not have a valid public IP, set proxy in order to expose the service externally.

    ``` bash
    kubectl proxy
    ```

3. When set, you can simply deploy the Purser UI and Dgraph database service using target `make deploy-purser`.

   _If you wish to however, deploy the database service and the UI service separately, execute the following targets respectively._

   ``` bash
   # deploy Dgraph database
   make kubectl-deploy-purser-db

   # deploy purser UI
   make kubectl-deploy-purser-ui
   ```

4. Once deployed, if proxy was set the UI service can be accessed from [this url](http://127.0.0.1:8001/api/v1/namespaces/default/services/http:purser-ui:4200/proxy/home).

    If public IP was available for your cluster, the UI service should be accessible from path `<External-Public-IP>:<NodePort>`.

    Eg. `http://<minishiftIP>:<NodePort>/home`

5. In order to drop the Dgraph entries from the database, delete the `Persistent Volume` corresponding to the `dgraph datadir`.