# mailchimp.v3
MailChimp API v3 lib in Go

## Create a client
A client is used when communicating with Mailchimp.
```
client := mailchimp.NewClient(MailchimpToken)
```

### Create a member
```
cmd := &mailchimp.CreateMember{
    EmailAddress: "demo@example.net",
    EmailType:    mailchimp.HTML,
    Status:       mailchimp.Subscribed,
    MergeFields: map[string]string{
        "FNAME": contact.Name,
        "LNAME": contact.Surname,
    },
    Vip: false,
}
member, err := client.CreateMember(cmd, listID)
if err != nil {
    // ...
}
```

### Create a list
```
list, err := s.client.NewList(&mailchimp.CreateList{
    Name:                "My list",
    Contact:             "demo@example.net",
    PermissionReminder:  "You subscribed to me",
    UseArchiveBar:       false,
    CampaignDefaults:    campaignDefaults,
    NotifyOnSubscribe:   "new_subscription@example.net",
    NotifyOnUnsubscribe: "lost_subscription@example.net",
    EmailTypeOption:     false,
    Visibility:          mailchimp.ListVisibilityPublic,
})

if err != nil {
    return nil, err
}
```

### Query stuff
```
list, err := client.GetList(s.RemoteID)
if err != nil {
    return nil, err
}

// mailchimp.Parameters is an alias for map[string]interface{}
params := &mailchimp.Parameters{"status": "subscribed"}

rlist, err := client.GetMembers(listID, params) // params is optional
if err != nil {
    return nil, err
}
```
