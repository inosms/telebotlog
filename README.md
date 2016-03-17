# telebotlog
a simple command for redirecting output from any command to your telegram account.

For Example:

`ls /home/magic | telebotlog send homegroup`

This sends the output of the `ls` command to all receivers listed in homegroup.

Each time `telebotlog` reads a `\n` a new message is sent.

# Usage

At first you have to create your own bot and [obtain an auth token from the @BotFather](https://core.telegram.org/bots).

After this register the bot with
`telebotlog register TOKEN`.

Now you have to register users. As Telegram bots can not initiate a conversation with a user one has to create an invitation link that the user presses in order to initiate the conversation.

In order to support easy naming the users are grouped. With this one can redirect the output of a command to a group of users. Users can only be invited into groups. So one has to create a group at first:

`telebotlog group create BOTNAME GROUPNAME`

This creates a new group with name `GROUPNAME` which sends all its inputs via the bot specified by `BOTNAME`.
The bot has to registered beforehand.

Finally the invitation for groups can be created!

`telebotlog group invite GROUPNAME`

This prints a link which has to be pressed in order to accept the invitation.

**NOTE** the command must not be interrupted while the user has not yet accepted the invitation, otherwise the link will become invalid!

After successfully pressing the link the user should be registered for the group. One can check this by listing the groups

`telebotlog group list`

which will print a list of groups and their users.


The configuration is complete now and one can use the command:

`ls /home/magic | telebotlog send GROUPNAME`

## Configuration

`telebotlog` stores all its configurations in `~/.telebotlog.conf`


# TODO
- Better User names than IDs
- Not only send on `\n`, make this configurable
- Send directly to user and not to group?
