# go-grpc-api
This is starter code for a user-management microservice written in go. It uses gRPC and protocol buffers internally, while exposing a REST API. Firebase manages authentication, and mySQL is used for persistence.

## Set up dependencies
This is a guide to running the app in a local sandbox. This will let you build on top of this framework. When you're ready to deploy it, follow this guide (TODO). The only prerequiesite is to have go installed.

### Step 1 - Renaming the project
To get started, clone this repo and rename it something suitable for your context. Then, run the rename script with your new project name, often of the pattern `github.com/<username>/<repository_name>`:
```
$ ./rename_project.sh <project_name>
```

Run `$ go mod tidy` to clean up dependencies.

### Step 2 - Set up Firebase auth
Now, to set up authentication, you need to [create a Firebase project](https://console.firebase.google.com/), and [enable authentication](https://firebase.google.com/docs/auth).

This project uses Firebase for user id generation, and auth token generation/verification. Firebase provides free, unlimited authentication :)

In order for your code to authenticate inbound requests, it needs to authenticate itself with Firebase. We do that with an API key, associated with a service account. So, from the Firebase console, you need to navigate to the service accounts page (project settings > Service accounts), then generate a private key. Once you've done so, download it to some file, e.g. ~/.project_key.json.

Finally, store the encoded API key in a special environment variable so the server can authenticate with firebase:
```
$ FIREBASE_JSON=$(base64 ~/.project_key.json)
```

### Step 3 - Set up the database
This project uses mySQL for storage. You need to install:
- The [daemon](https://dev.mysql.com/downloads/mysql/)
- The [shell](https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-install.html)

Once you do, create a local mysql user for your server to identify as:
```
$ mysql -u root -p 
mysql> CREATE USER 'db_user' IDENTIFIED BY 'db_password';
```

Then, set the proper environment variables for your code to authenticate as this user:
```
$ DB_USER=<db_user>
$ DB_PASSWORD=<db_password>
$ DB_IP_ADDR=localhost:3306 # (default port for mySQL)
```

Once you do, start the daeomon (in macOS, you do it in system preferences). Then, initialize a dev database for testing:
```
$ mysql -u root -p < db/up.sql
```
If you want to bring down the database, run `$ mysql -u root -p < db/down.sql` (Note: this deletes all data in the database).

### Step 4 - Build the server
If all goes well, you can use make to compile the go code:
```
$ make server
```

### Thanks
TODO
