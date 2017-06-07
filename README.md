# lambda-logger

> Better logging for AWS Lambda and Cloudwatch Logs

# What is this?
Ever gotten frustrated by the capabilities of the UI of *Cloudwatch Logs*? Ever developed a lambda function and got frustrated by the limited UI of cloudwatch logs?

This tool streams logs from cloudwatch logs to GCP Stackdriver Logs to provide superb log viewing UX.
It is simply a lambda function that gets streamed logs from AWS log groups (where all lambda functions send logs to) and sends them to external logging services (such as stackdriver). It is a zero-maintenance solution for getting a better log viewing UX.

# Why?
As an AWS Lambda function developer you are limited in your selection of log services. Right now the only viable solution is Cloudwatch Logs.
Cloudwatch Logs work out of the box for any lambda function which is nice, but the UI is very limited and leaves a lot to be desired.
To name just a few shortcomings of the UI: There are no levels (info, warn, error), you can't simply tail a log forever b/c it gets directed into smaller fragments, you can't view logs for more than a single function at a time (which is typical requirement when developer functions that work together).

Luckily, as much as the UI sucks, there's an easy way out.
It's easy to simply stream logs to a lambda function. This lambda function can do whatever you want with them. In our case we send them to stackdriver.

## Why stackdriver?
Stackdriver logs are the logs of choice for us at Yodas but if you take a quick look at the code you'd notice that it's quite modular so it's easy to add more loggers (log outputs) to replace stackdriver.

# How to use
There are 3 main steps (and we'll follow along with detailed explanation)

1. Create a GCP stackdriver logging account and download its credentials file
2. Deploy the lambda-logger function onto your AWS account.
3. Configure Cludwatch Log Groups to stream the logs into your new lambda function.


## 1. Create a GCP stackdriver logging account and download its credentials file

1. Create GCP account. (if you don't have one yet)
2. Create a GCP project. (if you don't have one yet)
3. Add Stackdriver Logging. (if you didn't add it yet)

![Enable Stackdriver logging](doc/enable-stackdriver.png?raw=true "Enable Stackdriver logging")

4. Create a service account and add the role Logs Writer to this account. Download the JSON credentials file.

![Create a serivce account](doc/create-service-account.png?raw=true "Create a serivce account")

![Create and download credentials file](doc/create-service-account-credentials.png?raw=true "Create and download credentials file")

5. Save the credentials json file and use it to deploy the functino. (see the Makefile for example as well as the next section here)

## 2. Deploy the lambda-logger function onto your AWS account.

Prerequisites: The function uses the following frawework so in oder to successfully deploy you'll have to install them.

1. The [Apex](http://apex.run/) framework
2. And [Go](https://golang.org/) for implementation
3. And [Glide](https://github.com/Masterminds/glide) for 3rd party go vendoring (dependency management)

You'll also need a `.aws/credentials` file or similar mechanism to authenticate to AWS for the sake of deployement, but since you're here, we assume you already have that.

Steps:

1. Create a role for the lambda and use it in `project.json` (replace `arn:aws:iam::032352966958:role/lambda-logger_lambda_function` with your created role's ARN). This role needs to be very simple, just allow running a lambda function and writing to cloudwatch logs (just in case... for monitoring this function's behavior)

The trust relationship is pretty standard for lambda functions:

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

2. Deploy the function (see in the `Makefile` the `deploy` target). You'll have to place the credentials file from stackdriver right next to the `main.go` file in `functions/stackdriver/`. This way it'll get automatically deplyment alongside the function itself.

Upon successful deployment you should see this:

![Lambda deployed](doc/lambda-deployed.png?raw=true "Lambda deployed")


## 3. Configure Cludwatch Log Groups to stream the logs into your new lambda function.
Now go to your cloudwatch logs and condigure each log group to stream its logs to you new lambda function.

Select your new Lambda funcation and select `JSON` as the format

![Setup streaming step 1](doc/cloudwatch-logs-setup1.png?raw=true "Setup streaming step 1")
![Setup streaming step 2](doc/cloudwatch-logs-setup2.png?raw=true "Setup streaming step 2")
![Setup streaming step 3 - user JSON](doc/cloudwatch-logs-setup3.png?raw=true "Setup streaming step 3 - use JSON")


That's how it should look when you're all done:
![Lambdas logs streamed](doc/cloudwatch-logs-streamed.png?raw=true "Lambdas logs streamed")


That's it, now you should start seeing logs in stackdriver!

# Costs
Note about costs and charges:  
AWS Lambda costs for each activation (it's not much, but >0), as well as outgoing network traffic from AWS, incoming network traffic to GCP and stackdriver logging service costs.  
For us at Yodas these costs are still negligible, but you should consider them.

