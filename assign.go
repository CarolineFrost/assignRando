package main
import "fmt"
import "github.com/google/go-github/github"
import "golang.org/x/oauth2"
// import "math"


func main(){

    user := "carolinefrost"
    userRepo := "go-ml"
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{ AccessToken: "a7e59dff324d25f75348bf5846af54f5e6249580"},
    )
    tc := oauth2.NewClient(oauth2.NoContext, ts)

    client := github.NewClient(tc)
    opt := &github.RepositoryListOptions{}
    repos, _, _ := client.Repositories.List(user, opt)
    var curRepo *github.Repository
    for _,repo := range repos {
        if *repo.Name == userRepo {
            curRepo = repo
        }
    }
    if (curRepo == nil) {return}
    possibleAssignees, _, _ := client.Issues.ListAssignees(user, *curRepo.Name, &github.ListOptions{})
    for _, assignee := range possibleAssignees {
        fmt.Println(*assignee.Login)
    }

    allIssues, _, _ := client.Issues.ListByRepo(user,*curRepo.Name, &github.IssueListByRepoOptions{})

    for _, issue := range allIssues {
        if len(issue.Assignees) == 0 {
            fmt.Println(*issue.HTMLURL)
        }
    }
}
