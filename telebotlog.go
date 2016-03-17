package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("telebotlog", "Redirect output via telegram bot. Use with pipes.\n\n"+
		"Example usage:\n\necho \"computer is exploding\" | telebotlog send mygroup")
	verbose = app.Flag("verbose", "Enables Verbose Output").Short('v').Default("false").Bool()

	registerBotCmd   = app.Command("register", "Register a bot from which to send.")
	registerBotToken = registerBotCmd.Arg("token", "Token of the bot.").Required().String()

	groupCmd = app.Command("group", "Configure groups.")

	createLogGroupCmd  = groupCmd.Command("create", "Create a new group of receivers from a bot.")
	createLogGroupBot  = createLogGroupCmd.Arg("bot", "Name of the bot from which messages are sent to the group.").Required().String()
	createLogGroupName = createLogGroupCmd.Arg("name", "Name of the log group. Must be unique.").Required().String()

	removeGroupCmd  = groupCmd.Command("remove", "Remove a whole group.")
	removeGroupName = removeGroupCmd.Arg("name", "Name of the group to remove.").Required().String()

	addUserToGroupCmd       = groupCmd.Command("invite", "Create an invitation linkt to add an user to a log group. This requires for the user to initiate communication with the bot.")
	adduserToGroupGroupName = addUserToGroupCmd.Arg("group", "Name of the group the user is attached to.").Required().String()

	uninviteUserFromGroupCmd       = groupCmd.Command("uninvite", "Removes a user from a specific group.")
	uninviteUserFromGroupCmdUserID = uninviteUserFromGroupCmd.Arg("id", "Id of the user to remove").Required().Int()
	uninviteUserFromGroupCmdGroup  = uninviteUserFromGroupCmd.Arg("group", "Group from which to remove the user").Required().String()

	listGroupCmd = groupCmd.Command("list", "List all groups.")

	sendCmd   = app.Command("send", "Sends the input read at StdIn to the given group.")
	sendGroup = sendCmd.Arg("group", "Name of the group to send input to.").Required().String()
)

func printError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case registerBotCmd.FullCommand():
		printError(registerBot(*registerBotToken))
	case createLogGroupCmd.FullCommand():
		printError(createGroup(*createLogGroupBot, *createLogGroupName))
	case addUserToGroupCmd.FullCommand():
		printError(addUserToGroup(*adduserToGroupGroupName))
	case listGroupCmd.FullCommand():
		printError(listGroups())
	case removeGroupCmd.FullCommand():
		printError(removeGroup(*removeGroupName))
	case uninviteUserFromGroupCmd.FullCommand():
		printError(uninviteUser(*uninviteUserFromGroupCmdUserID, *uninviteUserFromGroupCmdGroup))
	case sendCmd.FullCommand():
		printError(send(*sendGroup))
	}
}
