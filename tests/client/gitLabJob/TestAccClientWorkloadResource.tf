provider "aembit" {
}

resource "aembit_client_workload" "test" {
    name = "Unit Test 1 - gitLabJob"
    description = "Acceptance Test client workload"
    is_active = false
    identities = [
        {
            type = "gitlabIdTokenNamespacePath"
            value = "namespacePath"
        },
        {
            type = "gitlabIdTokenProjectPath"
            value = "projectPath"
        },
        {
            type = "gitlabIdTokenRefPath"
            value = "refPath"
        },
        {
            type = "gitlabIdTokenSubject"
            value = "subject"
        },
    ]
}

