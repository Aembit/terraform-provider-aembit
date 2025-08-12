provider "aembit" {
}

resource "aembit_client_workload" "test" {
    name = "Unit Test 1 - awsLambdaArn"
    description = "Acceptance Test client workload"
    is_active = false
    identities = [
        {
            type = "awsLambdaArn"
            value = "arn:aws:lambda:us-east-1:880961858887:function:helloworld"
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

