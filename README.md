# UniBot
Uni discord bot to do things

*Source code coming soon because I want to fix some bugs that are going on*

## About
Just another bot to do entertaining things to add some spice to your server
#####
She can be invited with this [link](https://discordapp.com/oauth2/authorize?client_id=462421724659580950&scope=bot&permissions=535948390)

### What can uni do?
If enabled she will have the ability to request posts from a subreddit, grab an image from the [inspirobot](http://inspirobot.me/), search for random images off derpibooru with the given tags and even be able to parse it's own lua code for each channel

## How do I enable these modules?
It's simple first you have to be the server owner or have the set admin role
#####
To enable the reddit search module: 
```
hey uni enable module reddit
```
To enable the inspirobot module: 
```
hey uni enable module inspire
```
To enable the derpibooru search module: 
```
hey uni enable module derpi
```
**(Her lua parsing comes later, still WIP)**

## What if I want to disable these?
reddit:
```
hey uni disable module reddit
```
inspirobot:
```
hey uni disable module inspire
```
derpibooru:
```
hey uni disable module derpi
```

**(Her lua parsing comes later, still WIP)**
#####

## Usage of the reddit module
##### To grab any current posts do 
```
hey uni find a post <in/on/from/within> <r/, /r/, or just the name of the subreddit>
```
Example:
```
hey uni find a post in r/yesyesyesno
```

##### For top posts:

```
hey uni find a top post <in/on/from/within> <r/, /r/, or just the name of the subreddit>
```
Example:
```
hey uni find a top post in r/aww
```

##### For new posts:

```
hey uni find a new post <in/on/from/within> <r/, /r/, or just the name of the subreddit>
```
Example:
```
hey uni find a new post in r/FloridaMan
```
###### Notes: If you want to have uni browse NSFW images set the channel's tag to be NSFW
## Usage of the inspiro bot module
```
hey uni inspire me
```
that's about it

#####
## Usage of the derpibooru module
##### To search for an image on derpibooru do
```
hey uni search on <derpi/derpibooru> <tags provided here>
```
#####


Example:
```
hey uni search on derpi artist:rodrigues404
```
Or
```
hey uni search on derpibooru first_seen_at.gt:3 days ago AND score.gte:77
```
#####
#####
##### To see the results of an image on derpibooru do
```
hey uni derpi image <ID of image>
```
Example:
```
hey uni derpi image 1761475
```

###### Notes: Uni will set a channel's filter to be "everything" if the channel is marked "NSFW" and has no set filter

## Admin related things
To be able to set roles as an "admin role" you must do
```
hey uni set admin role <role name or ping here>
```
#####
Someone being a little annoying or just anything that gets him/her banned?
#####
Well uni is able to use the ban hammer ability in case you want uni to do so
```
hey uni <ban/perish> <name or ping here>
```
##### Note: if she returns more than one user she will not ban all of the users listed
To be able to clear messages from chat you do
```
hey uni clear <number here>
```

##### If the derpibooru module is enabled and you want to swap out filters you'll do
```
hey uni set derpi filter <filter ID>
```

## Lua things
Uni also has the ability to parse lua code for every message that's posted in chat
For whatever reason you want to use it for like deleting messages if they contained something it's an available option
##### Notes: All commands related to her lua module will need admin role to modify
##### And her lua package was not mine, it's been slightly modified to have people not modify important variables [here's the original](https://github.com/yuin/gopher-lua)
#### Enabling
```
hey uni enable lua
```
##### Editing the lua script
It's a little awkward but you kinda need to post all of her lua code inside a message
But one way to do it is
```
hey uni rewrite lua ```lua
print("test")```
```
And in case you want to review your lua code you'll do
```
hey uni view lua
```