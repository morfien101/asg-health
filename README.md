# asg-health

This simple tool can be used to see if EC2 AutoScaling instances are `InService`.
It can also set the custom health of the instance.

## Help menu

```text
        Use this tool to check if the EC2 instance is 'InService' in it's Autoscaling Group.
        It can also set the instance's custom health attribute to 'Unhealthy' which will cause
        the AutoScaling Group to start the termination process.

        Only a single action can be invoked in a single run.
        It will consume credentials from instance policies or ENV vars.
        There is no provision for manually feeding in credentials and never will be.

  -h    Show the help menu
  -i string
        instance_id for the EC2 instance. If - is passed the instance ID is determined automatically from the metadata if available
  -in-service
        Checks to see if the instance is 'InService' in it's AutoScaling Group
  -set-unhealthy
        Set the instance to unhealthy in it's AutoScaling Group
  -v    Show the version
  -verbose
        Will log success statements as well as errors
```
