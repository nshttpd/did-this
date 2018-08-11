#### did-this

n.b. .... currently under development and a work in progress

In remote teams it's helpful to keep track of what you are working on or have completed during the day for
stand ups or status reports. Sometimes these are "online only" via Slack or some other mechanism. This is a CLI
that will help you keep track of what you have done to report on the next day.

Default location to store the DB and config file is ~/.did-this.

Current usage is something like :

`did-this add JIRA-667 rebuilt Kubernetes cluster`

the above will add that completed task into the DB under the current date. There is no going back and padding
previous days, so make sure you keep track of what's going on.

When it's time to report in you can list the tasks you completed :

`did-this list`

This will list all the completed tasks that you saved from the previous day. If you want to remind yourself
of what you've done today you can :

`did-this list today`

or on Monday when you need to list out what you did on Friday :

`did-this list 2018-08-10`

