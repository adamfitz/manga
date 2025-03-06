# manga

To connect to remote postgresql db server create a config file called `manga.config` in `~/.config` directory as follows:

```
$ cat ~/.config/manga.config 
{
	"db_server": "db server name or IP",
	"db_port": "5432",
	"db_user": "db username",
	"db_user_pass": "db users password",
	"db_name": "your database name",
}
```