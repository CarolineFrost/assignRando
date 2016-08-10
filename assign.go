package main
import "fmt"
import "github.com/google/go-github/github"
import "golang.org/x/oauth2"
import "math/rand"


func main(){

    org := "sourcegraph-beta"
    userRepo := "sourcegraph-desktop"
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{ AccessToken: "68ed6eb0012b7c189133900507ac2b5967a147c9"},
    )
    tc := oauth2.NewClient(oauth2.NoContext, ts)
    client := github.NewClient(tc)

    if ensureRepoExists(client, userRepo) != nil {return}

    possibleAssignees, _, _ := client.Issues.ListAssignees("sourcegraph", "sourcegraph-desktop", &github.ListOptions{})
    if len(possibleAssignees) == 0 {return}
    randGen := rand.New(rand.NewSource(1))

    for true {
        allIssues, _, _ := client.Issues.ListByRepo(org, userRepo, &github.IssueListByRepoOptions{})
        unassignedIssues := make(map[*github.Issue]bool)
        assignedIssues := make(map[*github.Issue]bool)

        for _, issue := range allIssues {
            if len(issue.Assignees) == 0 {
                unassignedIssues[issue] = true
            } else {
                assignedIssues[issue] = true
            }
        }
        for issue := range unassignedIssues {
            teamSize := len(possibleAssignees)
            randomIndex := randGen.Intn(teamSize)
            assignRandomly(org, userRepo, client, issue, possibleAssignees, randomIndex)
        }
    }   
}

func ensureRepoExists(client *github.Client, userRepo string) (err error) {
    repos, _, _ := client.Repositories.ListByOrg("sourcegraph-beta", &github.RepositoryListByOrgOptions{})
    if (len(repos) == 0) {return}

    var curRepo *github.Repository
    for _,repo := range repos {
        if *repo.Name == userRepo {
            curRepo = repo
        }
    }
    if (curRepo == nil) {return err}
    return nil
}

func assignRandomly(org string, repo string, client *github.Client, issue *github.Issue, possibleAssignees []*github.User, randomIndex int) {
    assignedDev := possibleAssignees[randomIndex]
    iss, _, _ := client.Issues.AddAssignees(org, repo, *issue.Number, []string{*assignedDev.Login})
    if len(iss.Assignees) == 0 {
        fmt.Println(fmt.Sprintf("Unable to assign %s to issue number %d", *assignedDev.Login, *iss.Number))
    } else {
        fmt.Println(fmt.Sprintf("Assigned %s to issue number %d", *assignedDev.Login, *issue.Number))
    }
}