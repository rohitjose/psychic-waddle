# psychic-waddle
Codebase for production ready serverless course

## Initial Setup
### Serverless Framework
 - Setup npm to install serverless framework
 - Install serverless framework
 - Setup go-serverless boilerplate code
  ```
  serverless create -t aws-go-mod -p myservice
  ```
### AWS Configuration for Serverless Framework
 - Setup new IAM user with Administrator access. This is done during the development stage. After an initial deployment the Cloudformation template can be analysed to define a Policy document of a more limited permissions.
 - Next, use the API key and secret to configure serverless framework access. If your aws CLI is already set, you can choose to add the ```--overwrite``` flag to overwrite the existing configuration.
```
serverless config credentials --provider aws --key <key> --secret <secret> --overwrite
```
### Go Dep Installation
The template here uses Go Dep for package management. The project needs to exist in the path ```$GOPATH/src/github.com```. You also need to ensure that the go environment variable ```GO111MODULE=auto```
```
brew install dep
brew upgrade dep
go env -w GO111MODULE=auto
```
### Build and deploy
Run the following commands
```
make
sls deploy
```
Output snippet:
```
go mod download github.com/aws/aws-lambda-go
➜  myservice git:(new-template) ✗ make
chmod u+x gomod.sh
./gomod.sh
export GO111MODULE=on
env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
➜ myservice git:(init-setup) ✗ sls deploy
Serverless: Packaging service...
Serverless: Excluding development dependencies...
Serverless: Creating Stack...
Serverless: Checking Stack create progress...
........
Serverless: Stack create finished...
Serverless: Uploading CloudFormation file to S3...
Serverless: Uploading artifacts...
Serverless: Uploading service myservice.zip file to S3 (4.91 MB)...
Serverless: Validating template...
Serverless: Updating Stack...
Serverless: Checking Stack update progress...
................................................
```
### Testing 
```
curl -X GET https://cm7di14fm8.execute-api.ap-southeast-2.amazonaws.com/hello 
curl -X GET https://cm7di14fm8.execute-api.ap-southeast-2.amazonaws.com/world
```

### Cleanup
Teardown the stack with the following command:
```
➜  myservice git:(init-setup) ✗ sls remove
Serverless: Getting all objects in S3 bucket...
Serverless: Removing objects in S3 bucket...
Serverless: Removing Stack...
Serverless: Checking Stack delete progress...
.................
Serverless: Stack delete finished...

Serverless: Stack delete finished...
```
