## Subscription Service Application.
This is built  using goroutines to 
- send email concurrently ,
- generate Invoice and a manual concurrently
- Shuts down the application by cleaning undone tasks and closing the channels 

The Main focus on building this simple app emphasizes the use of concurrency 
You can run it locally on your machine , But you must have the following configurations and setup on your machine 
- [Postgres]() for relational db management
- [MailHog]() for handling and sending mails locally 
- [Redis]() for session caching