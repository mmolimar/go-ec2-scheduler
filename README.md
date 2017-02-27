# Go EC2 scheduler

A simple app for starting/stopping EC2 instances in AWS.

## Configuring Credentials

Before using the app, ensure that you've configured credentials. The best
way to configure credentials on a development machine is to use the
`~/.aws/credentials` file, which might look like:
```
[default]
aws_access_key_id = MY-ACCESS-KEY-ID
aws_secret_access_key = MY-SECRET-ACCESS-KEY
```

## Usage

In your project's root dir:
- `go run ec2-manager.go -help` to show usage
- `go run ec2-manager.go -action=start -region=<REGION> -tag=<TAG> <VALUE1> <VALUE2> <VALUE N>` to start instances filtered by those tag and region.
- `go run ec2-manager.go -action=stop -region=<REGION> -tag=<TAG> <VALUE1> <VALUE2> <VALUE N>` to stop instances filtered by those tag and region.

## Installing Kala

    ```
	go get github.com/mmolimar/go-ec2-scheduler
	```
	
## Building
- Run `make` or `make build` to compile the app.
- Run `make install` to compile and install the app.
- Run `make clean` to clean up.

# TODO's

- [ ] Dockerizing the app
- [ ] Include cron scheduler

## Contributing

If you would like to add something to this app, you are welcome to do so!


## License

Released under the MIT license.