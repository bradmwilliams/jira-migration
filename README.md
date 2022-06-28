# jira-migration
To run either of these Jira Clients, you must provision an API Token, for yourself, and pass it into the binary:

## ci-search-jira-client
```
$ ./ci-search-jira-client --jira-endpoint https://issues.redhat.com --jira-bearer-token-file /tmp/api
```

## prow-jira-client
```
$ ./prow-jira-client --jira-endpoint https://issues.redhat.com --jira-bearer-token-file /tmp/api
```
