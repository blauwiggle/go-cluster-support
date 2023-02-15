# go-cluster-support

During my first week in cluster support I noticed that the constant change of the kubeconfig with our
specific Active Directory configuration is time consuming and that sometimes mistakes happen.
Therefore this tool does the work for me to connect to the cluster.

## add .env file

Add a .env file to the root of the project with the following content:

```env
# Define tenant ID
tenantId=<tenant id>

# Define subscription names for dev and prod
devSubscription="<DEV name>"
prodSubscription="<PROD name>"
```

First you need to login to azure with the following command for each subscription:
```bash
az login
```

To get the tenant ID you can use the following command:
```bash
az account show --query 'tenantId' -o tsv
```

To get the subscription names you can use the following command:
```bash
az account show --query 'name' -o tsv
```


## build project

```bash
go build -o bin/go-cluster-support
```

## run project

```bash
./bin/go-cluster-support
```

## install tools

```bash
./bin/go-cluster-support --tools
```

## add alias

Add this to your config file (e.g. ~/.bashrc or ~/.zshrc) depending on your shell.

```bash
alias cats=<path-to-binary>/go-cluster-support
```
