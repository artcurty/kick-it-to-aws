
# Kick-It-To-AWS

Kick-It-To-AWS is a Go-based project aimed at centralizing various deployment scripts for AWS. 
The project is designed to provide a unified interface for different types of AWS deployments, 
making it easier to automate and manage infrastructure changes across multiple services.

## Functionality Summary

1. [Static Site Deployment](cmd/deploy/frontend/static/README.md)
   Deploy static websites to AWS S3 and invalidate CloudFront cache. This functionality is managed by 
   a script located in the `cmd/deploy/frontend/static` folder. 
   It handles the upload of files to S3 and invalidates the CloudFront cache to ensure changes are reflected immediately.


## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
