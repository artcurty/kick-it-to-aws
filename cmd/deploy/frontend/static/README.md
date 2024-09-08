
# Static Site Deployment (S3 + CloudFront)

This section of the Kick-It-To-AWS project provides a script that handles the deployment of static websites to AWS S3 and invalidates the CloudFront cache.

## Functionality
- Uploads static files to a specified S3 bucket.
- Invalidates CloudFront cache to reflect changes immediately.

## Environment Variables

The following environment variables are required to configure the static site deployment:

| Variable Name               | Description                                           | Example                                      |
|-----------------------------|-------------------------------------------------------|----------------------------------------------|
| `AWS_ACCESS_KEY_ID`          | AWS access key ID for accessing S3 and other services.| `AKIAxxxxxxxxxxxxxxxx`                       |
| `AWS_SECRET_ACCESS_KEY`      | AWS secret access key for accessing S3 and other services.| `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`|
| `AWS_REGION`                 | AWS region where your resources are hosted.           | `us-west-1`                                  |
| `S3_BUCKET_NAME`             | The name of the S3 bucket where files are uploaded.   | `my-s3-bucket`                               |
| `CLOUDFRONT_DISTRIBUTION_ID` | The CloudFront distribution ID for cache invalidation.| `EDFDVBD6EXAMPLE`                            |
| `LOG_LEVEL`                  | Log level for application output.                     | `info`, `debug`, `error`                     |

## Pipeline using
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
        CLOUDFRONT_DISTRIBUTION_ID: E1234ABCDEF
      commands:
        - git clone https://github.com/artcurty/kick-it-to-aws.git
        - cd kick-it-to-aws
        - go build -o kick-it-to-aws cmd/deploy/frontend/static/main.go
        - ./kick-it-to-aws
```

This pipeline builds the Go application and runs the `kick-it-to-aws` script to upload static files to S3 and invalidate the CloudFront cache.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
