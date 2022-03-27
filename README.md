# go-grpc-api
This is starter code for a user-management microservice written in go. It uses gRPC and protocol buffers internally, while exposing a REST API. Firebase manages authentication, and mySQL is used for persistence.

## Set up dependencies
This is a guide to running the app in a local sandbox. This will let you build on top of this framework. When you're ready to deploy it, follow this guide (TODO). The only prerequiesite is to have go installed.

Note that configuration that contains sensative data is stored in environment variables, so they can more easily be used as secrets in a k8s cluster.

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
$ export FIREBASE_JSON=$(base64 ~/.project_key.json)
```

### Step 3 - Set up the database
This project uses mySQL for storage. You need to install:
- The [daemon](https://dev.mysql.com/downloads/mysql/)
- The [shell](https://dev.mysql.com/doc/mysql-shell/8.0/en/mysql-shell-install.html)

Once you do, create a local mysql user for your server and grant it root privledges:
```
$ mysql -u root -p 
mysql> CREATE USER 'db_user' IDENTIFIED BY 'db_password';
mysql> GRANT ALL PRIVILEGES ON *.* TO 'db_user'@'%';
```

Then, set the proper environment variables for your code to authenticate as this user:
```
$ export DB_USER=<db_user>
$ export DB_PASSWORD=<db_password>
$ export DB_IP_ADDR=localhost:3306 # (default port for mySQL)
$ export DB_NAME=dev # we are using the dev instance
```

Once you do, start the daemon (in macOS, you do it in system preferences). Then, initialize a dev database for testing:
```
$ mysql -u root -p < db/up.sql
```
If you want to bring down the database, run `$ mysql -u root -p < db/down.sql` (Note: this deletes all data in the database).

### Step 4 - Build the server
If all goes well, you can use make to compile the go code:
```
$ make server
```

## Testing the service
Now the all the dependencies are set up, we will test the CreateUser route of the API.

### Step 1 - Create Firebase user
In the Firebase console, navigate to (Authentication > Users > Add user). From there, choose a fake email/password and create the user.

### Step 2 - Generate a Firebase auth token
In this stage, we will use a test UI to authenticate our client with Firebase so that we can communicate with the app. To do so, head over to IllDepence@'s [token.html](https://gist.github.com/IllDepence/7c201287af52bd1f78bed65ec7737e84), and download the file. You need to change two fields in the html:
1. config.projectId -- Found in the Firebase console (Project settings > General)
2. config.apiKey -- Found in the GCP console. Navigate through Firebase console (Project settings > Service accounts > All service accounts) which should take you to GCP. Then, in GCP go to "Credentials" and you will find your default browser key.

Now that your test frontend is set up, open the html in a browser and sign in with the user credentials from step 1. Once you've signed in, click "show ID token" and copy the token into an environment variable:
```
$ export TOKEN=<token_from_ui>
```

### Step 3 - Run the thing
Now, we're ready to run the thing. Make sure you have all necessary env variables set:
- FIREBASE_JSON
- DB_USER
- DB_PASSWORD
- DB_IP_ADDR
- DB_NAME
- TOKEN

Go ahead and run the server. By default it will bind to port 8080:
```
$ ./bin/server
```

Next, send an authenticated curl request that will register our Firebase user with the service:
```
$ curl -X POST localhost:8080/v1/users \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"username": "test", "firstName": "test", "lastName": "user"}'
```

You should receieve a response like:
```
{"id":"3PS5ssaJYKaltvh3PuvqFpoBJgU2","username":"test","firstName":"test", "lastName":"user"}
```

To confirm the user is persisted, you can run:
```
$ mysql -u root -p
mysql> SELECT * FROM dev.Users;
```

Now that the scaffolding is working, start developing your app!

### Feedback
Bugs, changes, or improvements? Feel free to reach out and make a pull request.

### Thanks
TODO
