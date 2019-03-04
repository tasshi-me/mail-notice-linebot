# mail-notice-linebot
メールが来たのを教えるLINEBOT

# SETUP(Heroku)
```bash
#Heroku cli install
$ brew tap heroku/brew && brew install heroku
$ heroku login
$ heroku container:login

# Create heroku app
$ heroku apps:create <Your Heroku App name>
$ heroku stack:set container

# Set local varibales
$ heroku config:set LINE_CHANNEL_ID=<Your LINEBOT Channel ID>
$ heroku config:set LINE_CHANNEL_SECRET=<Your LINEBOT Channel Secret>
$ heroku config:set LINE_ACCESS_TOKEN=<Your LINEBOT Access Token>

# Push & Release Container
$ heroku container:push <Your image name>
$ heroku container:release <Your image name>

# Open your heroku app
$ heroku open

```

