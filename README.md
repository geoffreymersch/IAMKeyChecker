# IAMKeyChecker
Look for expired AWS Access Keys of all IAM users within an account or a specific one.

## Usage
### CLI
    export HOURS=10
    go run main.go cli default
        2023/05/25 16:43:56 AWS Profile name provided: default
        2023/05/25 16:43:56 Fetching all the IAM Users.
        2023/05/25 16:43:58 All the IAM Users fetched : 1002
        2023/05/25 16:44:11 List of all the IAM Users using an expired key: [geoffrey]

    go run main.go cli default test-user-generation-970
        2023/05/25 16:45:34 AWS Profile name provided: default
        2023/05/25 16:45:35 User test-user-generation-970 is not using an expired key.

### API
    go run main.go server default
   
### Result :

    curl -s http://127.0.0.1:8080/users/
        {"list_users_expired_key":["geoffrey"],"number_users_expired_key":1}

    curl -s http://127.0.0.1:8080/user/test-user-generation-933
        {"iam_user_name":"test-user-generation-933","is_using_expired_key":false}

    curl -s http://127.0.0.1:8080/user/geoffrey
        {"iam_user_name":"geoffrey","is_using_expired_key":true}

## Docker
If not running inside AWS, you can copy your AWS profile.

    docker build -t iamkeychecker .
    docker run -e HOURS=3 -it -m 64m -v ~/.aws:/root/.aws  iamkeychecker cli default
    docker run -e HOURS=3 -it -m 64m -v ~/.aws:/root/.aws  iamkeychecker cli default iam-user-name

## TODO
- Get rid of AWS static credentials
- Authentication / IP filtering
- Unit tests 
- Integrations tests