# authen-amazon-ecr

## Usage

```yml
jobs:
  authen-ecr:
    steps:
      - uses: aws-actions/configure-aws-credentials@v1
      - id: authen-amazon-ecr
        uses: aereal/authen-amazon-ecr@v1
      - run: docker login --username $_USERNAME --password $_PASSWORD $_SERVER
        env:
          _USERNAME: ${{ steps.authen-amazon-ecr.outputs.username }}
          _PASSWORD: ${{ steps.authen-amazon-ecr.outputs.password }}
          _SERVER: ${{ steps.authen-amazon-ecr.outputs.server }}
```

### Use with service container credentials

```yml
jobs:
  authen-ecr:
    outputs:
      username: ${{ steps.authen-amazon-ecr.outputs.username }}
      password: ${{ steps.authen-amazon-ecr.outputs.password }}
    steps:
      - uses: aws-actions/configure-aws-credentials@v1
      - id: authen-amazon-ecr
        uses: aereal/authen-amazon-ecr@v1
  test:
    needs:
      - authen-ecr
    services:
      mysql:
        image: $YOUR_PRIVATE_ECR_REPO_SERVER/$PRIVATE_REPO:latest
        credentials:
          username: ${{ needs.authen-ecr.outputs.username }}
          password: ${{ needs.authen-ecr.outputs.password }}
```
