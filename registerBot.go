package main

import (
	"fmt"

	"gopkg.in/telegram-bot-api.v1"
)

// isBotAlreadyRegistered returns whether the given token is already
// registered in the config file or an error if the config file could
// not be read
func isBotAlreadyRegistered(token string) (bool, error) {

	// read config file
	conf, err := readConfiguration()
	if err != nil {
		return false, err
	}

	// search if already registered
	for _, elem := range conf.Bots {
		if elem.Token == token {
			return true, nil
		}
	}
	return false, nil
}

// getBotNameByName returns the token of a bot from the config file
// given it exists and an error otherwise
func getBotNameByName(name string) (string, error) {
	conf, err := readConfiguration()
	if err != nil {
		return "", err
	}

	// search if already registered
	for _, elem := range conf.Bots {
		if elem.Name == name {
			return elem.Token, nil
		}
	}
	return "", fmt.Errorf("botname [%s] was not found in config file", name)
}

// registerBot tries to register a bot in the config file given a token
// if something fails this returns an error
func registerBot(token string) error {

	// check if bot is already registered
	isRegistered, err := isBotAlreadyRegistered(token)
	if err != nil {
		return err
	}
	if isRegistered {
		fmt.Printf("bot [%s] already registered\n", token)
		return nil
	}

	// if not try to register
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("[ERROR] could not create bot\n%s\n", err.Error())
		return err
	}

	// get more bot information such as user name
	botInfo, err := bot.GetMe()
	if err != nil {
		return err
	}
	if *verbose {
		fmt.Printf("[%s] Successfully logged in\n", botInfo.UserName)
	}

	// now save the bot in the configuration file
	// for this read the file at first
	conf, err := readConfiguration()
	if err != nil {
		return err
	}
	// then append the bot to the file
	conf.Bots = append(conf.Bots, botConfig{token, botInfo.UserName})
	// then write the file
	err = writeConfiguration(conf)
	if err != nil {
		return err
	}

	// if nothing crashed until now, the bot should be registered successfully
	fmt.Printf("[%s] Successfully registered\n", botInfo.UserName)
	return nil
}
