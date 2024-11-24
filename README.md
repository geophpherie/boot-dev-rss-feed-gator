# boot-dev-rss-feed-gator
RSS Feed aggregator after update

# Dependencies
You'll need `go` and `postgres` installed for this to run.

# Installation
You can't install the program from github as I did not properly name my go.mod file and am not going to fix it. But you can install after cloning this repo by running `go install ./cmd/.`.

# Config
You'll need to set up a config file for the program. It will look for a `.gatorconfig.json` file in your home directory. This config file needs two keys, `db_url` and `current_user_name`. However, the `current_user_name` will be set by the program. The `db_url` should point to your running postgres instance with sslmode=disabled.

# Commands.
The current list of available commands:
* `login`
* `register`
* `reset`
* `users`
* `agg`
* `addfeed`
* `feeds`
* `follow`
* `following`
* `unfollow`
* `browse`
