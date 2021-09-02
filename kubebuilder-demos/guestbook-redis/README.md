# GuestBook-Redis

1. init project
```bash
go mod init guestbook-redis
kubebuilder init --domain walk1ng.dev
kubebuilder create api --group webapp --version v1 --kind GuestBook
```
2. go on