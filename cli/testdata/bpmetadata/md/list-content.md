# h1 doc title
some content doc title for h1
some content doc title for h1

## h2 sub heading
some content sub heading for h2

## Documentation
- [document-01](http://google.com/doc-01)
- [document-02](http://google.com/doc-02)
- [document-03](http://google.com/doc-03)
- [document-04](http://google.com/doc-04)

### h3 sub sub heading

## Horizontal Rules

## Diagrams
- text-document-01
- text-document-02

## Description:
### tagline
Opinionated GCP project creation

### detailed
This blueprint allows you to create opinionated Google Cloud Platform projects.
It creates projects and configures aspects like Shared VPC connectivity, IAM access, Service Accounts, and API enablement to follow best practices.

### preDeploy
To deploy this blueprint you must have an active billing account and billing permissions.

### Architecture
1. User requests are sent to the front end, which is deployed on two Cloud Run services as containers to support high scalability applications.
2. The request then lands on the middle tier, which is the API layer that provides access to the backend. This is also deployed on Cloud Run for scalability and ease of deployment in multiple languages. This middleware is a Golang based API.

### Architecture
![Test Image](https://i.redd.it/w3kr4m2fi3111.png)
1. Step 1
2. Step 2
3. Step 3

### ArchitectureNotValid
1. Step 1
2. Step 2
3. Step 3
