parameter-store-exec
====================

Run a program with [AWS Systems Manager Parameter Store][ssmps] parameters injected as environment variables. Perfect for a Docker `ENTRYPOINT` to load secret and plain configuration in production AWS.

A parameter store path is read from the `PARAMETER_STORE_EXEC_PATH` environment variable, which is easy to inject at runtime with `docker run --env ...`. All parameters stored under that path are loaded into environment, with their keys converted to `ENV_FRIENDLY` names. For path `/one/two` the key `/one/two/three/four` becomes `THREE_FOUR`.

```sh
# which AWS region to use; sometimes pre-set e.g. under `aws-vault`.
export AWS_REGION=us-east-1

# Parameter Store path to load into environment.
export PARAMETER_STORE_EXEC_PATH=/foo/bar

# Simle command; just put `parameter-store-exec` in front of your command:
parameter-store-exec <command [args ...]>
```


Prior art
---------

[cultureamp/s3dotenv][s3dotenv] is a similar tool, but backs onto an S3 bucket containing a [dotenv][godotenv] file. It may still be useful, but [Parameter Store][ssmps] is more suitable for storing production runtime configuration, with per-key encryption, versioning and audit.

[segment/chamber][chamber] is a larger and more opinionated way to manage Parameter Store configuration. It's more restrictive regarding hierarchical path namespacing, a feature added to Parameter Store after chamber was first written.  It also needs a service name (parameter namespace) provided on the command line for `chamber exec`, making it harder to parameterize via environment at runtime. Segment did [an excellent write-up][segment-blog] on using Parameter Store for secure configuration on AWS, and particularly ECS.

[Droplr/aws-env][aws-env] also backs onto Parameter Store, but only outputs env vars ready for shell eval, making it less suitable (but workable) as a Dockerfile `ENTRYPOINT`.

For a more recent (and comprehensive) review of alternatives, see the table in the [COMPARISON](COMPARISON.md) file.


## Demo

Store some parameters:

```sh
aws ssm put-parameter --name /acme/staging/hello --value world --type String
aws ssm put-parameter --name /acme/staging/xyz/username --value acme-staging --type String
aws ssm put-parameter --name /acme/staging/xyz/password --value hunter2 --type SecureString
```

Run a command with environment injected from Parameter Store:

```sh
export AWS_REGION=us-east-1
export PARAMETER_STORE_EXEC_PATH=/acme/staging

parameter-store-exec env
# (stderr) 2017/12/06 15:18:12 /acme/staging/hello => HELLO
# (stderr) 2017/12/06 15:18:12 /acme/staging/xyz/username => XYZ_USERNAME
# (stderr) 2017/12/06 15:18:12 /acme/staging/xyz/password => XYZ_PASSWORD
# ...
# HELLO=world
# XYZ_USERNAME=acme-staging
# XYZ_PASSWORD=hunter2
# ...
```

Those log messages are on `stderr`; easy to hide:

```
parameter-store-exec env 2>/dev/null
# ...
# HELLO=world
# XYZ_USERNAME=acme-staging
# XYZ_PASSWORD=hunter2
# ...
```

In a Dockerfile:

```
ENTRYPOINT ["parameter-store-exec"]
```

Then you can inject the desired Parameter Store path at runtime:

```sh
docker run --env PARAMETER_STORE_EXEC_PATH=/acme/staging $image env
# (stderr) 2017/12/06 15:18:12 /acme/staging/hello => HELLO
# (stderr) 2017/12/06 15:18:12 /acme/staging/xyz/username => XYZ_USERNAME
# (stderr) 2017/12/06 15:18:12 /acme/staging/xyz/password => XYZ_PASSWORD
# ...
# HELLO=world
# XYZ_USERNAME=acme-staging
# XYZ_PASSWORD=hunter2
# ...
```

Or use that same Docker image in a non-AWS environment; just don't pass `PARAMETER_STORE_EXEC_PATH`:

```
docker run --env HELLO=local $image env
# (stderr) 2017/12/06 15:21:45 PARAMETER_STORE_EXEC_PATH not set
# HELLO=local
```


[aws-env]: https://github.com/Droplr/aws-env
[aws-vault]: https://github.com/99designs/aws-vault
[chamber]: https://github.com/segmentio/chamber
[godotenv]: https://github.com/joho/godotenv
[s3dotenv]: https://github.com/cultureamp/s3dotenv
[segment-blog]: https://segment.com/blog/the-right-way-to-manage-secrets/
[ssmps]: http://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html
