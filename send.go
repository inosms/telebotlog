package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"gopkg.in/telegram-bot-api.v1"
)

// send starts the loop to send everything that is written on
// stdin to all the memebers of the group
func send(groupName string) error {

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

	// https://gist.github.com/svett/a95595069e560173a3c8
	info, _ := os.Stdin.Stat()

	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return fmt.Errorf("The command is intended to work with pipes.")
	} else if info.Size() >= 0 {
		reader := bufio.NewReader(os.Stdin)
		redirect(reader, bot, group.Users)
	} else {
		fmt.Errorf("Negative size")
	}
	return nil
}

func redirect(reader *bufio.Reader, bot *tgbotapi.BotAPI, receiver []int) {
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}

		for _, id := range receiver {
			if _, err := bot.Send(tgbotapi.NewMessage(id, input)); err != nil {
				fmt.Printf("Error while sending:\n%s\n", err.Error())
			}
		}
	}
}
