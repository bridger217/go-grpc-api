# go-grpc-api
This is starter code for a user-management microservice written in go. It uses gRPC and protocol buffers internally, while exposing a REST API. Firebase manages authentication, and mySQL is used for persistence.

## Get Started

### Step 1 - Renaming the project
To get started, fork this repo and rename it something suitable for your context. Then, run the rename script with your new project name, often of the pattern `github.com/<username>/<repository_name>`:
```
$ ./rename_project.sh <project_name>
```

Run `go mod tidy` to clean up dependencies.

