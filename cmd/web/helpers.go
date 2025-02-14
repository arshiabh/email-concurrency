package main

func (app *application) sendEmail(msg Message) {
	//we add to wait gp here 
	// then done() it on the mailer we get value from chan 
	app.Wait.Add(1)
	//first write msg to this chan 
	app.Mailer.MailerChan <- msg
}
