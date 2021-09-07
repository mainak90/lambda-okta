package okta_actions

import (
	"context"
	"fmt"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func RemoveUserFromGroup(ctx context.Context, url string, token string, user string, group string) error {
	ctx, client, err := okta.NewClient(ctx, okta.WithOrgUrl(url), okta.WithToken(token), okta.WithRequestTimeout(45), okta.WithRateLimitMaxRetries(3))
	if err != nil {
		return err
	}
	_, err = client.Group.RemoveUserFromGroup(ctx, group, user)
	if err != nil {
		return err
	}
	fmt.Printf(user, group, "User %s added to group %s \n")
	return nil
}

