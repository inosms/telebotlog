# telelog
a simple command for redirecting output from any command via telegram supporting multiple receivers.

Each time `telelog` reads a `\n` a new message is sent.

### Examples

`ls /home/magic | telelog send homegroup`

This sends the output of the `ls` command to all receivers listed in homegroup via telegram.

One could for example use this to start a lengthy process on a remote ssh server and be notified when it is finished.

`createuniverse --superbig && echo "universe ready" | telelog send mydevice`

or simply log the output from a process to your device:

`command | telelog send mydevice`

# Usage

## Building
Get telelog via

`go get github.com/inosms/telelog`

and build it via

`go build github.com/inosms/telelog`

## The Command

At first you have to create your own bot and [obtain an auth token from the @BotFather](https://core.telegram.org/bots).

After this register the bot with

`telelog register TOKEN`.

Now you have to register users. As Telegram bots can not initiate a conversation with a user one has to create an invitation link that the user presses in order to initiate the conversation.

In order to support easy naming the users are grouped. With this one can redirect the output of a command to a group of users. Users can only be invited into groups. So one has to create a group at first:

`telelog group create BOTNAME GROUPNAME`

This creates a new group with name `GROUPNAME` which sends all its inputs via the bot specified by `BOTNAME`.
The bot has to registered beforehand.

Finally the invitation for groups can be created!

`telelog group invite GROUPNAME`

This prints a link which has to be pressed in order to accept the invitation.

**NOTE** the command must not be interrupted while the user has not yet accepted the invitation, otherwise the link will become invalid!

After successfully pressing the link the user should be registered for the group. One can check this by listing the groups

`telelog group list`

which will print a list of groups and their users.


The configuration is complete now and one can use the command:

`ls /home/magic | telelog send GROUPNAME`

## Configuration

`telelog` stores all its configurations in `~/.telelog.conf`


# TODO
- Better User names than IDs
- Not only send on `\n`, make this configurable
- Send directly to user and not to group?
