package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me al link to it.

In order to get a random page from your list, send me command /rnd.

In order to get all pages from your list, send me command /list
`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand     = "Unknown command 🤔"
	msgNoSavedPages       = "You have no saved pages 🙊"
	msgSaved              = "Saved! 👍"
	msgAlreadyExists      = "You already have this page in your list 👀"
	msgSomethingWentWrong = "❌ Something went wrong! Try again later ❌"
)
