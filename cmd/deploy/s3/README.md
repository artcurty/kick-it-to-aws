# Deploy to S3

This directory contains the code necessary to upload files to an AWS S3 bucket in parallel.

## Functionality
- Uploads parallel files to a specified S3 bucket.

## Environment Variables

Make sure to set the following environment variables before running the application:

| Environment Variable      | Description                                      | Example                      |
|---------------------------|--------------------------------------------------|------------------------------|
| `AWS_ACCESS_KEY_ID`       | AWS access key ID for accessing S3 and other services. | `AKIAxxxxxxxxxxxxxxxx`       |
| `AWS_SECRET_ACCESS_KEY`   | AWS secret access key for accessing S3 and other services. | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `S3_PATH`                 | Base path in S3 where the files will be stored.  | `my/path/s3`                 |
| `S3_FILES_PATH`           | Local path of the files to be uploaded.          | `/path/to/my/files`          |
| `S3_BUCKET`               | Name of the S3 bucket.                           | `my-bucket`                  |
| `MAX_PARALLEL_UPLOADS`    | Maximum number of parallel uploads allowed.      | `5`                          |

## Usage Example

Here’s how you can use the static site deployment script in a CI pipeline like Drone CI:

```yaml
kind: pipeline
type: docker
name: deploy

steps:
    - name: deploy
      image: golang:1.16
      environment:
        AWS_ACCESS_KEY_ID:
          from_secret: AWS_ACCESS_KEY_ID
        AWS_SECRET_ACCESS_KEY:
          from_secret: AWS_SECRET_ACCESS_KEY
        AWS_REGION: us-west-1
        S3_BUCKET_NAME: my-s3-bucket
        MAX_PARALLEL_UPLOADS: 5
      commands:
        - git clone https://github.com/artcurty/kick-it-to-aws.git
        - cd kick-it-to-aws
        - go build -o kick-it-to-aws cmd/deploy/s3/main.go
        - ./kick-it-to-aws
```

This pipeline builds the Go application and runs the `kick-it-to-aws` script to upload parallel files to S3.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
