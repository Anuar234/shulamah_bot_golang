package telegram

const msgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me a link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After tham this page will be removed from the list!`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnkownCommand = "Unknown command"
	msgNoSavedPages  = "You have no saved pages"
	msgSaved         = "Saved!"
	msgALreadyExists = "Ypu have already saved this page in the list"
	DownloadDir      = "downloads"
)