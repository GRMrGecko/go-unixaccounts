# go-unixaccounts
A simple UNIX account information parser of /etc/passwd and /etc/group for GoLang.

## Install
go get github.com/grmrgecko/go-unixaccounts

## Example
```go
import (
    "fmt"
    UNIXAccounts "github.com/grmrgecko/go-unixaccounts"
)

func main() {
    accounts := UNIXAccounts.NewUNIXAccounts()

    user := accounts.UserWithName("root")
    groups := accounts.UserMemberOf(user)

    var groupNames []string
    for _, group := range groups {
        groupNames = append(groupNames, group.Name)
    }

    fmt.Println("Found groups root is a member of:", groupNames)
}
```
