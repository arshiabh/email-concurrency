package main

func (app *application) sendEmail(msg Message) {
	app.Wait.Add(1)
	//first write msg to this chan 
	app.Mailer.MailerChan <- msg
}
