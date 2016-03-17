package main

import (
	"fmt"

	"github.com/dchest/uniuri"

	"gopkg.in/telegram-bot-api.v1"
)

// isGroupAlreadyExists returns whether the given group name was already
// used to create another group
func isGroupAlreadyExists(name string) (bool, error) {
	conf, err := readConfiguration()
	if err != nil {
		return false, err
	}

	for _, elem := range conf.Groups {
		if elem.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// isBotExists returns whether the given botname is already registered with the
// config file
func isBotExists(name string) (bool, error) {
	conf, err := readConfiguration()
	if err != nil {
		return false, err
	}

	for _, elem := range conf.Bots {
		if elem.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func createGroup(bot string, name string) error {
	// first check whether the group already exists
	alreadyExists, err := isGroupAlreadyExists(name)
	if err != nil {
		return err
	}
	if alreadyExists {
		return fmt.Errorf("The group [%s] already exists!", name)
	}

	// then check if the bot exists at all
	botExists, err := isBotExists(bot)
	if err != nil {
		return err
	}
	if !botExists {
		return fmt.Errorf("The bot [%s] does not exists!", bot)
	}

	// if everything went ok one can create the group
	newGroup := groupConfig{bot, name, nil}

	// now the group only has to be saved
	conf, err := readConfiguration()
	if err != nil {
		return err
	}
	conf.Groups = append(conf.Groups, newGroup)
	err = writeConfiguration(conf)
	if err != nil {
		return err
	}

	if *verbose {
		fmt.Printf("Created group [%s]\n", name)
	}
	return nil
}

// given a name this returns the group config object if existing and
// an error if not existing or an error occured
func getGroupByName(name string) (groupConfig, error) {
	conf, err := readConfiguration()
	if err != nil {
		return groupConfig{}, err
	}

	for _, elem := range conf.Groups {
		if elem.Name == name {
			return elem, nil
		}
	}
	return groupConfig{}, fmt.Errorf("There is no group [%s]", name)
}

// given a group id and a user id this adds the user to the group if not already
// existing
func addUserToGroupConfig(groupName string, uid int) error {
	conf, err := readConfiguration()
	if err != nil {
		return err
	}

	for i, group := range conf.Groups {

		// find correct group
		if group.Name == groupName {

			// check if user is already existing
			for _, member := range group.Users {
				// if yes, just do nothing
				if member == uid {
					return nil
				}
			}

			// if not add the user
			// here one has to access the original
			// conf variable, as range copies the group object
			conf.Groups[i].addUser(uid)
			return writeConfiguration(conf)
		}
	}

	// if one reaches here then the group has not existed
	return fmt.Errorf("The group [%s] was not found to append a user", groupName)
}

// addUserToGroup tries to add a new user to the given group
// for this it creates a bot and a link to join the bot.
// the bot then waits for the user to join
// if the user then clicks the link the user is registered
func addUserToGroup(groupName string) error {

	group, err := getGroupByName(groupName)
	if err != nil {
		return err
	}

	token, err := getBotNameByName(group.BotName)
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("could not create bot\n%s\n", err.Error())
	}

	uri := uniuri.NewLen(4)

	url := fmt.Sprintf("https://telegram.me/%s?start=%s", group.BotName, uri)

	fmt.Printf("Listening for user\nGo to page %s\n", url)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Text == "/start "+uri {

			err = addUserToGroupConfig(groupName, update.Message.Chat.ID)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully added user [%d]\n", update.Message.Chat.ID)

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("successfully added to group [%s]", groupName))
			bot.Send(msg)
			return nil
		}
	}
	return nil
}

func listGroups() error {
	conf, err := readConfiguration()
	if err != nil {
		return err
	}

	for _, group := range conf.Groups {
		fmt.Printf("[%s]\nfrom [%s] -> %v\n\n", group.Name, group.BotName, group.Users)
	}

	return nil
}

// given a group name this removes the group with the given name
func removeGroup(name string) error {
	conf, err := readConfiguration()
	if err != nil {
		return err
	}

	// find the group with the correct name
	for i, group := range conf.Groups {
		if group.Name == name {
			conf.Groups = append(conf.Groups[:i], conf.Groups[i+1:]...)
			err = writeConfiguration(conf)
			if err != nil {
				return err
			}
			if *verbose {
				fmt.Printf("Removed group [%s]\n", name)
			}
			return nil
		}
	}

	return fmt.Errorf("The group [%s] was not found", name)
}

// given a user id and a group string this removes the user from the given
// group, given that both exist
func uninviteUser(user int, groupName string) error {
	conf, err := readConfiguration()
	if err != nil {
		return err
	}

	for i, group := range conf.Groups {
		if group.Name == groupName {
			for j, _user := range group.Users {
				if _user == user {
					// remove the user
					conf.Groups[i].Users = append(conf.Groups[i].Users[:j], conf.Groups[i].Users[j+1:]...)

					// save configuration
					err = writeConfiguration(conf)
					if err != nil {
						return err
					}

					if *verbose {
						fmt.Printf("Removed user [%d] from group [%s]\n", user, groupName)
					}
					return nil
				}
			}
		}
	}
	return fmt.Errorf("There is no group/user combination for group [%s] and user [%d]", groupName, user)
}
