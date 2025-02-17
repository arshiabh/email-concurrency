package main

func (app *application) sendEmail(msg Message) {
	//we add to wait gp here 
	// then done() it on the mailer we get value from chan 
	app.Wait.Add(1)
	//trigger chan on blocking mail for loop 
	app.Mailer.MailerChan <- msg
}
